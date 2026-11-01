package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsv1ac "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
	metav1ac "k8s.io/client-go/applyconfigurations/meta/v1"
	networkingv1ac "k8s.io/client-go/applyconfigurations/networking/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	fieldManager = "dummysite-controller"
	fetcherImage = "curlimages/curl:8.23.0"
	serverImage  = "nginx:1.31.8-alpine"
	htmlDir      = "/usr/share/nginx/html"
)

var dummySites = schema.GroupVersionResource{Group: "stable.dwk", Version: "v1", Resource: "dummysites"}

// Site is what the controller needs to know about a DummySite
type Site struct {
	Name       string
	Namespace  string
	UID        string
	WebsiteURL string
}

// fatal logs an error and stops the app
func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

// requireEnv returns the value of an environment variable and stops the app if it is not set,
// so that a missing configuration fails at startup instead of falling back to a hidden default
func requireEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fatal("Environment variable is not set", "variable", name)
	}
	return value
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	port := requireEnv("PORT")

	config, err := rest.InClusterConfig()
	if err != nil {
		fatal("Not running in a cluster", "error", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fatal("Could not create the Kubernetes client", "error", err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fatal("Could not create the dynamic client", "error", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// The informer watches the API, and tells about every DummySite that exists or changes.
	// Handlers only queue the name, so that a failed reconcile can be retried with a delay.
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 10*time.Minute)
	informer := factory.ForResource(dummySites).Informer()
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	defer queue.ShutDown()

	enqueue := func(obj any) {
		if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
			queue.Add(key)
		}
	}
	if _, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    enqueue,
		UpdateFunc: func(_, obj any) { enqueue(obj) },
	}); err != nil {
		fatal("Could not watch DummySites", "error", err)
	}

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if !informer.HasSynced() {
			http.Error(w, "Not synced with the API yet", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprint(w, "OK")
	})
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fatal("Server stopped", "error", err)
		}
	}()

	factory.Start(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		fatal("Could not sync with the API")
	}
	slog.Info("Controller started")

	go func() {
		<-ctx.Done()
		queue.ShutDown()
	}()
	for {
		key, shutdown := queue.Get()
		if shutdown {
			return
		}
		if err := reconcile(ctx, client, informer.GetStore(), key); err != nil {
			slog.Error("Could not reconcile, retrying", "dummysite", key, "error", err)
			queue.AddRateLimited(key)
		} else {
			queue.Forget(key)
		}
		queue.Done(key)
	}
}

// reconcile makes the cluster match a DummySite. A deleted DummySite needs nothing here: Kubernetes
// removes the resources through their owner references.
func reconcile(ctx context.Context, client kubernetes.Interface, store cache.Store, key string) error {
	obj, exists, err := store.GetByKey(key)
	if err != nil || !exists {
		return err
	}
	site, err := parseSite(obj.(*unstructured.Unstructured))
	if err != nil {
		// Retrying does not help, so wait for the DummySite to change
		slog.Error("Invalid DummySite", "dummysite", key, "error", err)
		return nil
	}

	// Server-side apply creates a resource or updates it to this state, so running this twice is safe
	options := metav1.ApplyOptions{FieldManager: fieldManager, Force: true}
	if _, err := client.AppsV1().Deployments(site.Namespace).Apply(ctx, deployment(site), options); err != nil {
		return fmt.Errorf("deployment: %w", err)
	}
	if _, err := client.CoreV1().Services(site.Namespace).Apply(ctx, service(site), options); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if _, err := client.NetworkingV1().Ingresses(site.Namespace).Apply(ctx, ingress(site), options); err != nil {
		return fmt.Errorf("ingress: %w", err)
	}
	slog.Info("DummySite reconciled", "dummysite", key, "website_url", site.WebsiteURL)
	return nil
}

func parseSite(obj *unstructured.Unstructured) (Site, error) {
	websiteURL, _, err := unstructured.NestedString(obj.Object, "spec", "website_url")
	if err != nil {
		return Site{}, err
	}
	if err := validateURL(websiteURL); err != nil {
		return Site{}, err
	}
	return Site{Name: obj.GetName(), Namespace: obj.GetNamespace(), UID: string(obj.GetUID()), WebsiteURL: websiteURL}, nil
}

func validateURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil {
		return err
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return fmt.Errorf("website_url must be an http or https URL, got %q", value)
	}
	return nil
}

func owner(site Site) *metav1ac.OwnerReferenceApplyConfiguration {
	return metav1ac.OwnerReference().
		WithAPIVersion(dummySites.GroupVersion().String()).
		WithKind("DummySite").
		WithName(site.Name).
		WithUID(types.UID(site.UID)).
		WithController(true)
}

func labels(site Site) map[string]string {
	return map[string]string{"app": "dummysite-" + site.Name}
}

// deployment fetches the page in an init container and serves it with nginx. A new website_url
// changes the Pod template, so Kubernetes starts a new Pod that fetches the new page.
func deployment(site Site) *appsv1ac.DeploymentApplyConfiguration {
	return appsv1ac.Deployment(site.Name+"-dep", site.Namespace).
		WithOwnerReferences(owner(site)).
		WithLabels(labels(site)).
		WithSpec(appsv1ac.DeploymentSpec().
			WithReplicas(1).
			WithSelector(metav1ac.LabelSelector().WithMatchLabels(labels(site))).
			WithTemplate(corev1ac.PodTemplateSpec().
				WithLabels(labels(site)).
				WithSpec(corev1ac.PodSpec().
					WithVolumes(corev1ac.Volume().WithName("html").WithEmptyDir(corev1ac.EmptyDirVolumeSource())).
					WithInitContainers(corev1ac.Container().
						WithName("fetch").
						WithImage(fetcherImage).
						WithCommand("curl", "--fail", "--silent", "--show-error", "--location", "--max-time", "30",
							"--output", "/html/index.html", "--url", site.WebsiteURL).
						WithVolumeMounts(corev1ac.VolumeMount().WithName("html").WithMountPath("/html"))).
					WithContainers(corev1ac.Container().
						WithName("nginx").
						WithImage(serverImage).
						WithPorts(corev1ac.ContainerPort().WithContainerPort(80)).
						WithVolumeMounts(corev1ac.VolumeMount().WithName("html").WithMountPath(htmlDir).WithReadOnly(true)).
						WithReadinessProbe(corev1ac.Probe().
							WithHTTPGet(corev1ac.HTTPGetAction().WithPath("/").WithPort(intstr.FromInt32(80))))))))
}

func service(site Site) *corev1ac.ServiceApplyConfiguration {
	return corev1ac.Service(site.Name+"-svc", site.Namespace).
		WithOwnerReferences(owner(site)).
		WithLabels(labels(site)).
		WithSpec(corev1ac.ServiceSpec().
			WithType(corev1.ServiceTypeClusterIP).
			WithSelector(labels(site)).
			WithPorts(corev1ac.ServicePort().WithPort(80).WithTargetPort(intstr.FromInt32(80))))
}

// ingress publishes the site at <name>.localhost, which browsers resolve to this machine
func ingress(site Site) *networkingv1ac.IngressApplyConfiguration {
	return networkingv1ac.Ingress(site.Name+"-ingress", site.Namespace).
		WithOwnerReferences(owner(site)).
		WithLabels(labels(site)).
		WithSpec(networkingv1ac.IngressSpec().
			WithRules(networkingv1ac.IngressRule().
				WithHost(site.Name + ".localhost").
				WithHTTP(networkingv1ac.HTTPIngressRuleValue().
					WithPaths(networkingv1ac.HTTPIngressPath().
						WithPath("/").
						WithPathType(networkingv1.PathTypePrefix).
						WithBackend(networkingv1ac.IngressBackend().
							WithService(networkingv1ac.IngressServiceBackend().
								WithName(site.Name + "-svc").
								WithPort(networkingv1ac.ServiceBackendPort().WithNumber(80))))))))
}

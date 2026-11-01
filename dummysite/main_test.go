package main

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func dummySite(url string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "stable.dwk/v1",
		"kind":       "DummySite",
		"metadata":   map[string]any{"name": "example", "namespace": "default", "uid": "abc"},
		"spec":       map[string]any{"website_url": url},
	}}
}

func TestValidateURL(t *testing.T) {
	for _, value := range []string{"https://example.com/", "http://example.com/a?b=c"} {
		if err := validateURL(value); err != nil {
			t.Errorf("validateURL(%q) = %v", value, err)
		}
	}
	for _, value := range []string{"", "example.com", "file:///etc/passwd", "ftp://example.com/", "https://"} {
		if err := validateURL(value); err == nil {
			t.Errorf("validateURL(%q) accepted", value)
		}
	}
}

func TestReconcileCreatesOwnedResources(t *testing.T) {
	client := fake.NewClientset()
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	if err := store.Add(dummySite("https://example.com/")); err != nil {
		t.Fatal(err)
	}

	// Twice, since a reconcile has to be safe to repeat
	for range 2 {
		if err := reconcile(context.Background(), client, store, "default/example"); err != nil {
			t.Fatal(err)
		}
	}

	ctx := context.Background()
	dep, err := client.AppsV1().Deployments("default").Get(ctx, "example-dep", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := dep.OwnerReferences; len(got) != 1 || got[0].Kind != "DummySite" || got[0].Name != "example" {
		t.Errorf("owner references = %+v", got)
	}
	args := dep.Spec.Template.Spec.InitContainers[0].Command
	if args[len(args)-1] != "https://example.com/" {
		t.Errorf("the fetcher does not get the URL: %v", args)
	}
	if _, err := client.CoreV1().Services("default").Get(ctx, "example-svc", metav1.GetOptions{}); err != nil {
		t.Error(err)
	}
	if _, err := client.NetworkingV1().Ingresses("default").Get(ctx, "example-ingress", metav1.GetOptions{}); err != nil {
		t.Error(err)
	}
}

func TestReconcileOfDeletedSiteDoesNothing(t *testing.T) {
	client := fake.NewClientset()
	if err := reconcile(context.Background(), client, cache.NewStore(cache.MetaNamespaceKeyFunc), "default/gone"); err != nil {
		t.Fatal(err)
	}
	if actions := client.Actions(); len(actions) != 0 {
		t.Errorf("actions = %v", actions)
	}
}

func TestReconcileIgnoresInvalidURL(t *testing.T) {
	client := fake.NewClientset()
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	if err := store.Add(dummySite("file:///etc/passwd")); err != nil {
		t.Fatal(err)
	}
	if err := reconcile(context.Background(), client, store, "default/example"); err != nil {
		t.Fatal(err)
	}
	if actions := client.Actions(); len(actions) != 0 {
		t.Errorf("actions = %v", actions)
	}
}

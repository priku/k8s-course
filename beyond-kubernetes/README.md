# Beyond Kubernetes

## Platform comparison: Rancher vs OpenShift (exercise 5.5)

**My pick: Rancher.** The exercise asks to pick one and argue for it, so this is deliberately one-sided.

- **Works with the Kubernetes you already have.** Rancher manages any conformant cluster, including imported AKS, EKS and GKE clusters. OpenShift is a platform of its own that you move onto.
- **It is plain Kubernetes.** What is learned on this course applies directly: Deployments, Ingress and the Gateway API. OpenShift adds its own concepts on top, such as Routes, Projects and DeploymentConfigs, and its own `oc` command.
- **Fully open source, without a subscription to get started.** All of Rancher is free to run, and support is an optional purchase. OpenShift is a paid subscription product, and its community version OKD does not come with the same support.
- **Scales down as well as up.** The same family runs from K3s on a Raspberry Pi or in k3d, as on this course, up to RKE2 for strict environments. OpenShift has much larger minimum requirements for a full cluster.
- **Many clusters from one place.** Managing a fleet across clouds and data centres is the very thing Rancher was built for, including GitOps with Fleet.
- **Less lock-in.** Leaving Rancher leaves the clusters running as ordinary Kubernetes. Leaving OpenShift means migrating away from everything that is specific to it.

What OpenShift does better, to be fair:

- More is included and integrated from the start: image builds, a registry, pipelines and monitoring.
- Stricter security by default, such as not running containers as root.
- A single vendor supports the whole stack from the operating system up, which large regulated organizations often require.

For a small team that wants to stay close to upstream Kubernetes, these do not outweigh the cost and the lock-in.

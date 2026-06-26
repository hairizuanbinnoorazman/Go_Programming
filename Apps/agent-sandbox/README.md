# Agent Sandbox GKE UI

A small Go web application for demonstrating Agent Sandbox and GKE Pod Snapshot lifecycle operations from buttons in a browser.

The dashboard:

- Creates and destroys `SandboxClaim` resources.
- Suspends and resumes a Sandbox by setting `spec.replicas` to `0` or `1`.
- Creates GKE `PodSnapshotManualTrigger` resources and waits for checkpoint completion.
- Captures a memory snapshot and then suspends the Sandbox.
- Restores the latest ready snapshot by resuming a suspended Sandbox or deleting its current Pod so Agent Sandbox recreates it.
- Lists every live Pod in the namespace, including runtime class, node, IP, last log line, and GKE `PodRestored` condition.
- Lists GKE PodSnapshots and their Cloud Storage paths.

## Existing GKE environment

The repository defaults match the verified cluster:

| Setting | Value |
| --- | --- |
| Project | `new-demo-project-462517` |
| Cluster | `agent-sandbox-standard` |
| Zone | `asia-southeast1-a` |
| Namespace | `agent-sandbox-snapshots` |
| SandboxTemplate / warm pool | `snapshot-python` |
| Artifact Registry repository | `asia-southeast1-docker.pkg.dev/new-demo-project-462517/demo` |

The live cluster serves Agent Sandbox `v1alpha1`. This differs from the `v1beta1` examples in the June 24 blog draft. The manifests and code intentionally use the API actually installed by Agent Sandbox `v0.4.6`.

## Deploy

Prerequisites: authenticated `gcloud`, `kubectl`, Go 1.26+, Cloud Build access, and permission to push to the existing Artifact Registry repository.

```bash
make deploy
```

`make deploy` runs tests, builds and pushes the container with Cloud Build, refreshes GKE credentials, applies the namespaced RBAC/Deployment/Service, sets the exact image tag, and waits for rollout.

Override defaults when needed:

```bash
make deploy PROJECT_ID=my-project ZONE=asia-southeast1-a CLUSTER=my-cluster \
  NAMESPACE=agent-sandbox-snapshots AR_REPOSITORY=demo
```

The local workstation currently lacks `gke-gcloud-auth-plugin`. Make targets therefore pass a short-lived `gcloud auth print-access-token` token directly to `kubectl`.

## Access without Ingress or LoadBalancer

Forward the ClusterIP Service to localhost:

```bash
make port-forward
```

Open <http://localhost:8080>. Keep the terminal running while using the UI.

Equivalent direct command:

```bash
kubectl --token="$(gcloud auth print-access-token)" \
  --namespace=agent-sandbox-snapshots \
  port-forward service/agent-sandbox-demo 8080:8080
```

No Ingress, Gateway, public IP, or cloud load balancer is created.

## Lifecycle semantics

- **Suspend**: sets Sandbox replicas to zero. The Pod is removed. This does not checkpoint process memory.
- **Resume**: sets replicas to one. Without a compatible snapshot this starts normally; with a compatible latest GKE snapshot, GKE can restore it.
- **Snapshot memory**: creates a manual GKE Pod Snapshot. The configured policy uses `postCheckpoint: resume`, so the current Pod continues after the checkpoint.
- **Snapshot + suspend**: waits until the memory/rootfs snapshot is complete, then scales the Sandbox to zero.
- **Restore latest**: verifies a ready snapshot exists. If suspended, it scales the Sandbox to one. If already running, it deletes the Pod with zero grace period and Agent Sandbox recreates the same Pod template, allowing GKE to restore the latest compatible snapshot.
- **Destroy**: deletes the claim with foreground propagation. Claims created by this UI use `shutdownPolicy: DeleteForeground`.

The Pod table exposes the counter container's latest `MEMORY_STATE` log line. To demonstrate restoration, note its session UUID and count, create a snapshot, let the count advance, then restore. The recreated Pod should report `PodRestored=True`, retain the same UUID, and continue from the checkpointed count.

## Security and scope

The application is intended for local demonstration through `kubectl port-forward`. It has no user authentication. Its Service is `ClusterIP`, and RBAC is restricted to the configured namespace and only the resources/actions required by the UI.

Do not expose it publicly without adding authentication, authorization, CSRF protection, and TLS.

## Local development

With a working kubeconfig:

```bash
GOWORK=off go run . --kubeconfig="$HOME/.kube/config"
```

Then open <http://localhost:8080>.

## Remove the UI

```bash
make undeploy
```

This removes only the UI Service, Deployment, ServiceAccount, Role, and RoleBinding. It does not remove Agent Sandbox, the warm pool, claims, snapshots, bucket data, or the GKE cluster.

## Sources

- `~/Documents/blog/content/posts/20260624_deployAgentSandboxGKEStandard.md`
- `~/Documents/agent-sandbox/docs/gke-standard-gvisor-snapshot-validation-2026-06-21.md`
- `~/Documents/agent-sandbox` API, controller, and GKE snapshot client sources
- [GKE Pod snapshots](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-snapshots)
- [Agent Sandbox Pod snapshots](https://cloud.google.com/kubernetes-engine/docs/how-to/agent-sandbox-pod-snapshots)

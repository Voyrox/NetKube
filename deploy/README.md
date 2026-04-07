# Deploy NetKube

This folder contains the Kubernetes manifests for NetKube and quick commands for running it with Docker.

## Docker

Build the image:

```bash
docker build -t netkube:latest .
```

Run it locally:

```bash
docker run --rm -p 3000:3000 -v netkube-config:/app/config netkube:latest
```

Open `http://localhost:3000`.

## Kubernetes

1. Build and push the image to a registry your cluster can pull from.
2. Update the image value in `deploy/kubernetes/deployment.yaml`.
3. Apply the manifests.

Example:

```bash
docker build -t <registry>/netkube:latest .
docker push <registry>/netkube:latest
kubectl apply -f deploy/kubernetes/namespace.yaml
kubectl apply -f deploy/kubernetes/pvc.yaml
kubectl apply -f deploy/kubernetes/deployment.yaml
kubectl apply -f deploy/kubernetes/service.yaml
kubectl -n netkube port-forward svc/netkube 3000:3000
```

Then open `http://localhost:3000`.

## Notes

- The app stores uploaded kubeconfig files and selected contexts in `/app/config`.
- `deploy/kubernetes/pvc.yaml` creates persistent storage for that data.
- The default service is `ClusterIP`, so `kubectl port-forward` is the easiest first way to access it.

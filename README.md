# NetKube

## Run with Docker

```bash
docker build -t netkube:latest .
docker run --rm -p 3000:3000 -v netkube-config:/app/config netkube:latest
```

## Deploy to Kubernetes

Build the image and push it to a registry your cluster can pull from, then update the image in `deploy/kubernetes/deployment.yaml`.

```bash
docker build -t <registry>/netkube:latest .
docker push <registry>/netkube:latest
kubectl apply -f deploy/kubernetes/namespace.yaml
kubectl apply -f deploy/kubernetes/pvc.yaml
kubectl apply -f deploy/kubernetes/deployment.yaml
kubectl apply -f deploy/kubernetes/service.yaml
kubectl -n netkube port-forward svc/netkube 3000:3000
```

The app persists uploaded kubeconfig files and selected contexts under `/app/config`, so the deployment mounts a persistent volume claim for that data.

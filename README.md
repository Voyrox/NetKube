# NetKube

NetKube is a simple Kubernetes dashboard for viewing cluster health, nodes, pods, deployments, services, and multi-context kubeconfig data.

It provides a lightweight web interface for inspecting cluster resources, switching between stored kubeconfig contexts, and reviewing operational details without the overhead of a larger Kubernetes management platform.

Create a `.env` file before starting the app. NetKube requires `EMAIL` (or legacy `USERNAME`) and `PASSWORD` for login. If `SESSION_SECRET` is not set, NetKube generates and persists one automatically under `config/session_secret`.

## What does it look like?
<table>
	<tr>
		<td><img src="./assets/1.png" alt="NetKube screenshot 1" width="100%" /></td>
		<td><img src="./assets/2.png" alt="NetKube screenshot 2" width="100%" /></td>
	</tr>
	<tr>
		<td><img src="./assets/3.png" alt="NetKube screenshot 3" width="100%" /></td>
		<td><img src="./assets/4.png" alt="NetKube screenshot 4" width="100%" /></td>
	</tr>
</table>

## Run with Docker

```bash
docker build -t netkube:latest .
docker run --rm -p 3000:3000 --env-file .env -v netkube-config:/app/config netkube:latest
```

## Deploy to Kubernetes

```bash
kubectl apply -f netkube.yaml
kubectl -n netkube port-forward svc/netkube 3000:3000
```

The app persists uploaded kubeconfig files, selected contexts, and the generated session secret under `/app/config`, so the deployment mounts a persistent volume claim for that data.

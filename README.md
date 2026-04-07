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
```

The default manifest exposes NetKube with a `NodePort` service on port `30080`, so you can usually open `http://<node-ip>:30080` after deploying.

The app persists uploaded kubeconfig files, selected contexts, and the generated session secret under `/app/config`, so the deployment mounts a persistent volume claim for that data.

## Reverse proxy notes

NetKube sets secure session cookies automatically when requests arrive over HTTPS. If you deploy it behind a reverse proxy, make sure the proxy forwards `X-Forwarded-Proto` so NetKube can detect HTTPS correctly.

For local HTTP access, NetKube falls back to non-secure cookies so login still works without TLS.

### NGINX example

```nginx
server {
    listen 80;
    server_name netkube.example.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name netkube.example.com;

    ssl_certificate /etc/letsencrypt/live/netkube.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/netkube.example.com/privkey.pem;

    location / {
        proxy_pass http://192.168.1.114:30080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Pangolin example

If you are using Pangolin as the reverse proxy in front of NetKube, create a public HTTP resource that targets `http://192.168.1.114:30080`.

- Connection method: `Local` or the site that can reach your Kubernetes node
- Target address: `192.168.1.114:30080`
- Custom headers:

```text
X-Forwarded-Proto: https
```

If your NetKube resource is published on `https://netkube.example.com`, Pangolin should terminate TLS at the edge and forward requests to the NodePort target while preserving the original host.

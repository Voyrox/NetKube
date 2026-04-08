let deploymentItems = [];

const DEFAULT_DEPLOYMENT_YAML = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: default
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80`;

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("deploymentSearch");

  initCreateResourceModal({
    triggerId: "createDeploymentButton",
    title: "Create Deployment",
    description: "Edit the deployment YAML before confirming.",
    initialValue: DEFAULT_DEPLOYMENT_YAML,
    confirmLabel: "Confirm",
    pendingLabel: "Creating...",
    async onConfirm(content) {
      const data = await fetchClusterData("/api/workloads/deployments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content }),
      });

      window.setTimeout(() => window.location.reload(), 150);
      return {
        message: `Created deployment ${data.name || "resource"} in ${data.namespace || "default"}.`,
      };
    },
  });

  searchInput?.addEventListener("input", () =>
    renderFilteredDeployments(searchInput.value),
  );

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/deployments");
      deploymentItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "deploymentsNamespaceFilter" });
      setText(
        "deploymentsHeroTitle",
        `Deployments (${data.stats?.healthy || 0} healthy)`,
      );
      renderFilteredDeployments(searchInput?.value || "");
    } catch (error) {
      deploymentItems = [];
      renderDeploymentsTable([], error.message || "Failed to load deployments");
    }
  });
});

function renderFilteredDeployments(query) {
  const filteredItems = filterDeployments(deploymentItems, query);
  renderDeploymentsTable(
    filteredItems,
    filteredItems.length ? "" : "No deployments match your search",
  );
  setText("deploymentsTableCount", `${filteredItems.length} rows`);
}

function filterDeployments(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;

  return items.filter((item) =>
    [item.namespace, item.name, item.status].some((value) =>
      String(value || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}

function renderDeploymentsTable(items, message) {
  const body = document.getElementById("deploymentsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="8">${escapeHtml(message || "No deployments found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.ready)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(String(item.desired))}</td>
      <td>${escapeHtml(String(item.updated))}</td>
      <td>${escapeHtml(String(item.available))}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    row.addEventListener("click", () => openDeployment(item));
    row.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openDeployment(item);
      }
    });
    body.appendChild(row);
  });
}

function openDeployment(item) {
  const params = new URLSearchParams({
    namespace: item.namespace || "",
    name: item.name || "",
  });
  window.location.href = `/workloads/manage/deployment?${params.toString()}`;
}

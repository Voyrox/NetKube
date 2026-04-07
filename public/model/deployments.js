document.addEventListener("DOMContentLoaded", async () => {
  const searchInput = document.getElementById("deploymentSearch");

  try {
    const data = await fetchClusterData("/api/workloads/deployments");
    applyPageMeta(data.meta, { namespaceId: "deploymentsNamespaceFilter" });
    setText("deploymentsTableCount", `${data.count || 0} rows`);
    setText("deploymentsHeroTitle", `Deployments (${data.stats?.healthy || 0} healthy)`);
    renderDeploymentsTable(data.items || []);

    if (searchInput) {
      searchInput.addEventListener("input", () => {
        const filteredItems = filterDeployments(data.items || [], searchInput.value);
        renderDeploymentsTable(filteredItems, filteredItems.length ? "" : "No deployments match your search");
        setText("deploymentsTableCount", `${filteredItems.length} rows`);
      });
    }
  } catch (error) {
    renderDeploymentsTable([], error.message || "Failed to load deployments");
  }
});

function filterDeployments(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;

  return items.filter((item) => [item.namespace, item.name, item.status].some((value) => String(value || "").toLowerCase().includes(normalized)));
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
    name: item.name || ""
  });

  window.location.href = `/workloads/manage/deployment?${params.toString()}`;
}

document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/workloads/deployments");
    applyPageMeta(data.meta, { namespaceId: "deploymentsNamespaceFilter" });
    setText("deploymentsTableCount", `${data.count || 0} rows`);
    setText("deploymentsHeroTitle", `Deployments (${data.stats?.healthy || 0} healthy)`);
    renderDeploymentsTable(data.items || []);
  } catch (error) {
    renderDeploymentsTable([], error.message || "Failed to load deployments");
  }
});

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
    body.appendChild(row);
  });
}

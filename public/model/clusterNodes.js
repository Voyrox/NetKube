document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/cluster/nodes");
    applyPageMeta(data.meta);
    setText("nodesHeroTitle", `Nodes (${data.count || 0})`);
    setText("nodesTableCount", `${data.count || 0} rows`);
    renderNodesTable(data.items || []);
  } catch (error) {
    renderNodesTable([], error.message || "Failed to load nodes");
  }
});

function renderNodesTable(items, message) {
  const body = document.getElementById("nodesTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="6">${escapeHtml(message || "No nodes found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(item.role)}</td>
      <td>${escapeHtml(item.version)}</td>
      <td>${escapeHtml(item.internalIP)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    row.addEventListener("click", () => openNode(item));
    row.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openNode(item);
      }
    });
    body.appendChild(row);
  });
}

function openNode(item) {
  const params = new URLSearchParams({ name: item.name || "" });
  window.location.href = `/clusters/details/node?${params.toString()}`;
}

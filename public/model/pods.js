document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/workloads/pods");
    applyPageMeta(data.meta, { namespaceId: "podsNamespaceFilter" });
    setText("podsTableCount", `${data.count || 0} rows`);
    setText("podsHeroTitle", `Pods (${data.stats?.running || 0} running)`);
    renderPodsTable(data.items || []);
  } catch (error) {
    renderPodsTable([], error.message || "Failed to load pods");
  }
});

function renderPodsTable(items, message) {
  const body = document.getElementById("podsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="10">${escapeHtml(message || "No pods found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.ready)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(String(item.restarts))}</td>
      <td>${escapeHtml(item.lastRestart)}</td>
      <td>${escapeHtml(item.lastRestartReason)}</td>
      <td>${escapeHtml(item.node)}</td>
      <td>${escapeHtml(item.podIP)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

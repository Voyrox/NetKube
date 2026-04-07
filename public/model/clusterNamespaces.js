document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/cluster/namespaces");
    applyPageMeta(data.meta);
    setText("namespacesHeroTitle", `Namespaces (${data.count || 0})`);
    setText("namespacesTableCount", `${data.count || 0} rows`);
    renderNamespacesTable(data.items || []);
  } catch (error) {
    renderNamespacesTable([], error.message || "Failed to load namespaces");
  }
});

function renderNamespacesTable(items, message) {
  const body = document.getElementById("namespacesTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="3">${escapeHtml(message || "No namespaces found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.phase)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

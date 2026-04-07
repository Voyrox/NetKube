document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/cluster/leases");
    applyPageMeta(data.meta);
    setText("leasesHeroTitle", `Leases (${data.count || 0})`);
    setText("leasesTableCount", `${data.count || 0} rows`);
    renderLeasesTable(data.items || []);
  } catch (error) {
    renderLeasesTable([], error.message || "Failed to load leases");
  }
});

function renderLeasesTable(items, message) {
  const body = document.getElementById("leasesTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="5">${escapeHtml(message || "No leases found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.holder)}</td>
      <td>${escapeHtml(item.lastRenew)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

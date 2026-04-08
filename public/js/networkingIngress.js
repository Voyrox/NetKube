let ingressItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("ingressSearch");

  searchInput?.addEventListener("input", () => {
    renderFilteredIngress(searchInput.value);
  });

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/networking/ingress");
      ingressItems = data.items || [];
      applyPageMeta(data.meta);
      setText("ingressHeroTitle", `Ingress (${data.count || 0})`);
      renderFilteredIngress(searchInput?.value || "");
    } catch (error) {
      ingressItems = [];
      renderIngressTable([], error.message || "Failed to load ingress");
    }
  });
});

function renderFilteredIngress(query) {
  const filteredItems = filterIngress(ingressItems, query);
  renderIngressTable(filteredItems, filteredItems.length ? "" : "No ingress match your search");
  setText("ingressTableCount", `${filteredItems.length} rows`);
}

function filterIngress(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;

  return items.filter((item) => [item.namespace, item.name, item.class, item.hosts, item.address].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderIngressTable(items, message) {
  const body = document.getElementById("ingressTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="6">${escapeHtml(message || "No ingress found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.class)}</td>
      <td>${escapeHtml(item.hosts)}</td>
      <td>${escapeHtml(item.address)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

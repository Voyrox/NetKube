let configMapItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("configMapsSearch");

  searchInput?.addEventListener("input", () => {
    renderFilteredConfigMaps(searchInput.value);
  });

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/configuration/configmaps");
      configMapItems = data.items || [];
      applyPageMeta(data.meta);
      setText("configMapsHeroTitle", `ConfigMaps (${data.count || 0})`);
      renderFilteredConfigMaps(searchInput?.value || "");
    } catch (error) {
      configMapItems = [];
      renderConfigMapsTable([], error.message || "Failed to load configmaps");
    }
  });
});

function renderFilteredConfigMaps(query) {
  const filteredItems = filterConfigMaps(configMapItems, query);
  renderConfigMapsTable(filteredItems, filteredItems.length ? "" : "No configmaps match your search");
  setText("configMapsTableCount", `${filteredItems.length} rows`);
}

function filterConfigMaps(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;

  return items.filter((item) => [item.namespace, item.name, item.immutable].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderConfigMapsTable(items, message) {
  const body = document.getElementById("configMapsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="5">${escapeHtml(message || "No configmaps found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(String(item.data))}</td>
      <td>${escapeHtml(item.immutable)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

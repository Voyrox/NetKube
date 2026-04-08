let csiDriverItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("csiDriversSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredCsiDrivers(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/storage/csidrivers");
      csiDriverItems = data.items || [];
      applyPageMeta(data.meta);
      setText("csiDriversHeroTitle", `CSI Drivers (${data.count || 0})`);
      renderFilteredCsiDrivers(searchInput?.value || "");
    } catch (error) {
      csiDriverItems = [];
      renderCsiDriversTable([], error.message || "Failed to load CSI drivers");
    }
  });
});
function renderFilteredCsiDrivers(query) {
  const items = filterByQuery(csiDriverItems, query, [
    "name",
    "attachRequired",
    "modes",
  ]);
  renderCsiDriversTable(
    items,
    items.length ? "" : "No CSI drivers match your search",
  );
  setText("csiDriversTableCount", `${items.length} rows`);
}
function renderCsiDriversTable(items, message) {
  renderTableRows(
    "csiDriversTableBody",
    6,
    items,
    message || "No CSI drivers found",
    (item) =>
      `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.attachRequired)}</td><td>${escapeHtml(item.podInfoOnMount)}</td><td>${escapeHtml(item.storageCapacity)}</td><td>${escapeHtml(item.modes)}</td><td>${escapeHtml(item.age)}</td>`,
  );
}

let csiNodeItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("csiNodesSearch");
  searchInput?.addEventListener("input", () => renderFilteredCsiNodes(searchInput.value));
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/storage/csinodes");
      csiNodeItems = data.items || [];
      applyPageMeta(data.meta);
      setText("csiNodesHeroTitle", `CSI Nodes (${data.count || 0})`);
      renderFilteredCsiNodes(searchInput?.value || "");
    } catch (error) {
      csiNodeItems = [];
      renderCsiNodesTable([], error.message || "Failed to load CSI nodes");
    }
  });
});
function renderFilteredCsiNodes(query) { const items = filterByQuery(csiNodeItems, query, ["name"]); renderCsiNodesTable(items, items.length ? "" : "No CSI nodes match your search"); setText("csiNodesTableCount", `${items.length} rows`); }
function renderCsiNodesTable(items, message) { renderTableRows("csiNodesTableBody", 3, items, message || "No CSI nodes found", (item) => `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(String(item.drivers))}</td><td>${escapeHtml(item.age)}</td>`); }

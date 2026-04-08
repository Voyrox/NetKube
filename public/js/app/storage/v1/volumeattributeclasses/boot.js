let volumeAttributeClassItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("volumeAttributeClassesSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredVolumeAttributeClasses(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData(
        "/api/storage/volumeattributeclasses",
      );
      volumeAttributeClassItems = data.items || [];
      applyPageMeta(data.meta);
      setText(
        "volumeAttributeClassesHeroTitle",
        `Volume Attribute Classes (${data.count || 0})`,
      );
      renderFilteredVolumeAttributeClasses(searchInput?.value || "");
    } catch (error) {
      volumeAttributeClassItems = [];
      renderVolumeAttributeClassesTable(
        [],
        error.message || "Failed to load volume attribute classes",
      );
    }
  });
});
function renderFilteredVolumeAttributeClasses(query) {
  const items = filterByQuery(volumeAttributeClassItems, query, [
    "name",
    "driverName",
  ]);
  renderVolumeAttributeClassesTable(
    items,
    items.length ? "" : "No volume attribute classes match your search",
  );
  setText("volumeAttributeClassesTableCount", `${items.length} rows`);
}
function renderVolumeAttributeClassesTable(items, message) {
  renderTableRows(
    "volumeAttributeClassesTableBody",
    4,
    items,
    message || "No volume attribute classes found",
    (item) =>
      `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.driverName)}</td><td>${escapeHtml(String(item.parameters))}</td><td>${escapeHtml(item.age)}</td>`,
  );
}

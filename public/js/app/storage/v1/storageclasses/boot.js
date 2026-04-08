let storageClassItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("storageClassesSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredStorageClasses(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/storage/storageclasses");
      storageClassItems = data.items || [];
      applyPageMeta(data.meta);
      setText(
        "storageClassesHeroTitle",
        `Storage Classes (${data.count || 0})`,
      );
      renderFilteredStorageClasses(searchInput?.value || "");
    } catch (error) {
      storageClassItems = [];
      renderStorageClassesTable(
        [],
        error.message || "Failed to load storage classes",
      );
    }
  });
});
function renderFilteredStorageClasses(query) {
  const items = filterByQuery(storageClassItems, query, [
    "name",
    "provisioner",
    "reclaimPolicy",
    "bindingMode",
    "default",
  ]);
  renderStorageClassesTable(
    items,
    items.length ? "" : "No storage classes match your search",
  );
  setText("storageClassesTableCount", `${items.length} rows`);
}
function renderStorageClassesTable(items, message) {
  renderTableRows(
    "storageClassesTableBody",
    6,
    items,
    message || "No storage classes found",
    (item) =>
      `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.provisioner)}</td><td>${escapeHtml(item.reclaimPolicy)}</td><td>${escapeHtml(item.bindingMode)}</td><td>${escapeHtml(item.default)}</td><td>${escapeHtml(item.age)}</td>`,
  );
}

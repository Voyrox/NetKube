let persistentVolumeItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("persistentVolumesSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredPersistentVolumes(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/storage/persistentvolumes");
      persistentVolumeItems = data.items || [];
      applyPageMeta(data.meta);
      setText(
        "persistentVolumesHeroTitle",
        `Persistent Volumes (${data.count || 0})`,
      );
      renderFilteredPersistentVolumes(searchInput?.value || "");
    } catch (error) {
      persistentVolumeItems = [];
      renderPersistentVolumesTable(
        [],
        error.message || "Failed to load persistent volumes",
      );
    }
  });
});
function renderFilteredPersistentVolumes(query) {
  const items = filterByQuery(persistentVolumeItems, query, [
    "name",
    "status",
    "claim",
    "storageClass",
  ]);
  renderPersistentVolumesTable(
    items,
    items.length ? "" : "No persistent volumes match your search",
  );
  setText("persistentVolumesTableCount", `${items.length} rows`);
}
function renderPersistentVolumesTable(items, message) {
  renderTableRows(
    "persistentVolumesTableBody",
    8,
    items,
    message || "No persistent volumes found",
    (item) =>
      `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.status)}</td><td>${escapeHtml(item.capacity)}</td><td>${escapeHtml(item.accessModes)}</td><td>${escapeHtml(item.reclaimPolicy)}</td><td>${escapeHtml(item.claim)}</td><td>${escapeHtml(item.storageClass)}</td><td>${escapeHtml(item.age)}</td>`,
  );
}

let persistentVolumeClaimItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("persistentVolumeClaimsSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredPersistentVolumeClaims(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData(
        "/api/storage/persistentvolumeclaims",
      );
      persistentVolumeClaimItems = data.items || [];
      applyPageMeta(data.meta);
      setText(
        "persistentVolumeClaimsHeroTitle",
        `Persistent Volume Claims (${data.count || 0})`,
      );
      renderFilteredPersistentVolumeClaims(searchInput?.value || "");
    } catch (error) {
      persistentVolumeClaimItems = [];
      renderPersistentVolumeClaimsTable(
        [],
        error.message || "Failed to load persistent volume claims",
      );
    }
  });
});
function renderFilteredPersistentVolumeClaims(query) {
  const items = filterByQuery(persistentVolumeClaimItems, query, [
    "namespace",
    "name",
    "status",
    "volume",
    "storageClass",
  ]);
  renderPersistentVolumeClaimsTable(
    items,
    items.length ? "" : "No persistent volume claims match your search",
  );
  setText("persistentVolumeClaimsTableCount", `${items.length} rows`);
}
function renderPersistentVolumeClaimsTable(items, message) {
  renderTableRows(
    "persistentVolumeClaimsTableBody",
    8,
    items,
    message || "No persistent volume claims found",
    (item) =>
      `<td>${escapeHtml(item.namespace)}</td><td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.status)}</td><td>${escapeHtml(item.volume)}</td><td>${escapeHtml(item.capacity)}</td><td>${escapeHtml(item.accessModes)}</td><td>${escapeHtml(item.storageClass)}</td><td>${escapeHtml(item.age)}</td>`,
  );
}

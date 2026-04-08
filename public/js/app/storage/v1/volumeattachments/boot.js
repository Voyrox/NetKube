let volumeAttachmentItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("volumeAttachmentsSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredVolumeAttachments(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/storage/volumeattachments");
      volumeAttachmentItems = data.items || [];
      applyPageMeta(data.meta);
      setText(
        "volumeAttachmentsHeroTitle",
        `Volume Attachments (${data.count || 0})`,
      );
      renderFilteredVolumeAttachments(searchInput?.value || "");
    } catch (error) {
      volumeAttachmentItems = [];
      renderVolumeAttachmentsTable(
        [],
        error.message || "Failed to load volume attachments",
      );
    }
  });
});
function renderFilteredVolumeAttachments(query) {
  const items = filterByQuery(volumeAttachmentItems, query, [
    "name",
    "attacher",
    "node",
    "persistentVolume",
    "attached",
  ]);
  renderVolumeAttachmentsTable(
    items,
    items.length ? "" : "No volume attachments match your search",
  );
  setText("volumeAttachmentsTableCount", `${items.length} rows`);
}
function renderVolumeAttachmentsTable(items, message) {
  renderTableRows(
    "volumeAttachmentsTableBody",
    6,
    items,
    message || "No volume attachments found",
    (item) =>
      `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.attacher)}</td><td>${escapeHtml(item.node)}</td><td>${escapeHtml(item.persistentVolume)}</td><td>${escapeHtml(item.attached)}</td><td>${escapeHtml(item.age)}</td>`,
  );
}

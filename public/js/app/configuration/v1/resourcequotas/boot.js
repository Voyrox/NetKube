let resourceQuotaItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("resourceQuotasSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredResourceQuotas(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/configuration/resourcequotas");
      resourceQuotaItems = data.items || [];
      applyPageMeta(data.meta);
      setText(
        "resourceQuotasHeroTitle",
        `Resource Quotas (${data.count || 0})`,
      );
      renderFilteredResourceQuotas(searchInput?.value || "");
    } catch (error) {
      resourceQuotaItems = [];
      renderResourceQuotasTable(
        [],
        error.message || "Failed to load resource quotas",
      );
    }
  });
});
function renderFilteredResourceQuotas(query) {
  const filteredItems = filterResourceQuotas(resourceQuotaItems, query);
  renderResourceQuotasTable(
    filteredItems,
    filteredItems.length ? "" : "No resource quotas match your search",
  );
  setText("resourceQuotasTableCount", `${filteredItems.length} rows`);
}
function filterResourceQuotas(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;
  return items.filter((item) =>
    [item.namespace, item.name].some((value) =>
      String(value || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}
function renderResourceQuotasTable(items, message) {
  const body = document.getElementById("resourceQuotasTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="6">${escapeHtml(message || "No resource quotas found")}</td></tr>`;
    return;
  }
  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `<td>${escapeHtml(item.namespace)}</td><td>${escapeHtml(item.name)}</td><td>${escapeHtml(String(item.scopes))}</td><td>${escapeHtml(String(item.hard))}</td><td>${escapeHtml(String(item.used))}</td><td>${escapeHtml(item.age)}</td>`;
    body.appendChild(row);
  });
}

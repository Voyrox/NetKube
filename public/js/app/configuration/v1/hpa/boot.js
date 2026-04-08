let hpaItems = [];
document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("hpaSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredHpa(searchInput.value),
  );
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/configuration/hpa");
      hpaItems = data.items || [];
      applyPageMeta(data.meta);
      setText("hpaHeroTitle", `HPA (${data.count || 0})`);
      renderFilteredHpa(searchInput?.value || "");
    } catch (error) {
      hpaItems = [];
      renderHpaTable([], error.message || "Failed to load HPA");
    }
  });
});
function renderFilteredHpa(query) {
  const filteredItems = filterHpa(hpaItems, query);
  renderHpaTable(
    filteredItems,
    filteredItems.length ? "" : "No HPA match your search",
  );
  setText("hpaTableCount", `${filteredItems.length} rows`);
}
function filterHpa(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;
  return items.filter((item) =>
    [item.namespace, item.name, item.target].some((value) =>
      String(value || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}
function renderHpaTable(items, message) {
  const body = document.getElementById("hpaTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="7">${escapeHtml(message || "No HPA found")}</td></tr>`;
    return;
  }
  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `<td>${escapeHtml(item.namespace)}</td><td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.target)}</td><td>${escapeHtml(item.min)}</td><td>${escapeHtml(String(item.max))}</td><td>${escapeHtml(item.current)}</td><td>${escapeHtml(item.age)}</td>`;
    body.appendChild(row);
  });
}

let limitRangeItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("limitRangesSearch");
  searchInput?.addEventListener("input", () => renderFilteredLimitRanges(searchInput.value));

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/configuration/limitranges");
      limitRangeItems = data.items || [];
      applyPageMeta(data.meta);
      setText("limitRangesHeroTitle", `Limit Ranges (${data.count || 0})`);
      renderFilteredLimitRanges(searchInput?.value || "");
    } catch (error) {
      limitRangeItems = [];
      renderLimitRangesTable([], error.message || "Failed to load limit ranges");
    }
  });
});

function renderFilteredLimitRanges(query) {
  const filteredItems = filterLimitRanges(limitRangeItems, query);
  renderLimitRangesTable(filteredItems, filteredItems.length ? "" : "No limit ranges match your search");
  setText("limitRangesTableCount", `${filteredItems.length} rows`);
}

function filterLimitRanges(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;
  return items.filter((item) => [item.namespace, item.name].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderLimitRangesTable(items, message) {
  const body = document.getElementById("limitRangesTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="4">${escapeHtml(message || "No limit ranges found")}</td></tr>`;
    return;
  }
  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `<td>${escapeHtml(item.namespace)}</td><td>${escapeHtml(item.name)}</td><td>${escapeHtml(String(item.limits))}</td><td>${escapeHtml(item.age)}</td>`;
    body.appendChild(row);
  });
}

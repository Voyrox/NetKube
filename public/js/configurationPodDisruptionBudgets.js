let pdbItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("pdbSearch");
  searchInput?.addEventListener("input", () => renderFilteredPdb(searchInput.value));

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/configuration/poddisruptionbudgets");
      pdbItems = data.items || [];
      applyPageMeta(data.meta);
      setText("pdbHeroTitle", `Pod Disruption Budgets (${data.count || 0})`);
      renderFilteredPdb(searchInput?.value || "");
    } catch (error) {
      pdbItems = [];
      renderPdbTable([], error.message || "Failed to load pod disruption budgets");
    }
  });
});

function renderFilteredPdb(query) {
  const filteredItems = filterPdb(pdbItems, query);
  renderPdbTable(filteredItems, filteredItems.length ? "" : "No pod disruption budgets match your search");
  setText("pdbTableCount", `${filteredItems.length} rows`);
}

function filterPdb(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;
  return items.filter((item) => [item.namespace, item.name].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderPdbTable(items, message) {
  const body = document.getElementById("pdbTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="6">${escapeHtml(message || "No pod disruption budgets found")}</td></tr>`;
    return;
  }
  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `<td>${escapeHtml(item.namespace)}</td><td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.minAvailable)}</td><td>${escapeHtml(item.maxUnavailable)}</td><td>${escapeHtml(String(item.allowed))}</td><td>${escapeHtml(item.age)}</td>`;
    body.appendChild(row);
  });
}

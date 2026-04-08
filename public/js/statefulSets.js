let statefulSetItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("statefulSetSearch");

  searchInput?.addEventListener("input", () => {
    renderFilteredStatefulSets(searchInput.value);
  });

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/statefulsets");
      statefulSetItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "statefulSetsNamespaceFilter" });
      setText("statefulSetsHeroTitle", `StatefulSets (${data.stats?.healthy || 0} healthy)`);
      renderFilteredStatefulSets(searchInput?.value || "");
    } catch (error) {
      statefulSetItems = [];
      renderStatefulSetsTable([], error.message || "Failed to load statefulsets");
    }
  });
});

function renderFilteredStatefulSets(query) {
  const filteredItems = filterStatefulSets(statefulSetItems, query);
  renderStatefulSetsTable(filteredItems, filteredItems.length ? "" : "No statefulsets match your search");
  setText("statefulSetsTableCount", `${filteredItems.length} rows`);
}

function filterStatefulSets(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;

  return items.filter((item) => [item.namespace, item.name, item.status].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderStatefulSetsTable(items, message) {
  const body = document.getElementById("statefulSetsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="8">${escapeHtml(message || "No statefulsets found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.ready)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(String(item.desired))}</td>
      <td>${escapeHtml(String(item.current))}</td>
      <td>${escapeHtml(String(item.updated))}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

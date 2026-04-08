let daemonSetItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("daemonSetSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredDaemonSets(searchInput.value),
  );

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/daemonsets");
      daemonSetItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "daemonSetsNamespaceFilter" });
      setText(
        "daemonSetsHeroTitle",
        `DaemonSets (${data.stats?.healthy || 0} healthy)`,
      );
      renderFilteredDaemonSets(searchInput?.value || "");
    } catch (error) {
      daemonSetItems = [];
      renderDaemonSetsTable([], error.message || "Failed to load daemonsets");
    }
  });
});

function renderFilteredDaemonSets(query) {
  const filteredItems = filterDaemonSets(daemonSetItems, query);
  renderDaemonSetsTable(
    filteredItems,
    filteredItems.length ? "" : "No daemonsets match your search",
  );
  setText("daemonSetsTableCount", `${filteredItems.length} rows`);
}

function filterDaemonSets(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;
  return items.filter((item) =>
    [item.namespace, item.name, item.status].some((value) =>
      String(value || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}

function renderDaemonSetsTable(items, message) {
  const body = document.getElementById("daemonSetsTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="8">${escapeHtml(message || "No daemonsets found")}</td></tr>`;
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
      <td>${escapeHtml(String(item.available))}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

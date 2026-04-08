let replicaSetItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("replicaSetSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredReplicaSets(searchInput.value),
  );

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/replicasets");
      replicaSetItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "replicaSetsNamespaceFilter" });
      setText(
        "replicaSetsHeroTitle",
        `ReplicaSets (${data.stats?.healthy || 0} healthy)`,
      );
      renderFilteredReplicaSets(searchInput?.value || "");
    } catch (error) {
      replicaSetItems = [];
      renderReplicaSetsTable([], error.message || "Failed to load replicasets");
    }
  });
});

function renderFilteredReplicaSets(query) {
  const filteredItems = filterReplicaSets(replicaSetItems, query);
  renderReplicaSetsTable(
    filteredItems,
    filteredItems.length ? "" : "No replicasets match your search",
  );
  setText("replicaSetsTableCount", `${filteredItems.length} rows`);
}

function filterReplicaSets(items, query) {
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

function renderReplicaSetsTable(items, message) {
  const body = document.getElementById("replicaSetsTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="7">${escapeHtml(message || "No replicasets found")}</td></tr>`;
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
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

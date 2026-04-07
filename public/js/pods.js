let podItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("deploymentSearch");

  searchInput?.addEventListener("input", () => {
    renderFilteredPods(searchInput.value);
  });

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/pods");
      podItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "podsNamespaceFilter" });
      setText("podsHeroTitle", `Pods (${data.stats?.running || 0} running)`);
      renderFilteredPods(searchInput?.value || "");
    } catch (error) {
      podItems = [];
      renderPodsTable([], error.message || "Failed to load pods");
    }
  });
});

function renderFilteredPods(query) {
  const filteredItems = filterPods(podItems, query);
  renderPodsTable(filteredItems, filteredItems.length ? "" : "No pods match your search");
  setText("podsTableCount", `${filteredItems.length} rows`);
}

function filterPods(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;

  return items.filter((item) => [item.namespace, item.name, item.status, item.node, item.podIP].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderPodsTable(items, message) {
  const body = document.getElementById("podsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="10">${escapeHtml(message || "No pods found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.ready)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(String(item.restarts))}</td>
      <td>${escapeHtml(item.lastRestart)}</td>
      <td>${escapeHtml(item.lastRestartReason)}</td>
      <td>${escapeHtml(item.node)}</td>
      <td>${escapeHtml(item.podIP)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    row.addEventListener("click", () => openPod(item));
    row.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openPod(item);
      }
    });
    body.appendChild(row);
  });
}

function openPod(item) {
  const params = new URLSearchParams({
    namespace: item.namespace || "",
    name: item.name || ""
  });

  window.location.href = `/workloads/manage/pod?${params.toString()}`;
}

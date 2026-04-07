document.addEventListener("DOMContentLoaded", async () => {
  const searchInput = document.getElementById("deploymentSearch");

  try {
    const data = await fetchClusterData("/api/workloads/pods");
    applyPageMeta(data.meta, { namespaceId: "podsNamespaceFilter" });
    setText("podsTableCount", `${data.count || 0} rows`);
    setText("podsHeroTitle", `Pods (${data.stats?.running || 0} running)`);
    renderPodsTable(data.items || []);

    if (searchInput) {
      searchInput.addEventListener("input", () => {
        const filteredItems = filterPods(data.items || [], searchInput.value);
        renderPodsTable(filteredItems, filteredItems.length ? "" : "No pods match your search");
        setText("podsTableCount", `${filteredItems.length} rows`);
      });
    }
  } catch (error) {
    renderPodsTable([], error.message || "Failed to load pods");
  }
});

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

document.addEventListener("DOMContentLoaded", () => {
  const params = new URLSearchParams(window.location.search);
  const name = params.get("name") || "";
  const query = name ? `?${new URLSearchParams({ name }).toString()}` : "";

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData(`/api/cluster/node${query}`);
      applyPageMeta(data.meta);
      renderNode(data.item || {});
    } catch (error) {
      renderNodeError(error.message || "Failed to load node details");
    }
  });
});

function renderNode(item) {
  setText("nodeHeroTitle", item.name || "Node");
  setText("nodeStatusMeta", item.status || "-");
  setText("nodeRoleMeta", item.role || "-");
  setText("nodeCpuPressure", item.cpuPressure || "-");
  setText("nodeMemoryPressure", item.memoryPressure || "-");
  setText("nodeDiskPressure", item.diskPressure || "-");
  setText("nodeKubeletVersion", item.kubeletVersion || "-");
  setText("nodeContainerRuntime", item.containerRuntime || "-");
  setText("nodeOsKernel", item.osKernel || "-");
  setText("nodeArchitecture", item.architecture || "-");
  setText("nodeInternalIP", item.internalIP || "-");
  setText("nodePodCIDR", item.podCIDR || "-");
  setText("nodeAllocatableCPU", item.allocatableCPU || "-");
  setText("nodeAllocatableMemory", item.allocatableMemory || "-");
  setText("nodeAllocatablePods", item.allocatablePods || "-");
  setText("nodeAllocatableStore", item.allocatableStore || "-");
  applyStatusChip("nodeLiveStatus", item.status || "Unknown");

  renderNodeLabels(item.labels || {});
  renderNodeConditions(item.conditions || []);
  renderNodeTimeline(item.timeline || []);
}

function renderNodeLabels(labels) {
  const element = document.getElementById("nodeLabelsList");
  if (!element) return;
  const entries = Object.entries(labels).sort((a, b) => a[0].localeCompare(b[0]));
  if (!entries.length) {
    element.innerHTML = '<li><span class="cluster-tag-key">No labels</span><strong class="cluster-tag-value">-</strong></li>';
    return;
  }
  element.innerHTML = entries.map(([key, value]) => `<li><span class="cluster-tag-key">${escapeHtml(key)}</span><strong class="cluster-tag-value">${escapeHtml(value || "-")}</strong></li>`).join("");
}

function renderNodeConditions(items) {
  const element = document.getElementById("nodeConditionsList");
  if (!element) return;
  if (!items.length) {
    element.innerHTML = '<p class="cluster-empty">No conditions reported for this node.</p>';
    return;
  }
  element.innerHTML = items.map((item) => `<article class="cluster-timeline-item"><div class="cluster-timeline-top"><span class="cluster-timeline-title">${escapeHtml(item.type || "Condition")}</span><span class="cluster-event-badge ${badgeClass(item.status)}">${escapeHtml(item.status || "Unknown")}</span></div><p class="cluster-timeline-message">${escapeHtml(item.reason || "-")}${item.message && item.message !== "-" ? ` - ${escapeHtml(item.message)}` : ""}</p></article>`).join("");
}

function renderNodeTimeline(items) {
  const element = document.getElementById("nodeTimelineList");
  if (!element) return;
  if (!items.length) {
    element.innerHTML = '<p class="cluster-empty">No recent node activity found.</p>';
    return;
  }
  element.innerHTML = items.map((item) => `<article class="cluster-timeline-item"><div class="cluster-timeline-top"><span class="cluster-timeline-title">${escapeHtml(item.title || "Event")}</span><span class="cluster-timeline-time">${escapeHtml(item.age || "-")}</span></div><p class="cluster-timeline-message">${escapeHtml(item.message || "-")}</p></article>`).join("");
}

function badgeClass(status) {
  const normalized = String(status || "").toLowerCase();
  if (normalized === "true" || normalized.includes("ready")) return "cluster-event-badge--healthy";
  if (normalized === "false") return "cluster-event-badge--warning";
  return "cluster-event-badge--danger";
}

function renderNodeError(message) {
  setText("nodeHeroTitle", "Node unavailable");
  setText("nodeStatusMeta", message);
  applyStatusChip("nodeLiveStatus", "Unavailable");
  ["nodeLabelsList", "nodeConditionsList", "nodeTimelineList"].forEach((id) => {
    const element = document.getElementById(id);
    if (element) {
      element.innerHTML = `<p class="cluster-empty">${escapeHtml(message)}</p>`;
    }
  });
}

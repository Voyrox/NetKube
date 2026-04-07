const CLUSTER_CONTEXT_HEADER = "X-NetKube-Context";

async function fetchClusterData(path, options = {}) {
  const contextId = window.NetKubeStorage?.getActiveContextId();
  if (!contextId) {
    throw new Error("No cluster selected");
  }

  const headers = {
    [CLUSTER_CONTEXT_HEADER]: contextId,
    ...(options.headers || {})
  };

  const response = await fetch(path, {
    ...options,
    headers
  });

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(data.error || "Failed to load cluster data");
  }

  return data;
}

function setText(id, value) {
  const element = document.getElementById(id);
  if (element) {
    element.textContent = value;
  }
}

function setWidth(id, value) {
  const element = document.getElementById(id);
  if (element) {
    element.style.width = value;
  }
}

function toggleSummaryDetail(prefix, section, label, value, options = {}) {
  const labelElement = document.getElementById(`${prefix}${section}Label`);
  const valueElement = document.getElementById(`${prefix}${section}`);
  const item = valueElement?.parentElement;
  const shouldShow = options.forceShow || Number(value || 0) > 0;

  if (labelElement) {
    labelElement.textContent = label;
  }

  if (valueElement) {
    valueElement.textContent = value || 0;
  }

  if (item) {
    item.hidden = !shouldShow;
  }

  return shouldShow;
}

function updateSummaryHealthBar(prefix, total, values) {
  const totalElement = document.getElementById(`${prefix}Total`);
  const card = totalElement?.closest(".summary-card");
  if (!card) return;

  const [healthySegment, warningSegment, dangerSegment] = card.querySelectorAll(".health-bar .segment");
  if (!healthySegment || !warningSegment || !dangerSegment) return;

  const safeTotal = Math.max(Number(total || 0), 0);
  if (safeTotal <= 0) {
    healthySegment.style.width = "0%";
    warningSegment.style.width = "0%";
    dangerSegment.style.width = "0%";
    return;
  }

  const primary = Math.max(Number(values.primary || 0), 0);
  const warning = Math.max(Number(values.warning || 0), 0);
  const danger = Math.max(Number(values.danger || 0), 0);
  const other = Math.max(Number(values.other || 0), 0);

  healthySegment.style.width = `${(primary / safeTotal) * 100}%`;
  warningSegment.style.width = `${(warning / safeTotal) * 100}%`;
  dangerSegment.style.width = `${((danger + other) / safeTotal) * 100}%`;
}

function updateSummaryCardVisibility(prefix, total, visibleItems, options = {}) {
  const totalElement = document.getElementById(`${prefix}Total`);
  const card = totalElement?.closest(".summary-card");
  if (!card) return;

  const breakdown = card.querySelector(".summary-breakdown");
  const healthBar = card.querySelector(".health-bar");
  const hasVisibleItems = visibleItems > 0;
  const shouldShowBreakdown = options.showBreakdown !== false && hasVisibleItems;

  if (breakdown) {
    breakdown.hidden = !shouldShowBreakdown;
  }

  if (healthBar) {
    healthBar.hidden = Number(total || 0) <= 0;
  }
}

function applyPageMeta(meta, options = {}) {
  if (!meta) return;

  setText(options.userId || "heroUser", meta.userName || meta.contextName || "Unknown user");
  setText(options.contextId || "selectedContext", meta.clusterName || meta.contextName || "Unknown cluster");
  setText(options.refreshId || "lastRefresh", formatRefresh(meta.lastRefresh));

  if (options.namespaceId) {
    setText(options.namespaceId, meta.namespace || "All namespaces");
  }
}

function formatRefresh(value) {
  if (!value) return "just now";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit"
  });
}

function applyStatusChip(id, status) {
  const element = document.getElementById(id);
  if (!element) return;

  element.textContent = status;
  element.classList.remove("status-chip--healthy", "status-chip--warning", "status-chip--danger");

  const normalized = String(status || "").toLowerCase();
  if (normalized.includes("healthy")) {
    element.classList.add("status-chip--healthy");
    return;
  }

  if (normalized.includes("watch") || normalized.includes("update")) {
    element.classList.add("status-chip--warning");
    return;
  }

  element.classList.add("status-chip--danger");
}

function renderWarnings(listId, emptyId, warnings) {
  const list = document.getElementById(listId);
  const empty = document.getElementById(emptyId);
  if (!list || !empty) return;

  list.innerHTML = "";

  if (!Array.isArray(warnings) || warnings.length === 0) {
    empty.textContent = "No warning events reported for the selected cluster.";
    return;
  }

  empty.textContent = "";
  warnings.forEach((warning) => {
    const article = document.createElement("article");
    article.className = "abnormality-item";
    article.innerHTML = `
      <div>
        <span class="abnormality-tag abnormality-tag--warning">${escapeHtml(warning.reason || "Warning")}</span>
        <h3>${escapeHtml(warning.message || warning.name || "Cluster warning")}</h3>
        <p>${escapeHtml(warning.name || "resource")} in namespace <strong>${escapeHtml(warning.namespace || "default")}</strong>.</p>
      </div>
      <span class="abnormality-value">${escapeHtml(warning.age || "-")}</span>
    `;
    list.appendChild(article);
  });
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function renderResourceUsage(data) {
  if (!data) return;

  renderResourceUsageMetric("usageCpu", data.usageCapacity?.cpu);
  renderResourceUsageMetric("usageMemory", data.usageCapacity?.memory);
  renderResourceUsageMetric("requestsPods", data.requestsAllocate?.pods);
  renderResourceUsageMetric("requestsCpu", data.requestsAllocate?.cpu);
  renderResourceUsageMetric("requestsMemory", data.requestsAllocate?.memory);
}

function renderResourceUsageMetric(prefix, metric) {
  if (!metric) return;

  const percent = Math.max(Math.min(Number(metric.percent || 0), 100), 0);
  setText(`${prefix}Percent`, `${Math.round(percent)}%`);
  setWidth(`${prefix}Fill`, `${percent}%`);
  setText(`${prefix}Value`, `${metric.used || "0"} / ${metric.total || "0"}`);
}

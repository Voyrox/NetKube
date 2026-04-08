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
  if (!card) {
    return;
  }

  const [healthySegment, warningSegment, dangerSegment] = card.querySelectorAll(
    ".health-bar .segment",
  );
  if (!healthySegment || !warningSegment || !dangerSegment) {
    return;
  }

  const safeTotal = Math.max(Number(values.healthTotal || total || 0), 0);
  if (safeTotal <= 0) {
    healthySegment.style.width = "0%";
    warningSegment.style.width = "0%";
    dangerSegment.style.width = "0%";
    return;
  }

  const primary = Math.max(Number(values.primary || 0), 0);
  const warning = Math.max(Number(values.warning || 0), 0);
  const danger = Math.max(Number(values.danger || 0), 0);
  healthySegment.style.width = `${(primary / safeTotal) * 100}%`;
  warningSegment.style.width = `${(warning / safeTotal) * 100}%`;
  dangerSegment.style.width = `${(danger / safeTotal) * 100}%`;
}

function updateSummaryCardVisibility(
  prefix,
  total,
  visibleItems,
  options = {},
) {
  const totalElement = document.getElementById(`${prefix}Total`);
  const card = totalElement?.closest(".summary-card");
  if (!card) {
    return;
  }

  const breakdown = card.querySelector(".summary-breakdown");
  const healthBar = card.querySelector(".health-bar");
  const hasVisibleItems = visibleItems > 0;
  const shouldShowBreakdown =
    options.showBreakdown !== false && hasVisibleItems;

  if (breakdown) {
    breakdown.hidden = !shouldShowBreakdown;
  }
  if (healthBar) {
    healthBar.hidden = Number(total || 0) <= 0;
  }
}

function renderWarnings(listId, emptyId, warnings) {
  const list = document.getElementById(listId);
  const empty = document.getElementById(emptyId);
  if (!list || !empty) {
    return;
  }

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

function renderResourceUsage(data) {
  if (!data) {
    return;
  }

  renderResourceUsageMetric("usageCpu", data.usageCapacity?.cpu);
  renderResourceUsageMetric("usageMemory", data.usageCapacity?.memory);
  renderResourceUsageMetric("requestsPods", data.requestsAllocate?.pods);
  renderResourceUsageMetric("requestsCpu", data.requestsAllocate?.cpu);
  renderResourceUsageMetric("requestsMemory", data.requestsAllocate?.memory);
}

function renderResourceUsageMetric(prefix, metric) {
  if (!metric) {
    return;
  }

  const percent = Math.max(Math.min(Number(metric.percent || 0), 100), 0);
  setText(`${prefix}Percent`, `${Math.round(percent)}%`);
  setWidth(`${prefix}Fill`, `${percent}%`);
  setText(`${prefix}Value`, `${metric.used || "0"} / ${metric.total || "0"}`);
}

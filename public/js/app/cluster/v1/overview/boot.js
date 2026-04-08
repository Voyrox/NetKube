document.addEventListener("DOMContentLoaded", () => {
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/cluster/overview");
      applyPageMeta(data.meta);

      renderMetricCard("nodes", data.nodes, [
        "Ready",
        "Pending",
        "Not ready",
        "Other",
      ]);
      renderMetricCard("pv", data.persistentVolumes, [
        "Bound",
        "Pending",
        "Failed",
        "Other",
      ]);
      renderMetricCard("crd", data.customResources, [
        "Established",
        "Pending",
        "Terminating",
        "Other",
      ]);
      renderResourceUsage(data.resourceUsage);
      renderWarnings(
        "clusterWarningsList",
        "clusterWarningsEmpty",
        data.warnings,
      );
    } catch (error) {
      setText(
        "clusterWarningsEmpty",
        error.message || "Failed to load cluster overview",
      );
    }
  });
});

function renderMetricCard(prefix, metric, labels) {
  if (!metric) return;
  const total = Math.max(Number(metric.total || 0), 0);
  const primary = Math.max(Number(metric.primary || 0), 0);
  const warning = Math.max(Number(metric.warning || 0), 0);
  const danger = Math.max(Number(metric.danger || 0), 0);
  const other = Math.max(total - primary - warning - danger, 0);
  const primaryOnly =
    total > 0 &&
    primary === total &&
    warning === 0 &&
    danger === 0 &&
    other === 0;

  applyStatusChip(`${prefix}Status`, metric.status || "Unknown");
  setText(`${prefix}Total`, total);

  let visibleItems = 0;
  visibleItems += toggleSummaryDetail(prefix, "Primary", labels[0], primary, {
    forceShow: !primaryOnly && total > 0 && primary > 0,
  });
  visibleItems += toggleSummaryDetail(prefix, "Warning", labels[1], warning);
  visibleItems += toggleSummaryDetail(prefix, "Danger", labels[2], danger);
  visibleItems += toggleSummaryDetail(prefix, "Other", labels[3], other);
  updateSummaryHealthBar(prefix, total, { primary, warning, danger, other });
  updateSummaryCardVisibility(prefix, total, visibleItems, {
    showBreakdown: !primaryOnly,
  });
}

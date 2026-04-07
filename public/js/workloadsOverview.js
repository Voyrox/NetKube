document.addEventListener("DOMContentLoaded", () => {
  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/overview");
      applyPageMeta(data.meta);

      renderWorkloadMetric("pods", data.pods, ["Running", "Pending", "Failed", "Other"]);
      renderWorkloadMetric("deployments", data.deployments, ["Healthy", "Warning", "Pending", "Other"], { hideStatus: true });
      renderWorkloadMetric("replicaSets", data.replicaSets, ["Ready", "Pending", "Issues", "Other"]);
      renderWorkloadMetric("daemonSets", data.daemonSets, ["Ready", "Pending", "Issues", "Other"]);
      renderWorkloadMetric("statefulSets", data.statefulSets, ["Ready", "Updating", "Issues", "Other"]);
      renderWorkloadMetric("cronJobs", data.cronJobs, ["Tracked", "Warning", "Issues", "Other"]);
      renderWorkloadMetric("jobs", data.jobs, ["Succeeded", "Active", "Failed", "Other"]);
      renderWorkloadMetric("resourceQuotas", data.resourceQuotas, ["Tracked", "Warning", "Issues", "Other"]);
      renderWarnings("workloadWarningsList", "workloadWarningsEmpty", data.warnings);
    } catch (error) {
      setText("workloadWarningsEmpty", error.message || "Failed to load workloads overview");
    }
  });
});

function renderWorkloadMetric(prefix, metric, labels, options = {}) {
  if (!metric) return;

  const total = Math.max(Number(metric.total || 0), 0);
  const primary = Math.max(Number(metric.primary || 0), 0);
  const warning = Math.max(Number(metric.warning || 0), 0);
  const danger = Math.max(Number(metric.danger || 0), 0);
  const other = Math.max(total - primary - warning - danger, 0);
  const primaryOnly = total > 0 && primary === total && warning === 0 && danger === 0 && other === 0;

  if (!options.hideStatus) {
    applyStatusChip(`${prefix}Status`, metric.status || "Unknown");
  }

  setText(`${prefix}Total`, total);

  let visibleItems = 0;
  visibleItems += toggleSummaryDetail(prefix, "Primary", labels[0], primary, { forceShow: !primaryOnly && total > 0 && primary > 0 });
  visibleItems += toggleSummaryDetail(prefix, "Warning", labels[1], warning);
  visibleItems += toggleSummaryDetail(prefix, "Danger", labels[2], danger);
  visibleItems += toggleSummaryDetail(prefix, "Other", labels[3], other);

  updateSummaryHealthBar(prefix, total, { primary, warning, danger, other });
  updateSummaryCardVisibility(prefix, total, visibleItems, { showBreakdown: !primaryOnly });
}

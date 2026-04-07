document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/workloads/overview");
    applyPageMeta(data.meta);

    renderWorkloadMetric("pods", data.pods, ["Running", "Pending", "Failed", "Other"]);
    renderWorkloadMetric("deployments", data.deployments, ["Healthy", "Warning", "Pending", "Other"]);
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

function renderWorkloadMetric(prefix, metric, labels) {
  if (!metric) return;

  applyStatusChip(`${prefix}Status`, metric.status || "Unknown");
  setText(`${prefix}Total`, metric.total || 0);
  setText(`${prefix}PrimaryLabel`, labels[0]);
  setText(`${prefix}Primary`, metric.primary || 0);
  setText(`${prefix}WarningLabel`, labels[1]);
  setText(`${prefix}Warning`, metric.warning || 0);
  setText(`${prefix}DangerLabel`, labels[2]);
  setText(`${prefix}Danger`, metric.danger || 0);
  setText(`${prefix}OtherLabel`, labels[3]);
  setText(`${prefix}Other`, Math.max((metric.total || 0) - (metric.primary || 0) - (metric.warning || 0) - (metric.danger || 0), 0));
}

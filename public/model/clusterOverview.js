document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/cluster/overview");
    applyPageMeta(data.meta);

    renderMetricCard("nodes", data.nodes, ["Ready", "Pending", "Not ready", "Other"]);
    renderMetricCard("pv", data.persistentVolumes, ["Bound", "Pending", "Failed", "Other"]);
    renderMetricCard("crd", data.customResources, ["Established", "Pending", "Terminating", "Other"]);
    renderWarnings("clusterWarningsList", "clusterWarningsEmpty", data.warnings);
  } catch (error) {
    setText("clusterWarningsEmpty", error.message || "Failed to load cluster overview");
  }
});

function renderMetricCard(prefix, metric, labels) {
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

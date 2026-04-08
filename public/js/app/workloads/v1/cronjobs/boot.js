let cronJobItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("cronJobsSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredCronJobs(searchInput.value),
  );

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/cronjobs");
      cronJobItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "cronJobsNamespaceFilter" });
      setText(
        "cronJobsHeroTitle",
        `CronJobs (${data.stats?.scheduled || 0} scheduled)`,
      );
      renderFilteredCronJobs(searchInput?.value || "");
    } catch (error) {
      cronJobItems = [];
      renderCronJobsTable([], error.message || "Failed to load cronjobs");
    }
  });
});

function renderFilteredCronJobs(query) {
  const filteredItems = filterCronJobs(cronJobItems, query);
  renderCronJobsTable(
    filteredItems,
    filteredItems.length ? "" : "No cronjobs match your search",
  );
  setText("cronJobsTableCount", `${filteredItems.length} rows`);
}

function filterCronJobs(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;
  return items.filter((item) =>
    [item.namespace, item.name, item.status, item.schedule].some((value) =>
      String(value || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}

function renderCronJobsTable(items, message) {
  const body = document.getElementById("cronJobsTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="8">${escapeHtml(message || "No cronjobs found")}</td></tr>`;
    return;
  }
  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.schedule)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(item.suspend)}</td>
      <td>${escapeHtml(String(item.active))}</td>
      <td>${escapeHtml(item.lastSchedule)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

let jobItems = [];

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("jobsSearch");
  searchInput?.addEventListener("input", () =>
    renderFilteredJobs(searchInput.value),
  );

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/jobs");
      jobItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "jobsNamespaceFilter" });
      setText(
        "jobsHeroTitle",
        `Jobs (${data.stats?.succeeded || 0} succeeded)`,
      );
      renderFilteredJobs(searchInput?.value || "");
    } catch (error) {
      jobItems = [];
      renderJobsTable([], error.message || "Failed to load jobs");
    }
  });
});

function renderFilteredJobs(query) {
  const filteredItems = filterJobs(jobItems, query);
  renderJobsTable(
    filteredItems,
    filteredItems.length ? "" : "No jobs match your search",
  );
  setText("jobsTableCount", `${filteredItems.length} rows`);
}

function filterJobs(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;
  return items.filter((item) =>
    [item.namespace, item.name, item.status].some((value) =>
      String(value || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}

function renderJobsTable(items, message) {
  const body = document.getElementById("jobsTableBody");
  if (!body) return;
  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="7">${escapeHtml(message || "No jobs found")}</td></tr>`;
    return;
  }
  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(item.completions)}</td>
      <td>${escapeHtml(String(item.active))}</td>
      <td>${escapeHtml(item.duration)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

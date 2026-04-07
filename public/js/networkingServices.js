document.addEventListener("DOMContentLoaded", async () => {
  try {
    const data = await fetchClusterData("/api/networking/services");
    applyPageMeta(data.meta);
    setText("servicesHeroTitle", `Services (${data.count || 0})`);
    setText("servicesTableCount", `${data.count || 0} rows`);
    renderServicesTable(data.items || []);
  } catch (error) {
    renderServicesTable([], error.message || "Failed to load services");
  }
});

function renderServicesTable(items, message) {
  const body = document.getElementById("servicesTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="7">${escapeHtml(message || "No services found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.type)}</td>
      <td>${escapeHtml(item.externalIP)}</td>
      <td>${escapeHtml(item.ports)}</td>
      <td>${escapeHtml(item.selector)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    body.appendChild(row);
  });
}

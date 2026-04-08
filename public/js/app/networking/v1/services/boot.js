document.addEventListener("DOMContentLoaded", () => {
  startAutoRefresh(async () => {
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
    row.innerHTML = `<td>${escapeHtml(item.namespace)}</td><td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.type)}</td><td>${escapeHtml(item.externalIP)}</td><td>${formatServicePorts(item.ports)}</td><td>${escapeHtml(item.selector)}</td><td>${escapeHtml(item.age)}</td>`;
    body.appendChild(row);
  });
}
function formatServicePorts(value) {
  const text = String(value || "").trim();
  if (!text || text === "-") return escapeHtml(text || "-");
  const entries = text.split(/\s*,\s*/).filter(Boolean);
  if (!entries.length) return escapeHtml(text);
  return entries.map((entry) => formatServicePortEntry(entry)).join("");
}
function formatServicePortEntry(entry) {
  const match = entry.match(
    /^([A-Z]+)\/(\S+?)(?:\s*->\s*(\S+))?(?:\s*\(([^)]+)\))?$/,
  );
  if (!match)
    return `<div class="service-port-entry">${escapeHtml(entry)}</div>`;
  const [, protocol, port, targetPort, name] = match;
  const target = targetPort || port;
  const label = name ? ` (${name})` : "";
  return `<div class="service-port-entry"><span class="service-port-entry__value">${escapeHtml(`${port}:${target}`)}</span><span class="service-port-entry__meta">/${escapeHtml(protocol)}${escapeHtml(label)}</span></div>`;
}

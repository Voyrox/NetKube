document.addEventListener("DOMContentLoaded", async () => {
  initializeLeaseDrawer();

  try {
    const data = await fetchClusterData("/api/cluster/leases");
    applyPageMeta(data.meta);
    setText("leasesHeroTitle", `Leases (${data.count || 0})`);
    setText("leasesTableCount", `${data.count || 0} rows`);
    renderLeasesTable(data.items || []);
  } catch (error) {
    renderLeasesTable([], error.message || "Failed to load leases");
  }
});

function renderLeasesTable(items, message) {
  const body = document.getElementById("leasesTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="5">${escapeHtml(message || "No leases found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.holder)}</td>
      <td>${escapeHtml(item.lastRenew)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    row.addEventListener("click", () => openLeaseYaml(item));
    row.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openLeaseYaml(item);
      }
    });
    body.appendChild(row);
  });
}

function initializeLeaseDrawer() {
  const close = document.getElementById("leaseYamlClose");
  const backdrop = document.getElementById("leaseYamlBackdrop");
  if (close) {
    close.addEventListener("click", closeLeaseDrawer);
  }
  if (backdrop) {
    backdrop.addEventListener("click", closeLeaseDrawer);
  }
}

async function openLeaseYaml(item) {
  const drawer = document.getElementById("leaseYamlDrawer");
  const backdrop = document.getElementById("leaseYamlBackdrop");
  const title = document.getElementById("leaseYamlTitle");
  const content = document.getElementById("leaseYamlContent");
  if (!drawer || !title || !content) return;

  drawer.hidden = false;
  if (backdrop) backdrop.hidden = false;
  title.textContent = item.name || "Lease";
  setDrawerText("leaseYamlContent", "Loading YAML...");

  try {
    const query = new URLSearchParams({ namespace: item.namespace || "", name: item.name || "" }).toString();
    const data = await fetchClusterData(`/api/cluster/lease/yaml?${query}`);
    renderYamlDrawerContent("leaseYamlContent", data.content || "No YAML available.");
  } catch (error) {
    setDrawerText("leaseYamlContent", error.message || "Failed to load lease YAML.");
  }
}

function closeLeaseDrawer() {
  const drawer = document.getElementById("leaseYamlDrawer");
  const backdrop = document.getElementById("leaseYamlBackdrop");
  if (drawer) {
    drawer.hidden = true;
  }
  if (backdrop) {
    backdrop.hidden = true;
  }
}

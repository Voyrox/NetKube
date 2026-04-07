document.addEventListener("DOMContentLoaded", () => {
  initializeNamespaceDrawer();

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/cluster/namespaces");
      applyPageMeta(data.meta);
      setText("namespacesHeroTitle", `Namespaces (${data.count || 0})`);
      setText("namespacesTableCount", `${data.count || 0} rows`);
      renderNamespacesTable(data.items || []);
    } catch (error) {
      renderNamespacesTable([], error.message || "Failed to load namespaces");
    }
  });
});

function renderNamespacesTable(items, message) {
  const body = document.getElementById("namespacesTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="3">${escapeHtml(message || "No namespaces found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.phase)}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    row.addEventListener("click", () => openNamespaceYaml(item));
    row.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openNamespaceYaml(item);
      }
    });
    body.appendChild(row);
  });
}

function initializeNamespaceDrawer() {
  const close = document.getElementById("namespaceYamlClose");
  const backdrop = document.getElementById("namespaceYamlBackdrop");
  if (close) {
    close.addEventListener("click", closeNamespaceDrawer);
  }
  if (backdrop) {
    backdrop.addEventListener("click", closeNamespaceDrawer);
  }
}

async function openNamespaceYaml(item) {
  const drawer = document.getElementById("namespaceYamlDrawer");
  const backdrop = document.getElementById("namespaceYamlBackdrop");
  const title = document.getElementById("namespaceYamlTitle");
  const content = document.getElementById("namespaceYamlContent");
  if (!drawer || !title || !content) return;

  drawer.hidden = false;
  if (backdrop) backdrop.hidden = false;
  title.textContent = item.name || "Namespace";
  setDrawerText("namespaceYamlContent", "Loading YAML...");

  try {
    const query = new URLSearchParams({ name: item.name || "" }).toString();
    const data = await fetchClusterData(`/api/cluster/namespace/yaml?${query}`);
    renderYamlDrawerContent("namespaceYamlContent", data.content || "No YAML available.");
  } catch (error) {
    setDrawerText("namespaceYamlContent", error.message || "Failed to load namespace YAML.");
  }
}

function closeNamespaceDrawer() {
  const drawer = document.getElementById("namespaceYamlDrawer");
  const backdrop = document.getElementById("namespaceYamlBackdrop");
  if (drawer) {
    drawer.hidden = true;
  }
  if (backdrop) {
    backdrop.hidden = true;
  }
}

let secretItems = [];

document.addEventListener("DOMContentLoaded", () => {
  initializeSecretDrawer();
  const searchInput = document.getElementById("secretsSearch");

  searchInput?.addEventListener("input", () => {
    renderFilteredSecrets(searchInput.value);
  });

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/configuration/secrets");
      secretItems = data.items || [];
      applyPageMeta(data.meta);
      setText("secretsHeroTitle", `Secrets (${data.count || 0})`);
      renderFilteredSecrets(searchInput?.value || "");
    } catch (error) {
      secretItems = [];
      renderSecretsTable([], error.message || "Failed to load secrets");
    }
  });
});

function renderFilteredSecrets(query) {
  const filteredItems = filterSecrets(secretItems, query);
  renderSecretsTable(filteredItems, filteredItems.length ? "" : "No secrets match your search");
  setText("secretsTableCount", `${filteredItems.length} rows`);
}

function filterSecrets(items, query) {
  const normalized = String(query || "").trim().toLowerCase();
  if (!normalized) return items;

  return items.filter((item) => [item.namespace, item.name, item.type].some((value) => String(value || "").toLowerCase().includes(normalized)));
}

function renderSecretsTable(items, message) {
  const body = document.getElementById("secretsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="5">${escapeHtml(message || "No secrets found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.type)}</td>
      <td>${escapeHtml(String(item.data))}</td>
      <td>${escapeHtml(item.age)}</td>
    `;
    row.addEventListener("click", () => openSecretData(item));
    row.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openSecretData(item);
      }
    });
    body.appendChild(row);
  });
}

function initializeSecretDrawer() {
  const close = document.getElementById("secretDataClose");
  const backdrop = document.getElementById("secretDataBackdrop");
  if (close) {
    close.addEventListener("click", closeSecretDrawer);
  }
  if (backdrop) {
    backdrop.addEventListener("click", closeSecretDrawer);
  }
}

async function openSecretData(item) {
  const drawer = document.getElementById("secretDataDrawer");
  const backdrop = document.getElementById("secretDataBackdrop");
  const title = document.getElementById("secretDataTitle");
  if (!drawer || !title) return;

  drawer.hidden = false;
  if (backdrop) backdrop.hidden = false;
  title.textContent = `${item.namespace || "default"} / ${item.name || "Secret"}`;
  setDrawerText("secretDataContent", "Loading secret data...");

  try {
    const query = new URLSearchParams({
      namespace: item.namespace || "",
      name: item.name || ""
    }).toString();
    const data = await fetchClusterData(`/api/configuration/secret?${query}`);
    renderYamlDrawerContent("secretDataContent", data.content || "No secret data available.");
  } catch (error) {
    setDrawerText("secretDataContent", error.message || "Failed to load secret data.");
  }
}

function closeSecretDrawer() {
  const drawer = document.getElementById("secretDataDrawer");
  const backdrop = document.getElementById("secretDataBackdrop");
  if (drawer) {
    drawer.hidden = true;
  }
  if (backdrop) {
    backdrop.hidden = true;
  }
}

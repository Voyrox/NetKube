const SOURCES_API = "/api/config/sources";
const SELECTED_CONTEXTS_API = "/api/config/selected-contexts";
const CONTEXTS_API = "/api/config/contexts";
const CLUSTER_OVERVIEW_PATH = "/clusters/overview";

document.addEventListener("DOMContentLoaded", async () => {
  const modal = document.getElementById("setupModal");
  const openButton = document.getElementById("openSetupModal");
  const closeButton = document.getElementById("closeSetupModal");
  const addFilesButton = document.getElementById("addFilesButton");
  const addFolderButton = document.getElementById("addFolderButton");
  const fileInput = document.getElementById("fileInput");
  const folderInput = document.getElementById("folderInput");
  const sourceList = document.getElementById("setupSourceList");
  const emptyState = document.getElementById("setupEmptyState");

  const contextList = document.getElementById("contextList");
  const contextHelperText = document.getElementById("contextHelperText");
  const connectButton = document.getElementById("connectButton");
  const searchInput = document.getElementById("contextSearchInput");
  const refreshContextsButton = document.getElementById("refreshContextsButton");

  let sources = [];
  let selectedContextIds = await normalizeSelectedContextIds(
    await loadSelectedContextIds()
  );

  await refreshSources();
  let parsedContexts = await loadContexts();

  renderSources();
  renderContexts();

  openButton?.addEventListener("click", () => {
    modal?.classList.add("open");
  });

  closeButton?.addEventListener("click", () => {
    modal?.classList.remove("open");
  });

  modal?.addEventListener("click", (event) => {
    if (event.target === modal) {
      modal.classList.remove("open");
    }
  });

  addFilesButton?.addEventListener("click", () => {
    fileInput?.click();
  });

  addFolderButton?.addEventListener("click", () => {
    folderInput?.click();
  });

  searchInput?.addEventListener("input", () => {
    renderContexts();
  });

  refreshContextsButton?.addEventListener("click", async () => {
    parsedContexts = await loadContexts();
    await pruneMissingSelections();
    renderContexts();
  });

  connectButton?.addEventListener("click", () => {
    const activeContextId = selectedContextIds[0];
    if (!activeContextId) return;

    window.NetKubeStorage?.setActiveContextId(activeContextId);
    window.location.href = CLUSTER_OVERVIEW_PATH;
  });

  fileInput?.addEventListener("change", async (event) => {
    const files = Array.from(event.target.files || []);
    if (!files.length) return;

    const items = await createSourceItems(files, "file");
    await mergeSources(items);

    fileInput.value = "";
  });

  folderInput?.addEventListener("change", async (event) => {
    const files = Array.from(event.target.files || []);
    if (!files.length) return;

    const items = await createSourceItems(files, "folder-file");
    await mergeSources(items);

    folderInput.value = "";
  });

  async function createSourceItems(files, type) {
    const textFiles = files.filter(isLikelyKubeconfigFile);
    const items = [];

    for (const file of textFiles) {
      const content = await file.text();

      items.push({
        id: createId((file.webkitRelativePath || file.name) + file.size + file.lastModified),
        type,
        name: file.name,
        size: file.size,
        lastModified: file.lastModified,
        content
      });
    }

    return items;
  }

  function isLikelyKubeconfigFile(file) {
    const path = (file.webkitRelativePath || file.name || "").toLowerCase();

    return (
      path.endsWith(".yaml") ||
      path.endsWith(".yml") ||
      path.endsWith(".conf") ||
      path.endsWith(".config") ||
      path.endsWith("config")
    );
  }

  async function mergeSources(newItems) {
    const existingIds = new Set(sources.map((item) => item.id));
    const merged = [...sources];

    for (const item of newItems) {
      const existingIndex = merged.findIndex((source) => source.id === item.id);
      if (existingIndex >= 0) {
        merged[existingIndex] = item;
      } else if (!existingIds.has(item.id)) {
        merged.push(item);
      }
    }

    await saveSources(merged);

    await refreshSources();
    parsedContexts = await loadContexts();
    await pruneMissingSelections();
    renderSources();
    renderContexts();
  }

  function renderSources() {
    if (!sourceList || !emptyState) return;

    sourceList.innerHTML = "";

    if (!sources.length) {
      emptyState.style.display = "block";
      return;
    }

    emptyState.style.display = "none";

    sources.forEach((item) => {
      const row = document.createElement("div");
      row.className = "setupSourceItem";

      const left = document.createElement("div");
      left.className = "setupSourceLeft";

      const icon = document.createElement("div");
      icon.className = "setupSourceIcon";

      const iconElement = document.createElement("i");
      iconElement.className =
        item.type === "file" ? "fa fa-file-text-o" : "fa fa-folder-open-o";

      icon.appendChild(iconElement);

      const textWrap = document.createElement("div");
      textWrap.className = "setupSourceText";

      const name = document.createElement("div");
      name.className = "setupSourceName";
      name.textContent = item.name;

      const meta = document.createElement("div");
      meta.className = "setupSourceMeta";
      meta.textContent = formatSourceMeta(item);

      textWrap.appendChild(name);
      textWrap.appendChild(meta);

      left.appendChild(icon);
      left.appendChild(textWrap);

      const deleteButton = document.createElement("button");
      deleteButton.className = "deleteSourceButton";
      deleteButton.type = "button";
      deleteButton.setAttribute("aria-label", `Delete ${item.name}`);
      deleteButton.innerHTML = '<i class="fa fa-trash"></i>';

      deleteButton.addEventListener("click", async () => {
        const updated = sources.filter((source) => source.id !== item.id);
        await saveSources(updated);

        await refreshSources();
        parsedContexts = await loadContexts();
        await pruneMissingSelections();
        renderSources();
        renderContexts();
      });

      row.appendChild(left);
      row.appendChild(deleteButton);

      sourceList.appendChild(row);
    });
  }

  function renderContexts() {
    if (!contextList || !contextHelperText || !connectButton) return;

    contextList.innerHTML = "";

    const query = (searchInput?.value || "").trim().toLowerCase();
    const filtered = parsedContexts.filter((item) => {
      const haystack = [
        item.contextName,
        item.userName,
        item.clusterName,
        item.namespace,
        item.sourceName,
        item.server
      ]
        .filter(Boolean)
        .join(" ")
        .toLowerCase();

      return haystack.includes(query);
    });

    if (!parsedContexts.length) {
      contextHelperText.textContent =
        "No kubeconfig contexts loaded yet. Open the setup modal from the cog to add files.";
      contextList.innerHTML = `
        <div class="contextItem empty">
          <div class="text">
            <a class="top">No contexts available</a>
            <a class="bottom">Add one or more kubeconfig files to get started</a>
          </div>
        </div>
      `;
      connectButton.disabled = true;
      return;
    }

    contextHelperText.textContent = `${filtered.length} context${filtered.length === 1 ? "" : "s"} found`;

    if (!filtered.length) {
      contextList.innerHTML = `
        <div class="contextItem empty">
          <div class="text">
            <a class="top">No matching contexts</a>
            <a class="bottom">Try a different search term</a>
          </div>
        </div>
      `;
      updateConnectButton();
      return;
    }

    filtered.forEach((item) => {
      const row = document.createElement("div");
      row.className = "contextItem";
      row.dataset.id = item.id;

      if (selectedContextIds.includes(item.id)) {
        row.classList.add("selected");
      }

      row.innerHTML = `
        <div class="icon">
          <img src="./public/logo/kubernetes.svg" alt="Kubernetes" />
        </div>

        <div class="text">
          <a class="top">${escapeHtml(item.userName || item.contextName)}</a>
          <a class="bottom">${escapeHtml(item.clusterName || item.contextName)}</a>
          <a class="meta">${escapeHtml(item.clusterName ? `Server: ${item.server}` : "no cluster")}</a>
        </div>
      `;

      row.addEventListener("click", async () => {
        await toggleContextSelection(item.id);
        row.classList.toggle("selected");
      });

      contextList.appendChild(row);
    });

    updateConnectButton();
  }

  async function toggleContextSelection(id) {
    if (selectedContextIds.includes(id)) {
      selectedContextIds = selectedContextIds.filter((item) => item !== id);
    } else {
      selectedContextIds = [id];
    }

    await saveSelectedContextIds();
    renderContexts();
  }

  function updateConnectButton() {
    if (connectButton) {
      connectButton.disabled = selectedContextIds.length === 0;
    }
  }

  async function pruneMissingSelections() {
    const validIds = new Set(parsedContexts.map((item) => item.id));
    const pruned = selectedContextIds.filter((id) => validIds.has(id));

    if (pruned.length !== selectedContextIds.length) {
      selectedContextIds = pruned;
      await saveSelectedContextIds();
    }
  }

  async function loadContexts() {
    try {
      const response = await fetch(CONTEXTS_API);
      if (!response.ok) return [];

      const data = await response.json();
      return Array.isArray(data.contexts) ? data.contexts : [];
    } catch (error) {
      console.error("Failed to load contexts from server", error);
      return [];
    }
  }

  async function saveSources(nextSources) {
    const response = await fetch(SOURCES_API, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ sources: nextSources })
    });

    if (!response.ok) {
      throw new Error("Failed to save sources");
    }
  }

  async function loadSources() {
    try {
      const response = await fetch(SOURCES_API);
      if (!response.ok) return [];

      const data = await response.json();
      return Array.isArray(data.sources) ? data.sources : [];
    } catch (error) {
      console.error("Failed to load kubeconfig sources from server", error);
      return [];
    }
  }

  async function refreshSources() {
    const loaded = await loadSources();
    sources.length = 0;
    sources.push(...loaded);
  }

  async function saveSelectedContextIds() {
    const response = await fetch(SELECTED_CONTEXTS_API, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ selectedContextIds })
    });

    if (!response.ok) {
      throw new Error("Failed to save selected context ids");
    }
  }

  async function loadSelectedContextIds() {
    try {
      const response = await fetch(SELECTED_CONTEXTS_API);
      if (!response.ok) return [];

      const data = await response.json();
      return Array.isArray(data.selectedContextIds) ? data.selectedContextIds : [];
    } catch (error) {
      console.error("Failed to load selected context ids from server", error);
      return [];
    }
  }

  async function normalizeSelectedContextIds(ids) {
    const normalized = Array.isArray(ids) && ids.length ? [ids[0]] : [];

    if (JSON.stringify(normalized) !== JSON.stringify(ids || [])) {
      selectedContextIds = normalized;
      await saveSelectedContextIds();
    }

    return normalized;
  }

  function formatSourceMeta(item) {
    const parts = [];

    if (typeof item.size === "number" && item.size > 0) {
      parts.push(formatBytes(item.size));
    }

    if (typeof item.lastModified === "number" && item.lastModified > 0) {
      parts.push(new Date(item.lastModified).toLocaleString());
    }

    return parts.join(" • ") || "Saved locally";
  }

  function formatBytes(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }

  function createId(value) {
    return btoa(unescape(encodeURIComponent(value))).replace(/=/g, "");
  }

  function escapeHtml(value) {
    return String(value)
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }
});

let podLogPoller = null;

document.addEventListener("DOMContentLoaded", async () => {
  const params = new URLSearchParams(window.location.search);
  const namespace = params.get("namespace") || "";
  const name = params.get("name") || "";

  initializePodTabs();
  bindCopyYAML();
  bindDeletePod(namespace, name);

  if (!namespace || !name) {
    renderPodError("Missing pod name or namespace in the URL.");
    return;
  }

  const requestQuery = new URLSearchParams({ namespace, name }).toString();

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData(`/api/workloads/pod?${requestQuery}`);
      applyPageMeta(data.meta);
      renderPodDetail(data.item || {});

      const containerNames = (data.item?.containers || [])
        .map((item) => item.name)
        .filter(Boolean);
      const selectedContainer = getSelectedLogContainer();
      const preferredContainer = containerNames.includes(selectedContainer)
        ? selectedContainer
        : containerNames[0] || "";

      setupLogSelectors(containerNames, (container) =>
        loadPodLogs(requestQuery, container),
      );
      await Promise.all([
        loadPodLogs(requestQuery, preferredContainer),
        loadPodEvents(requestQuery),
        loadPodYAML(requestQuery),
      ]);
    } catch (error) {
      renderPodError(error.message || "Failed to load pod details");
    }
  });
});

window.addEventListener("beforeunload", () => {
  if (podLogPoller) {
    window.clearInterval(podLogPoller);
  }
});

function initializePodTabs() {
  const tabs = document.querySelectorAll("[data-pod-tab]");
  const panels = document.querySelectorAll("[data-pod-panel]");
  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const selectedTab = tab.getAttribute("data-pod-tab");
      tabs.forEach((item) =>
        item.classList.toggle("pod-tab--active", item === tab),
      );
      panels.forEach((panel) =>
        panel.classList.toggle(
          "pod-tab-panel--active",
          panel.getAttribute("data-pod-panel") === selectedTab,
        ),
      );
    });
  });
}

function renderPodDetail(item) {
  setText("podHeroTitle", item.name || "Pod");
  setText("podNamespace", item.namespace || "-");
  setText("podName", item.name || "-");
  setText("podReady", item.ready || "-");
  setText("podRestarts", item.restarts || 0);
  setText("podNode", item.node || "-");
  setText("podAge", item.age || "-");
  setText("podPhase", item.phase || "-");
  setText("podServiceAccount", item.serviceAccount || "-");
  setText("podQosClass", item.qosClass || "-");
  setText("podPodIP", item.podIP || "-");
  setText("podHostIP", item.hostIP || "-");
  setText("podStartTime", item.startTime || "-");
  setText("podLastRestart", item.lastRestart || "-");
  setText("podLastRestartReason", item.lastRestartReason || "-");
  applyPodStatus("podStatus", item.status || item.phase || "Unknown");

  renderContainers(item.containers || []);
  renderConditions(item.conditions || []);
  renderKeyValueList(
    "podLabelsList",
    item.labels || {},
    "No labels attached to this pod.",
  );
  renderKeyValueList(
    "podAnnotationsList",
    item.annotations || {},
    "No annotations attached to this pod.",
  );
}

function applyPodStatus(id, status) {
  const element = document.getElementById(id);
  if (!element) return;
  element.textContent = status;
  element.classList.remove(
    "status-chip--healthy",
    "status-chip--warning",
    "status-chip--danger",
  );

  const normalized = String(status || "").toLowerCase();
  if (
    normalized.includes("running") ||
    normalized.includes("healthy") ||
    normalized.includes("succeeded")
  ) {
    element.classList.add("status-chip--healthy");
    return;
  }
  if (
    normalized.includes("pending") ||
    normalized.includes("creating") ||
    normalized.includes("init") ||
    normalized.includes("containercreating")
  ) {
    element.classList.add("status-chip--warning");
    return;
  }
  element.classList.add("status-chip--danger");
}

function renderContainers(items) {
  setText(
    "podContainersCount",
    `${items.length} container${items.length === 1 ? "" : "s"}`,
  );
  const containerList = document.getElementById("podContainersList");
  if (!containerList) return;
  if (!items.length) {
    containerList.innerHTML =
      '<p class="pod-empty">No containers reported for this pod.</p>';
    return;
  }
  containerList.innerHTML = items
    .map(
      (item) => `
    <article class="pod-collection-item">
      <div class="pod-collection-item__header">
        <strong>${escapeHtml(item.name || "container")}</strong>
        <span class="status-chip ${item.ready ? "status-chip--healthy" : "status-chip--warning"}">${item.ready ? "Ready" : "Not ready"}</span>
      </div>
      <div class="pod-collection-item__body">
        <div><span>Image</span><strong>${escapeHtml(item.image || "-")}</strong></div>
        <div><span>State</span><strong>${escapeHtml(item.state || "-")}</strong></div>
        <div><span>Restarts</span><strong>${escapeHtml(String(item.restarts || 0))}</strong></div>
      </div>
    </article>
  `,
    )
    .join("");
}

function renderConditions(items) {
  const conditionList = document.getElementById("podConditionsList");
  if (!conditionList) return;
  if (!items.length) {
    conditionList.innerHTML =
      '<p class="pod-empty">No conditions reported for this pod.</p>';
    return;
  }
  conditionList.innerHTML = items
    .map(
      (item) => `
    <article class="pod-collection-item">
      <div class="pod-collection-item__header">
        <strong>${escapeHtml(item.type || "Condition")}</strong>
        <span class="status-chip ${String(item.status || "").toLowerCase() === "true" ? "status-chip--healthy" : "status-chip--warning"}">${escapeHtml(item.status || "Unknown")}</span>
      </div>
      <div class="pod-collection-item__body">
        <div><span>Reason</span><strong>${escapeHtml(item.reason || "-")}</strong></div>
        <div><span>Message</span><strong>${escapeHtml(item.message || "-")}</strong></div>
      </div>
    </article>
  `,
    )
    .join("");
}

function renderKeyValueList(id, values, emptyMessage) {
  const element = document.getElementById(id);
  if (!element) return;
  const entries = Object.entries(values).sort((left, right) =>
    left[0].localeCompare(right[0]),
  );
  if (!entries.length) {
    element.innerHTML = `<p class="pod-empty">${escapeHtml(emptyMessage)}</p>`;
    return;
  }
  element.innerHTML = entries
    .map(
      ([key, value]) => `
    <div class="pod-keyvalue-item">
      <span>${escapeHtml(key)}</span>
      <strong>${escapeHtml(value || "-")}</strong>
    </div>
  `,
    )
    .join("");
}

function setupLogSelectors(containerNames, onChange) {
  const selectors = [
    document.getElementById("podLogContainerSelect"),
    document.getElementById("podLogsContainerSelect"),
  ];
  selectors.forEach((select) => {
    if (!select) return;
    select.innerHTML = containerNames.length
      ? containerNames
          .map(
            (name) =>
              `<option value="${escapeHtml(name)}">${escapeHtml(name)}</option>`,
          )
          .join("")
      : '<option value="">No containers</option>';
    select.onchange = () => {
      selectors.forEach((other) => {
        if (other && other !== select) {
          other.value = select.value;
        }
      });
      onChange(select.value);
    };
  });
}

function getSelectedLogContainer() {
  return (
    document.getElementById("podLogContainerSelect")?.value ||
    document.getElementById("podLogsContainerSelect")?.value ||
    ""
  );
}

async function loadPodLogs(requestQuery, container) {
  if (podLogPoller) {
    window.clearInterval(podLogPoller);
  }
  const run = async () => {
    try {
      const query = new URLSearchParams(requestQuery);
      if (container) query.set("container", container);
      const data = await fetchClusterData(
        `/api/workloads/pod/logs?${query.toString()}`,
      );
      const content =
        stripAnsiSequences(data.content || "") ||
        "No logs available for this container.";
      setText(
        "podLogsStatus",
        data.container ? `Live - ${data.container}` : "Live",
      );
      setText(
        "podLogsStatusFull",
        data.container ? `Live - ${data.container}` : "Live",
      );
      setCodeText("podOverviewLogs", content);
      setCodeText("podLogsFull", content);
      syncLogSelectors(data.container || container || "");
    } catch (error) {
      setText("podLogsStatus", "Unavailable");
      setText("podLogsStatusFull", "Unavailable");
      setCodeText("podOverviewLogs", error.message || "Failed to load logs");
      setCodeText("podLogsFull", error.message || "Failed to load logs");
    }
  };
  await run();
  podLogPoller = window.setInterval(run, 5000);
}

function syncLogSelectors(container) {
  [
    document.getElementById("podLogContainerSelect"),
    document.getElementById("podLogsContainerSelect"),
  ].forEach((select) => {
    if (select && container) {
      select.value = container;
    }
  });
}

async function loadPodEvents(requestQuery) {
  const element = document.getElementById("podEventsList");
  if (!element) return;
  try {
    const data = await fetchClusterData(
      `/api/workloads/pod/events?${requestQuery}`,
    );
    const items = data.items || [];
    if (!items.length) {
      element.innerHTML =
        '<p class="pod-empty">No events reported for this pod.</p>';
      return;
    }
    element.innerHTML = items
      .map(
        (item) => `
      <article class="pod-event-item">
        <div class="pod-event-item__meta">
          <span class="status-chip ${String(item.type || "").toLowerCase() === "warning" ? "status-chip--warning" : "status-chip--healthy"}">${escapeHtml(item.type || "Normal")}</span>
          <strong>${escapeHtml(item.reason || "-")}</strong>
          <span>${escapeHtml(item.age || "-")}</span>
        </div>
        <p>${escapeHtml(item.message || "-")}</p>
      </article>
    `,
      )
      .join("");
  } catch (error) {
    element.innerHTML = `<p class="pod-empty">${escapeHtml(error.message || "Failed to load events")}</p>`;
  }
}

async function loadPodYAML(requestQuery) {
  try {
    const data = await fetchClusterData(
      `/api/workloads/pod/yaml?${requestQuery}`,
    );
    renderYAMLContent("podYamlContent", data.content || "No YAML available.");
  } catch (error) {
    renderYAMLContent("podYamlContent", error.message || "Failed to load YAML");
  }
}

function bindCopyYAML() {
  const button = document.getElementById("podYamlCopyButton");
  const yaml = document.getElementById("podYamlContent");
  if (!button || !yaml) return;
  button.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(yaml.textContent || "");
      button.textContent = "Copied";
      window.setTimeout(() => {
        button.textContent = "Copy YAML";
      }, 1500);
    } catch {
      button.textContent = "Copy failed";
      window.setTimeout(() => {
        button.textContent = "Copy YAML";
      }, 1500);
    }
  });
}

function bindDeletePod(namespace, name) {
  const button = document.getElementById("deletePodButton");
  if (!button) return;

  button.addEventListener("click", async () => {
    if (!namespace || !name) return;

    const confirmed = window.confirm(
      `Delete pod ${name} from namespace ${namespace}?`,
    );
    if (!confirmed) return;

    const defaultLabel = button.textContent;
    button.disabled = true;
    button.textContent = "Deleting...";

    try {
      await fetchClusterData(
        `/api/workloads/pod?${new URLSearchParams({ namespace, name }).toString()}`,
        { method: "DELETE" },
      );
      window.location.href = "/workloads/pods";
    } catch (error) {
      window.alert(error.message || "Failed to delete pod");
      button.disabled = false;
      button.textContent = defaultLabel;
    }
  });
}

function setCodeText(id, value) {
  const element = document.getElementById(id);
  if (element) {
    element.textContent = value;
  }
}
function renderYAMLContent(id, value) {
  const element = document.getElementById(id);
  if (!element) return;
  element.innerHTML = window.NetKubeYaml?.renderHighlightedYaml(value) || "";
}
function stripAnsiSequences(value) {
  return String(value || "")
    .replace(/\u001b\[[0-9;]*[A-Za-z]/g, "")
    .replace(/\u001b\][^\u0007]*(\u0007|\u001b\\)/g, "")
    .trim();
}
function renderPodError(message) {
  setText("podHeroTitle", "Pod unavailable");
  setText("podName", message);
  applyPodStatus("podStatus", "Unavailable");
  setCodeText("podOverviewLogs", message);
  setCodeText("podLogsFull", message);
  renderYAMLContent("podYamlContent", message);
  [
    "podContainersList",
    "podConditionsList",
    "podLabelsList",
    "podAnnotationsList",
    "podEventsList",
  ].forEach((id) => {
    const element = document.getElementById(id);
    if (element) {
      element.innerHTML = `<p class="pod-empty">${escapeHtml(message)}</p>`;
    }
  });
}

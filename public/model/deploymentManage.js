document.addEventListener("DOMContentLoaded", async () => {
  const params = new URLSearchParams(window.location.search);
  const namespace = params.get("namespace") || "";
  const name = params.get("name") || "";

  initializeDeploymentTabs();

  if (!namespace || !name) {
    renderDeploymentError("Missing deployment name or namespace in the URL.");
    return;
  }

  const requestQuery = new URLSearchParams({ namespace, name }).toString();

  try {
    const data = await fetchClusterData(`/api/workloads/deployment?${requestQuery}`);
    applyPageMeta(data.meta);
    renderDeploymentDetail(data.item || {});

    await Promise.all([
      loadDeploymentEvents(requestQuery),
      loadDeploymentYAML(requestQuery)
    ]);
    bindCopyYAML("deploymentYamlCopyButton", "deploymentYamlContent");
  } catch (error) {
    renderDeploymentError(error.message || "Failed to load deployment details");
  }
});

function initializeDeploymentTabs() {
  const tabs = document.querySelectorAll("[data-deployment-tab]");
  const panels = document.querySelectorAll("[data-deployment-panel]");

  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const selected = tab.getAttribute("data-deployment-tab");
      tabs.forEach((item) => item.classList.toggle("pod-tab--active", item === tab));
      panels.forEach((panel) => panel.classList.toggle("pod-tab-panel--active", panel.getAttribute("data-deployment-panel") === selected));
    });
  });
}

function renderDeploymentDetail(item) {
  setText("deploymentHeroTitle", item.name || "Deployment");
  setText("deploymentNamespace", item.namespace || "-");
  setText("deploymentName", item.name || "-");
  setText("deploymentReady", item.ready || "-");
  setText("deploymentDesired", item.desired || 0);
  setText("deploymentAvailable", item.available || 0);
  setText("deploymentAge", item.age || "-");
  setText("deploymentStatusDetail", item.status || "-");
  setText("deploymentStrategy", item.strategy || "-");
  setText("deploymentSelector", item.selector || "-");
  setText("deploymentUpdated", item.updated || 0);
  setText("deploymentAvailableDetail", item.available || 0);
  setText("deploymentUnavailable", item.unavailable || 0);
  applyDeploymentStatus("deploymentStatus", item.status || "Unknown");

  renderDeploymentConditions(item.conditions || []);
  renderReplicaSets(item.replicaSets || []);
  renderPods(item.pods || []);
  renderTagList("deploymentLabelsList", item.labels || {}, "No labels attached to this deployment.");
  renderAnnotationList("deploymentAnnotationsList", item.annotations || {}, "No annotations attached to this deployment.");
}

function applyDeploymentStatus(id, status) {
  const element = document.getElementById(id);
  if (!element) return;

  element.textContent = status;
  element.classList.remove("status-chip--healthy", "status-chip--warning", "status-chip--danger");

  const normalized = String(status || "").toLowerCase();
  if (normalized.includes("healthy")) {
    element.classList.add("status-chip--healthy");
    return;
  }

  if (normalized.includes("pending") || normalized.includes("updating")) {
    element.classList.add("status-chip--warning");
    return;
  }

  element.classList.add("status-chip--danger");
}

function renderDeploymentConditions(items) {
  const list = document.getElementById("deploymentConditionsList");
  if (!list) return;

  if (!items.length) {
    list.innerHTML = '<p class="pod-empty">No conditions reported for this deployment.</p>';
    return;
  }

  list.innerHTML = items.map((item) => `
    <article class="deployment-condition-item">
      <div class="deployment-condition-item__header">
        <strong>${escapeHtml(item.type || "Condition")}</strong>
        <span class="status-chip ${String(item.status || "").toLowerCase() === "true" ? "status-chip--healthy" : "status-chip--warning"}">${escapeHtml(item.status || "Unknown")}</span>
      </div>
      <p>${escapeHtml(item.reason || "-")}${item.message && item.message !== "-" ? ` - ${escapeHtml(item.message)}` : ""}</p>
    </article>
  `).join("");
}

function renderReplicaSets(items) {
  setText("deploymentReplicaSetsCount", `${items.length} replicasets`);
  const list = document.getElementById("deploymentReplicaSetsList");
  if (!list) return;

  if (!items.length) {
    list.innerHTML = '<p class="pod-empty">No replica sets matched this deployment.</p>';
    return;
  }

  list.innerHTML = items.map((item) => `
    <article class="deployment-runtime-item">
      <strong>${escapeHtml(item.name || "replicaset")}</strong>
      <div class="deployment-runtime-item__meta">
        <span>Ready ${escapeHtml(item.ready || "-")}</span>
        <span>Desired ${escapeHtml(String(item.desired || 0))}</span>
        <span>${escapeHtml(item.age || "-")}</span>
      </div>
    </article>
  `).join("");
}

function renderPods(items) {
  setText("deploymentPodsCount", `${items.length} pods`);
  const list = document.getElementById("deploymentPodsList");
  if (!list) return;

  if (!items.length) {
    list.innerHTML = '<p class="pod-empty">No pods matched this deployment.</p>';
    return;
  }

  list.innerHTML = items.map((item) => `
    <article class="deployment-runtime-item">
      <div class="deployment-runtime-item__header">
        <strong>${escapeHtml(item.name || "pod")}</strong>
        <span class="status-chip ${podStatusClass(item.status)}">${escapeHtml(item.status || "Unknown")}</span>
      </div>
      <div class="deployment-runtime-item__meta">
        <span>Ready ${escapeHtml(item.ready || "-")}</span>
        <span>${escapeHtml(item.node || "-")}</span>
        <span>${escapeHtml(item.age || "-")}</span>
      </div>
    </article>
  `).join("");
}

function podStatusClass(status) {
  const normalized = String(status || "").toLowerCase();
  if (normalized.includes("running")) return "status-chip--healthy";
  if (normalized.includes("pending") || normalized.includes("creating")) return "status-chip--warning";
  return "status-chip--danger";
}

async function loadDeploymentEvents(requestQuery) {
  const list = document.getElementById("deploymentEventsList");
  if (!list) return;

  try {
    const data = await fetchClusterData(`/api/workloads/deployment/events?${requestQuery}`);
    const items = data.items || [];
    if (!items.length) {
      list.innerHTML = '<p class="pod-empty">No events reported for this deployment.</p>';
      return;
    }

    list.innerHTML = items.map((item) => `
      <article class="pod-event-item">
        <div class="pod-event-item__meta">
          <span class="status-chip ${String(item.type || "").toLowerCase() === "warning" ? "status-chip--warning" : "status-chip--healthy"}">${escapeHtml(item.type || "Normal")}</span>
          <strong>${escapeHtml(item.reason || "-")}</strong>
          <span>${escapeHtml(item.age || "-")}</span>
        </div>
        <p>${escapeHtml(item.message || "-")}</p>
      </article>
    `).join("");
  } catch (error) {
    list.innerHTML = `<p class="pod-empty">${escapeHtml(error.message || "Failed to load events")}</p>`;
  }
}

async function loadDeploymentYAML(requestQuery) {
  try {
    const data = await fetchClusterData(`/api/workloads/deployment/yaml?${requestQuery}`);
    renderYAMLContent("deploymentYamlContent", data.content || "No YAML available.");
  } catch (error) {
    renderYAMLContent("deploymentYamlContent", error.message || "Failed to load YAML");
  }
}

function renderTagList(id, values, emptyMessage) {
  const element = document.getElementById(id);
  if (!element) return;

  const entries = Object.entries(values).sort((left, right) => left[0].localeCompare(right[0]));
  if (!entries.length) {
    element.innerHTML = `<p class="pod-empty">${escapeHtml(emptyMessage)}</p>`;
    return;
  }

  element.innerHTML = entries.map(([key, value]) => `
    <div class="deployment-tag-item">
      <span>${escapeHtml(key)}</span>
      <strong>${escapeHtml(value || "-")}</strong>
    </div>
  `).join("");
}

function renderAnnotationList(id, values, emptyMessage) {
  const element = document.getElementById(id);
  if (!element) return;

  const entries = Object.entries(values)
    .filter(([key]) => key !== "kubectl.kubernetes.io/last-applied-configuration")
    .sort((left, right) => left[0].localeCompare(right[0]));

  if (!entries.length) {
    element.innerHTML = `<p class="pod-empty">${escapeHtml(emptyMessage)}</p>`;
    return;
  }

  element.innerHTML = entries.map(([key, value]) => `
    <div class="deployment-annotation-item">
      <span>${escapeHtml(key)}</span>
      <strong>${escapeHtml(value || "-")}</strong>
    </div>
  `).join("");
}

function bindCopyYAML(buttonId, contentId) {
  const button = document.getElementById(buttonId);
  const yaml = document.getElementById(contentId);
  if (!button || !yaml) return;

  button.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(yaml.textContent || "");
      button.textContent = "Copied";
      window.setTimeout(() => { button.textContent = "Copy YAML"; }, 1500);
    } catch {
      button.textContent = "Copy failed";
      window.setTimeout(() => { button.textContent = "Copy YAML"; }, 1500);
    }
  });
}

function renderYAMLContent(id, value) {
  const element = document.getElementById(id);
  if (!element) return;

  const lines = String(value || "").replace(/\r\n/g, "\n").split("\n");
  element.innerHTML = lines.map((line, index) => `<div class="pod-yaml-line"><span class="pod-yaml-line-number">${index + 1}</span><span class="pod-yaml-line-content">${highlightYAML(line)}</span></div>`).join("");
}

function highlightYAML(line) {
  const escaped = escapeHtml(line).replace(/ /g, "&nbsp;");
  if (!escaped.trim()) return "&nbsp;";
  if (escaped.trimStart().startsWith("#")) return `<span class="pod-yaml-comment">${escaped}</span>`;

  const keyMatch = escaped.match(/^(\s*-\s*)?([^:&][^:]*?):\s*(.*)$/);
  if (!keyMatch) return highlightYAMLValue(escaped);

  const prefix = keyMatch[1] || "";
  const key = keyMatch[2] || "";
  const value = keyMatch[3] || "";
  return `${prefix}<span class="pod-yaml-key">${key}</span>:${value ? ` ${highlightYAMLValue(value)}` : ""}`;
}

function highlightYAMLValue(value) {
  const trimmed = value.replace(/&nbsp;/g, " ").trim();
  if (!trimmed) return value;

  const unquoted = trimmed.replace(/^"|"$/g, "");
  if (/^\d{4}-\d{2}-\d{2}t\d{2}:\d{2}:\d{2}(?:\.\d+)?z$/i.test(unquoted)) return `<span class="pod-yaml-timestamp">${value}</span>`;
  if (/^(true|false|null)$/i.test(trimmed)) return `<span class="pod-yaml-boolean">${value}</span>`;
  if (/^-?\d+(\.\d+)?$/.test(trimmed)) return `<span class="pod-yaml-number">${value}</span>`;
  return `<span class="pod-yaml-string">${value}</span>`;
}

function renderDeploymentError(message) {
  setText("deploymentHeroTitle", "Deployment unavailable");
  setText("deploymentName", message);
  applyDeploymentStatus("deploymentStatus", "Unavailable");

  ["deploymentConditionsList", "deploymentReplicaSetsList", "deploymentPodsList", "deploymentEventsList", "deploymentLabelsList", "deploymentAnnotationsList"].forEach((id) => {
    const element = document.getElementById(id);
    if (element) element.innerHTML = `<p class="pod-empty">${escapeHtml(message)}</p>`;
  });

  renderYAMLContent("deploymentYamlContent", message);
}

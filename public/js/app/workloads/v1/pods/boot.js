let podItems = [];

const DEFAULT_POD_YAML = `apiVersion: v1
kind: Pod
metadata:
  name: nginx
  namespace: default
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`;

document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("deploymentSearch");

  initCreateResourceModal({
    triggerId: "createPodButton",
    title: "Create Pod",
    description: "Edit the pod YAML before confirming.",
    initialValue: DEFAULT_POD_YAML,
    confirmLabel: "Confirm",
    pendingLabel: "Creating...",
    async onConfirm(content) {
      const data = await fetchClusterData("/api/workloads/pods", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content }),
      });

      window.setTimeout(() => window.location.reload(), 150);
      return {
        message: `Created pod ${data.name || "resource"} in ${data.namespace || "default"}.`,
      };
    },
  });

  searchInput?.addEventListener("input", () =>
    renderFilteredPods(searchInput.value),
  );

  startAutoRefresh(async () => {
    try {
      const data = await fetchClusterData("/api/workloads/pods");
      podItems = data.items || [];
      applyPageMeta(data.meta, { namespaceId: "podsNamespaceFilter" });
      setText("podsHeroTitle", `Pods (${data.stats?.running || 0} running)`);
      renderFilteredPods(searchInput?.value || "");
    } catch (error) {
      podItems = [];
      renderPodsTable([], error.message || "Failed to load pods");
    }
  });
});

function renderFilteredPods(query) {
  const filteredItems = filterPods(podItems, query);
  renderPodsTable(
    filteredItems,
    filteredItems.length ? "" : "No pods match your search",
  );
  setText("podsTableCount", `${filteredItems.length} rows`);
}

function filterPods(items, query) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) return items;

  return items.filter((item) =>
    [item.namespace, item.name, item.status, item.node, item.podIP].some(
      (value) =>
        String(value || "")
          .toLowerCase()
          .includes(normalized),
    ),
  );
}

function renderPodsTable(items, message) {
  const body = document.getElementById("podsTableBody");
  if (!body) return;

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="11">${escapeHtml(message || "No pods found")}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.className = "deployments-table__row-link";
    row.tabIndex = 0;
    row.innerHTML = `
      <td>${escapeHtml(item.namespace)}</td>
      <td>${escapeHtml(item.name)}</td>
      <td>${escapeHtml(item.ready)}</td>
      <td>${escapeHtml(item.status)}</td>
      <td>${escapeHtml(String(item.restarts))}</td>
      <td>${escapeHtml(item.lastRestart)}</td>
      <td>${escapeHtml(item.lastRestartReason)}</td>
      <td>${escapeHtml(item.node)}</td>
      <td>${escapeHtml(item.podIP)}</td>
      <td>${escapeHtml(item.age)}</td>
      <td class="deployments-table__actions">
        <button class="action-button action-button--danger action-button--compact" type="button" data-pod-delete>
          Delete
        </button>
      </td>
    `;
    row.addEventListener("click", () => openPod(item));
    row.addEventListener("keydown", (event) => {
      if (event.target instanceof HTMLElement && event.target.closest("button")) {
        return;
      }
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        openPod(item);
      }
    });
    const deleteButton = row.querySelector("[data-pod-delete]");
    deleteButton?.addEventListener("click", (event) => {
      event.stopPropagation();
      void deletePodFromTable(item, deleteButton);
    });
    body.appendChild(row);
  });
}

async function deletePodFromTable(item, button) {
  const confirmed = window.confirm(
    `Delete pod ${item.name || "resource"} from namespace ${item.namespace || "default"}?`,
  );
  if (!confirmed) return;

  const defaultLabel = button.textContent;
  button.disabled = true;
  button.textContent = "Deleting...";

  try {
    await fetchClusterData(
      `/api/workloads/pod?${new URLSearchParams({ namespace: item.namespace || "", name: item.name || "" }).toString()}`,
      { method: "DELETE" },
    );
    podItems = podItems.filter(
      (candidate) =>
        !(
          candidate.namespace === item.namespace && candidate.name === item.name
        ),
    );
    renderFilteredPods(
      document.getElementById("deploymentSearch")?.value || "",
    );
  } catch (error) {
    window.alert(error.message || "Failed to delete pod");
    button.disabled = false;
    button.textContent = defaultLabel;
  }
}

function openPod(item) {
  const params = new URLSearchParams({
    namespace: item.namespace || "",
    name: item.name || "",
  });
  window.location.href = `/workloads/manage/pod?${params.toString()}`;
}

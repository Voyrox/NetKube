const CLUSTER_CONTEXT_HEADER = "X-NetKube-Context";

async function fetchClusterData(path, options = {}) {
  const contextId = window.NetKubeStorage?.getActiveContextId();
  if (!contextId) {
    throw new Error("No cluster selected");
  }

  const headers = {
    [CLUSTER_CONTEXT_HEADER]: contextId,
    ...(options.headers || {}),
  };

  const response = await fetch(path, {
    ...options,
    headers,
  });

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(data.error || "Failed to load cluster data");
  }

  return data;
}

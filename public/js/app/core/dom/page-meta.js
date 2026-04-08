function applyPageMeta(meta, options = {}) {
  if (!meta) {
    return;
  }

  setText(
    options.userId || "heroUser",
    meta.userName || meta.contextName || "Unknown user",
  );
  setText(
    options.contextId || "selectedContext",
    meta.clusterName || meta.contextName || "Unknown cluster",
  );
  setText(options.refreshId || "lastRefresh", formatRefresh(meta.lastRefresh));

  if (options.namespaceId) {
    setText(options.namespaceId, meta.namespace || "All namespaces");
  }
}

function formatRefresh(value) {
  if (!value) {
    return "just now";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}

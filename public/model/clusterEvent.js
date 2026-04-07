document.addEventListener("DOMContentLoaded", async () => {
  const params = new URLSearchParams(window.location.search);
  const query = new URLSearchParams();
  ["namespace", "name", "reason", "kind"].forEach((key) => {
    if (params.get(key)) {
      query.set(key, params.get(key));
    }
  });

  try {
    const data = await fetchClusterData(`/api/cluster/event${query.toString() ? `?${query.toString()}` : ""}`);
    applyPageMeta(data.meta);
    renderClusterEvent(data.item || {});
  } catch (error) {
    renderClusterEventError(error.message || "Failed to load event details");
  }
});

function renderClusterEvent(item) {
  setText("eventHeroTitle", item.title || "Event");
  setText("eventTypeMeta", item.type || "-");
  setText("eventNamespaceMeta", item.namespace || "-");
  setText("eventTitle", item.title || "-");
  setText("eventSummary", item.message || "-");
  setText("eventReason", item.reason || "-");
  setText("eventObject", item.involvedObject || "-");
  setText("eventSource", item.source || "-");
  setText("eventFirstSeen", item.firstSeen || "-");
  setText("eventLastSeen", item.lastSeen || "-");
  setText("eventCount", item.count || 0);
  setText("eventKind", item.kind || "-");
  setText("eventName", item.name || "-");
  setText("eventNamespace", item.namespace || "-");
  setText("eventNode", item.node || "-");
  setText("eventMessage", item.message || "-");

  const badge = document.getElementById("eventTypeBadge");
  if (badge) {
    badge.textContent = item.type || "-";
    badge.className = `cluster-event-badge ${eventBadgeClass(item.type)}`;
  }

  renderEventTimeline(item.timeline || []);
}

function renderEventTimeline(items) {
  const element = document.getElementById("eventTimelineList");
  if (!element) return;
  if (!items.length) {
    element.innerHTML = '<p class="cluster-empty">No related event history found.</p>';
    return;
  }
  element.innerHTML = items.map((item) => `<article class="cluster-timeline-item"><div class="cluster-timeline-top"><span class="cluster-timeline-title">${escapeHtml(item.title || "Event")}</span><span class="cluster-timeline-time">${escapeHtml(item.age || "-")}</span></div><p class="cluster-timeline-message">${escapeHtml(item.message || "-")}</p></article>`).join("");
}

function eventBadgeClass(type) {
  const normalized = String(type || "").toLowerCase();
  if (normalized === "normal") return "cluster-event-badge--healthy";
  if (normalized === "warning") return "cluster-event-badge--warning";
  return "cluster-event-badge--danger";
}

function renderClusterEventError(message) {
  setText("eventHeroTitle", "Event unavailable");
  setText("eventTitle", message);
  setText("eventSummary", message);
  const badge = document.getElementById("eventTypeBadge");
  if (badge) {
    badge.textContent = "Unavailable";
    badge.className = "cluster-event-badge cluster-event-badge--danger";
  }
  const timeline = document.getElementById("eventTimelineList");
  if (timeline) {
    timeline.innerHTML = `<p class="cluster-empty">${escapeHtml(message)}</p>`;
  }
}

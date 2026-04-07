const CONTEXTS_API = "/api/config/contexts";

document.addEventListener("DOMContentLoaded", async () => {
  const heroUser = document.getElementById("heroUser");
  const selectedContext = document.getElementById("selectedContext");
  const lastRefresh = document.getElementById("lastRefresh");

  const activeContextId = window.NetKubeStorage?.getActiveContextId();
  if (!activeContextId) {
    applyOverviewFallback(heroUser, selectedContext, lastRefresh);
    return;
  }

  try {
    const response = await fetch(CONTEXTS_API);
    if (!response.ok) {
      applyOverviewFallback(heroUser, selectedContext, lastRefresh);
      return;
    }

    const data = await response.json();
    const contexts = Array.isArray(data.contexts) ? data.contexts : [];
    const context = contexts.find((item) => item.id === activeContextId);

    if (!context) {
      window.NetKubeStorage?.clearActiveContextId();
      applyOverviewFallback(heroUser, selectedContext, lastRefresh);
      return;
    }

    if (heroUser) {
      heroUser.textContent = context.userName || context.contextName || "Unknown user";
    }

    if (selectedContext) {
      selectedContext.textContent = context.clusterName || context.contextName || "Unknown cluster";
    }

    if (lastRefresh) {
      lastRefresh.textContent = new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit"
      });
    }
  } catch (error) {
    console.error("Failed to load active context for overview", error);
    applyOverviewFallback(heroUser, selectedContext, lastRefresh);
  }
});

function applyOverviewFallback(heroUser, selectedContext, lastRefresh) {
  if (heroUser) {
    heroUser.textContent = "No cluster selected";
  }

  if (selectedContext) {
    selectedContext.textContent = "Choose a cluster from home";
  }

  if (lastRefresh) {
    lastRefresh.textContent = "not connected";
  }
}

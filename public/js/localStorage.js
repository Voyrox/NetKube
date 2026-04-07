const ACTIVE_CONTEXT_STORAGE_KEY = "netkube.activeContextId";
const CONTEXT_LAYOUT_STORAGE_KEY = "netkube.contextLayout";

function setActiveContextId(contextId) {
  try {
    if (!contextId) {
      sessionStorage.removeItem(ACTIVE_CONTEXT_STORAGE_KEY);
      return;
    }

    sessionStorage.setItem(ACTIVE_CONTEXT_STORAGE_KEY, contextId);
  } catch (error) {
    console.error("Failed to save active context id", error);
  }
}

function getActiveContextId() {
  try {
    return sessionStorage.getItem(ACTIVE_CONTEXT_STORAGE_KEY) || "";
  } catch (error) {
    console.error("Failed to load active context id", error);
    return "";
  }
}

function clearActiveContextId() {
  try {
    sessionStorage.removeItem(ACTIVE_CONTEXT_STORAGE_KEY);
  } catch (error) {
    console.error("Failed to clear active context id", error);
  }
}

function setContextLayout(layout) {
  try {
    if (!layout) {
      localStorage.removeItem(CONTEXT_LAYOUT_STORAGE_KEY);
      return;
    }

    localStorage.setItem(CONTEXT_LAYOUT_STORAGE_KEY, layout);
  } catch (error) {
    console.error("Failed to save context layout", error);
  }
}

function getContextLayout() {
  try {
    return localStorage.getItem(CONTEXT_LAYOUT_STORAGE_KEY) || "";
  } catch (error) {
    console.error("Failed to load context layout", error);
    return "";
  }
}

window.NetKubeStorage = {
  setActiveContextId,
  getActiveContextId,
  clearActiveContextId,
  setContextLayout,
  getContextLayout
};

const ACTIVE_CONTEXT_STORAGE_KEY = "netkube.activeContextId";

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

window.NetKubeStorage = {
  setActiveContextId,
  getActiveContextId,
  clearActiveContextId
};

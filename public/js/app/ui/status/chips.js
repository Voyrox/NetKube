function applyStatusChip(id, status) {
  const element = document.getElementById(id);
  if (!element) {
    return;
  }

  element.textContent = status;
  element.classList.remove(
    "status-chip--healthy",
    "status-chip--warning",
    "status-chip--danger",
  );

  const normalized = String(status || "").toLowerCase();
  if (normalized.includes("healthy")) {
    element.classList.add("status-chip--healthy");
    return;
  }

  if (normalized.includes("watch") || normalized.includes("update")) {
    element.classList.add("status-chip--warning");
    return;
  }

  element.classList.add("status-chip--danger");
}

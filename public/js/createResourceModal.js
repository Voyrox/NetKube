function initCreateResourceModal(config) {
  const trigger = document.getElementById(config.triggerId);
  const backdrop = document.getElementById("createResourceBackdrop");
  const modal = document.getElementById("createResourceModal");
  const title = document.getElementById("createResourceTitle");
  const description = document.getElementById("createResourceDescription");
  const lineNumbers = document.getElementById("createResourceLineNumbers");
  const highlight = document.getElementById("createResourceHighlight");
  const input = document.getElementById("createResourceInput");
  const closeButton = document.getElementById("createResourceClose");
  const cancelButton = document.getElementById("createResourceCancel");
  const confirmButton = document.getElementById("createResourceConfirm");
  const status = document.getElementById("createResourceStatus");

  if (!trigger || !backdrop || !modal || !title || !description || !lineNumbers || !highlight || !input || !closeButton || !cancelButton || !confirmButton || !status) {
    return;
  }

  let lastFocusedElement = null;
  const defaultConfirmLabel = confirmButton.textContent;

  function setStatus(message, state = "") {
    status.textContent = message || "";
    if (state) {
      status.dataset.state = state;
      return;
    }

    delete status.dataset.state;
  }

  function renderEditor(value) {
    const normalizedValue = String(value || "").replace(/\r\n/g, "\n");
    const lineCount = Math.max(normalizedValue.split("\n").length, 1);

    lineNumbers.innerHTML = Array.from({ length: lineCount }, (_, index) => {
      return `<div class="create-resource-modal__line-number">${index + 1}</div>`;
    }).join("");
    highlight.innerHTML = window.NetKubeYaml?.renderHighlightedYamlContent(normalizedValue) || "";
  }

  function syncScroll() {
    highlight.scrollTop = input.scrollTop;
    highlight.scrollLeft = input.scrollLeft;
    lineNumbers.scrollTop = input.scrollTop;
  }

  function closeModal() {
    backdrop.hidden = true;
    modal.hidden = true;
    document.body.classList.remove("resource-modal-open");
    setStatus("");
    if (lastFocusedElement instanceof HTMLElement) {
      lastFocusedElement.focus();
    }
  }

  function openModal() {
    lastFocusedElement = document.activeElement;
    title.textContent = config.title;
    description.textContent = config.description;
    input.value = config.initialValue;
    renderEditor(input.value);
    syncScroll();
    setStatus("");
    confirmButton.disabled = false;
    confirmButton.textContent = config.confirmLabel || defaultConfirmLabel;

    backdrop.hidden = false;
    modal.hidden = false;
    document.body.classList.add("resource-modal-open");

    window.requestAnimationFrame(() => {
      input.focus();
      input.setSelectionRange(input.value.length, input.value.length);
    });
  }

  trigger.addEventListener("click", openModal);
  closeButton.addEventListener("click", closeModal);
  cancelButton.addEventListener("click", closeModal);
  backdrop.addEventListener("click", closeModal);

  input.addEventListener("input", () => {
    renderEditor(input.value);
    syncScroll();
  });
  input.addEventListener("scroll", syncScroll);

  confirmButton.addEventListener("click", async () => {
    if (typeof config.onConfirm !== "function") {
      closeModal();
      return;
    }

    confirmButton.disabled = true;
    confirmButton.textContent = config.pendingLabel || "Creating...";
    setStatus("");

    try {
      const result = await config.onConfirm(input.value);
      setStatus(result?.message || "Created successfully.", "success");
      closeModal();
    } catch (error) {
      setStatus(error?.message || "Failed to create resource.", "error");
      confirmButton.disabled = false;
      confirmButton.textContent = config.confirmLabel || defaultConfirmLabel;
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && !modal.hidden) {
      event.preventDefault();
      closeModal();
    }
  });
}

window.initCreateResourceModal = initCreateResourceModal;

function renderYamlDrawerContent(id, value) {
  const element = document.getElementById(id);
  if (!element) {
    return;
  }

  element.classList.add("resource-drawer__content--yaml");
  element.innerHTML = window.NetKubeYaml?.renderHighlightedYaml(value) || "";
}

function setDrawerText(id, value) {
  const element = document.getElementById(id);
  if (!element) {
    return;
  }

  element.classList.remove("resource-drawer__content--yaml");
  element.textContent = value;
}

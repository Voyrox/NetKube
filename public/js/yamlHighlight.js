function escapeYamlHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function highlightYamlLine(line) {
  const escaped = escapeYamlHtml(line).replace(/ /g, "&nbsp;");
  if (!escaped.trim()) {
    return "&nbsp;";
  }

  if (escaped.trimStart().startsWith("#")) {
    return `<span class="nk-yaml-comment">${escaped}</span>`;
  }

  const keyMatch = escaped.match(/^(\s*-\s*)?([^:&][^:]*?):\s*(.*)$/);
  if (!keyMatch) {
    return highlightYamlValue(escaped);
  }

  const prefix = keyMatch[1] || "";
  const key = keyMatch[2] || "";
  const value = keyMatch[3] || "";

  return `${prefix}<span class="nk-yaml-key">${key}</span>:${value ? ` ${highlightYamlValue(value)}` : ""}`;
}

function highlightYamlValue(value) {
  const trimmed = value.replace(/&nbsp;/g, " ").trim();
  if (!trimmed) {
    return value;
  }

  const unquoted = trimmed.replace(/^"|"$/g, "");

  if (/^\d{4}-\d{2}-\d{2}t\d{2}:\d{2}:\d{2}(?:\.\d+)?z$/i.test(unquoted)) {
    return `<span class="nk-yaml-timestamp">${value}</span>`;
  }

  if (/^(true|false|null)$/i.test(trimmed)) {
    return `<span class="nk-yaml-boolean">${value}</span>`;
  }

  if (/^-?\d+(\.\d+)?$/.test(trimmed)) {
    return `<span class="nk-yaml-number">${value}</span>`;
  }

  return `<span class="nk-yaml-string">${value}</span>`;
}

function renderHighlightedYaml(value) {
  return String(value || "")
    .replace(/\r\n/g, "\n")
    .split("\n")
    .map((line, index) => {
      return `<div class="nk-yaml-line"><span class="nk-yaml-line-number">${index + 1}</span><span class="nk-yaml-line-content">${highlightYamlLine(line)}</span></div>`;
    })
    .join("");
}

window.NetKubeYaml = {
  renderHighlightedYaml,
  highlightYamlLine,
  highlightYamlValue,
  escapeYamlHtml
};

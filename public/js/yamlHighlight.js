function escapeYamlHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function nbsp(value) {
  return escapeYamlHtml(value).replace(/ /g, "&nbsp;");
}

function highlightYamlPunctuation(value) {
  return value.replace(/[{}\[\]]/g, (char) => {
    return `<span class="nk-yaml-punctuation">${char}</span>`;
  });
}

function splitYamlKeyValue(line) {
  const listPrefixMatch = line.match(/^(\s*-\s*)/);
  const prefix = listPrefixMatch ? listPrefixMatch[1] : "";
  const rest = line.slice(prefix.length);

  const lastColon = rest.lastIndexOf(":");
  if (lastColon === -1) {
    return null;
  }

  const key = rest.slice(0, lastColon);
  const afterColon = rest.slice(lastColon + 1);

  if (!key.trim()) {
    return null;
  }

  if (afterColon.length > 0 && !/^\s/.test(afterColon)) {
    return null;
  }

  return {
    prefix,
    key,
    value: afterColon.trimStart()
  };
}

function highlightYamlLine(line) {
  if (!line.trim()) {
    return "&nbsp;";
  }

  if (line.trimStart().startsWith("#")) {
    return `<span class="nk-yaml-comment">${highlightYamlPunctuation(nbsp(line))}</span>`;
  }

  const parts = splitYamlKeyValue(line);
  if (!parts) {
    return highlightYamlValue(highlightYamlPunctuation(nbsp(line)));
  }

  const { prefix, key, value } = parts;

  return `${highlightYamlPunctuation(nbsp(prefix))}<span class="nk-yaml-key">${highlightYamlPunctuation(nbsp(key))}</span>:${
    value ? ` ${highlightYamlValue(highlightYamlPunctuation(nbsp(value)))}` : ""
  }`;
}

function highlightYamlValue(value) {
  const trimmed = value.replace(/&nbsp;/g, " ").trim();
  if (!trimmed) {
    return value;
  }

  const unquoted = trimmed.replace(/^"|"$/g, "");

  if (/^[|>]([+-])?$/.test(trimmed)) {
    return `<span class="nk-yaml-operator">${value}</span>`;
  }

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

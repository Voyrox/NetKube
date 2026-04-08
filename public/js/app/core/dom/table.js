function filterByQuery(items, query, fields) {
  const normalized = String(query || "")
    .trim()
    .toLowerCase();
  if (!normalized) {
    return items;
  }

  return items.filter((item) =>
    fields.some((field) =>
      String(item?.[field] || "")
        .toLowerCase()
        .includes(normalized),
    ),
  );
}

function renderTableRows(bodyId, columnCount, items, emptyMessage, renderRow) {
  const body = document.getElementById(bodyId);
  if (!body) {
    return;
  }

  body.innerHTML = "";
  if (!items.length) {
    body.innerHTML = `<tr><td colspan="${columnCount}">${escapeHtml(emptyMessage)}</td></tr>`;
    return;
  }

  items.forEach((item) => {
    const row = document.createElement("tr");
    row.innerHTML = renderRow(item);
    body.appendChild(row);
  });
}

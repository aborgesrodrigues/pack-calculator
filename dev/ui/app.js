function parseSizes(raw) {
  const sizes = raw
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s) => {
      const n = Number(s);
      if (!Number.isInteger(n) || n <= 0) {
        throw new Error(`Invalid size: ${s}`);
      }
      return n;
    });

  if (sizes.length === 0) {
    throw new Error("At least one pack size is required");
  }

  return sizes;
}

function showResponse(containerId, statusId, bodyId, status, bodyText) {
  const container = document.getElementById(containerId);
  const statusEl = document.getElementById(statusId);
  const bodyEl = document.getElementById(bodyId);

  container.hidden = false;
  statusEl.textContent = `${status}`;
  statusEl.className = `status ${status >= 200 && status < 300 ? "ok" : "err"}`;
  bodyEl.textContent = bodyText;
}

async function request(url, payload) {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const text = await res.text();
  let formatted = text;
  try {
    formatted = JSON.stringify(JSON.parse(text), null, 2);
  } catch {
    /* keep raw text */
  }

  return { status: res.status, body: formatted };
}

document.getElementById("pack-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const btn = e.target.querySelector("button");
  btn.disabled = true;

  try {
    const sizes = parseSizes(document.getElementById("pack-sizes").value);
    const { status, body } = await request("/pack_size/batch", { sizes });
    showResponse("pack-response", "pack-status", "pack-body", status, body);
  } catch (err) {
    showResponse("pack-response", "pack-status", "pack-body", 0, err.message);
    document.getElementById("pack-status").className = "status err";
    document.getElementById("pack-status").textContent = "Error";
  } finally {
    btn.disabled = false;
  }
});

document.getElementById("calculate-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const btn = e.target.querySelector("button");
  btn.disabled = true;

  try {
    const items = Number(document.getElementById("items").value);
    if (!Number.isInteger(items) || items <= 0) {
      throw new Error("Items must be a positive integer");
    }

    const { status, body } = await request("/calculate", { items });
    showResponse("calculate-response", "calculate-status", "calculate-body", status, body);
  } catch (err) {
    showResponse("calculate-response", "calculate-status", "calculate-body", 0, err.message);
    document.getElementById("calculate-status").className = "status err";
    document.getElementById("calculate-status").textContent = "Error";
  } finally {
    btn.disabled = false;
  }
});

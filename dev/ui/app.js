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

function formatSizes(sizes) {
  return sizes.join(", ");
}

function setPackSizesInput(sizes) {
  document.getElementById("pack-sizes").value =
    sizes && sizes.length > 0 ? formatSizes(sizes) : "";
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

async function request(url, options = {}) {
  const { method = "POST", payload } = options;
  const init = {
    method,
    headers: {},
  };

  if (payload !== undefined) {
    init.headers["Content-Type"] = "application/json";
    init.body = JSON.stringify(payload);
  }

  const res = await fetch(url, init);
  const text = await res.text();
  let formatted = text;
  let data;

  try {
    data = JSON.parse(text);
    formatted = JSON.stringify(data, null, 2);
  } catch {
    /* keep raw text */
  }

  return { status: res.status, body: formatted, data };
}

async function loadPackSizes() {
  const input = document.getElementById("pack-sizes");
  input.disabled = true;

  try {
    const { status, data } = await request("/pack_size/batch", { method: "GET" });
    if (status >= 200 && status < 300 && data && Array.isArray(data.sizes)) {
      setPackSizesInput(data.sizes);
    }
  } catch {
    /* leave field empty on load failure */
  } finally {
    input.disabled = false;
  }
}

document.getElementById("pack-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const btn = e.target.querySelector("button");
  btn.disabled = true;

  try {
    const sizes = parseSizes(document.getElementById("pack-sizes").value);
    const { status, body } = await request("/pack_size/batch", { payload: { sizes } });
    showResponse("pack-response", "pack-status", "pack-body", status, body);
    if (status >= 200 && status < 300) {
      setPackSizesInput(sizes);
    }
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

    const { status, body } = await request("/calculate", { payload: { items } });
    showResponse("calculate-response", "calculate-status", "calculate-body", status, body);
  } catch (err) {
    showResponse("calculate-response", "calculate-status", "calculate-body", 0, err.message);
    document.getElementById("calculate-status").className = "status err";
    document.getElementById("calculate-status").textContent = "Error";
  } finally {
    btn.disabled = false;
  }
});

loadPackSizes();

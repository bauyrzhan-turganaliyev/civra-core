const API_BASE = ""; // same origin (gateway serves UI)

const el = (id) => document.getElementById(id);
const state = {
  userId: localStorage.getItem("civra.userId") || "",
  kingdomId: localStorage.getItem("civra.kingdomId") || "",
};

function setError(msg) {
  el("lastErr").textContent = msg || "";
}
function setLoginError(msg) {
  el("loginErr").textContent = msg || "";
}

async function api(path, opts = {}) {
  setError("");
  const res = await fetch(API_BASE + path, {
    ...opts,
    headers: {
      "Content-Type": "application/json",
      ...(opts.headers || {}),
    },
  });

  const text = await res.text();
  let data = null;
  try { data = text ? JSON.parse(text) : null; } catch { data = { raw: text }; }

  if (!res.ok) {
    const msg = data?.error ? data.error : `HTTP ${res.status}`;
    throw new Error(msg);
  }
  return data;
}

function profToResource(prof) {
  if (prof === "farmer") return "food";
  if (prof === "miner") return "iron";
  if (prof === "lumber") return "wood";
  return "food";
}

function renderInv(containerId, invObj) {
  const root = el(containerId);
  root.innerHTML = "";
  const entries = Object.entries(invObj || {});
  if (entries.length === 0) {
    root.innerHTML = `<div class="muted">Empty</div>`;
    return;
  }
  for (const [name, qty] of entries) {
    const div = document.createElement("div");
    div.className = "item";
    div.innerHTML = `<div class="name">${name}</div><div class="qty">${qty}</div>`;
    root.appendChild(div);
  }
}

function renderOrders(orders) {
  const root = el("orders");
  root.innerHTML = "";
  if (!orders || orders.length === 0) {
    root.innerHTML = `<div class="muted">No orders</div>`;
    return;
  }

  for (const o of orders) {
    const div = document.createElement("div");
    div.className = "order";
    div.innerHTML = `
      <div class="top">
        <div>
          <div class="id">${o.id}</div>
        </div>
        <div class="muted small">${new Date(o.createdAt).toLocaleString()}</div>
      </div>
      <div class="meta">
        <div><b>resource</b> ${o.resource}</div>
        <div><b>qty</b> ${o.quantity}</div>
        <div><b>price</b> ${o.price}</div>
        <div><b>seller</b> ${o.sellerId}</div>
      </div>
      <div class="actions">
        <button class="secondary" data-buy="${o.id}">Buy (as ${state.userId})</button>
        <button class="danger" data-cancel="${o.id}">Cancel (as seller)</button>
      </div>
    `;
    root.appendChild(div);
  }

  root.querySelectorAll("[data-buy]").forEach(btn => {
    btn.addEventListener("click", async () => {
      try {
        await api("/economy/market/buy", {
          method: "POST",
          body: JSON.stringify({ orderId: btn.dataset.buy, buyerId: state.userId }),
        });
        await refreshAll();
      } catch (e) { setError(String(e.message || e)); }
    });
  });

  root.querySelectorAll("[data-cancel]").forEach(btn => {
    btn.addEventListener("click", async () => {
      try {
        await api("/economy/market/cancel", {
          method: "POST",
          body: JSON.stringify({ orderId: btn.dataset.cancel, sellerId: state.userId }),
        });
        await refreshAll();
      } catch (e) { setError(String(e.message || e)); }
    });
  });
}

async function refreshAll() {
  if (!state.userId || !state.kingdomId) return;

  el("userIdView").textContent = state.userId;
  el("kingdomIdView").textContent = state.kingdomId;

  const [p, k, orders] = await Promise.all([
    api(`/economy/personal-inventory?userId=${encodeURIComponent(state.userId)}`),
    api(`/economy/kingdom-inventory?kingdomId=${encodeURIComponent(state.kingdomId)}`),
    api(`/economy/market/orders?kingdomId=${encodeURIComponent(state.kingdomId)}`),
  ]);

  renderInv("personalInv", p.inventory || {});
  renderInv("kingdomInv", k.inventory || {});
  renderOrders(orders.orders || []);
}

function showApp() {
  el("loginCard").classList.add("hidden");
  el("app").classList.remove("hidden");
}

function showLogin() {
  el("loginCard").classList.remove("hidden");
  el("app").classList.add("hidden");
}

function saveSession() {
  localStorage.setItem("civra.userId", state.userId);
  localStorage.setItem("civra.kingdomId", state.kingdomId);
}

function clearSession() {
  localStorage.removeItem("civra.userId");
  localStorage.removeItem("civra.kingdomId");
  state.userId = "";
  state.kingdomId = "";
}

async function onLogin() {
  const userId = el("userIdInput").value.trim();
  const kingdomId = el("kingdomIdInput").value.trim();
  setLoginError("");

  if (!userId || !kingdomId) {
    setLoginError("userId and kingdomId are required");
    return;
  }

  state.userId = userId;
  state.kingdomId = kingdomId;
  saveSession();

  showApp();
  await refreshAll();
}

async function onGather() {
  try {
    const prof = el("profSelect").value;
    const res = profToResource(prof);
    const amount = Number(el("gatherAmount").value || 0);
    const body = {
      userId: state.userId,
      kingdomId: state.kingdomId,
      profession: prof,
      resource: res,
      amount,
    };

    const out = await api("/economy/gather", {
      method: "POST",
      body: JSON.stringify(body),
    });

    el("gatherOut").textContent = JSON.stringify(out, null, 2);
    await refreshAll();
  } catch (e) {
    setError(String(e.message || e));
  }
}

async function onSell() {
  try {
    const body = {
      kingdomId: state.kingdomId,
      sellerId: state.userId,
      resource: el("sellRes").value,
      quantity: Number(el("sellQty").value || 0),
      price: Number(el("sellPrice").value || 0),
    };

    const out = await api("/economy/market/sell", {
      method: "POST",
      body: JSON.stringify(body),
    });

    el("sellOut").textContent = JSON.stringify(out, null, 2);
    await refreshAll();
  } catch (e) {
    setError(String(e.message || e));
  }
}

function init() {
  el("fillDemoBtn").addEventListener("click", () => {
    el("userIdInput").value = "u1";
    el("kingdomIdInput").value = "k1";
  });

  el("loginBtn").addEventListener("click", onLogin);
  el("logoutBtn").addEventListener("click", () => {
    clearSession();
    showLogin();
  });

  el("refreshBtn").addEventListener("click", () => refreshAll().catch(e => setError(e.message)));
  el("reloadOrdersBtn").addEventListener("click", () => refreshAll().catch(e => setError(e.message)));

  el("gatherBtn").addEventListener("click", onGather);
  el("sellBtn").addEventListener("click", onSell);

  // auto session
  if (state.userId && state.kingdomId) {
    showApp();
    refreshAll().catch(e => setError(e.message));
  } else {
    showLogin();
  }
}

init();

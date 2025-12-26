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
     credentials: "same-origin",
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
  await loadItems();
  await loadItemOrders();


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

  if (!userId || !kingdomId) {
    setLoginError("userId and kingdomId are required");
    return;
  }

  await api("/auth/login", {
    method: "POST",
    body: JSON.stringify({ userId, kingdomId }),
  });

  state.userId = userId;
  state.kingdomId = kingdomId;

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
el("logoutBtn").addEventListener("click", async () => {
  try { await api("/auth/logout", { method: "POST", body: "{}" }); } catch {}
  clearSession();
  showLogin();
});


  el("refreshBtn").addEventListener("click", () => refreshAll().catch(e => setError(e.message)));
  el("reloadOrdersBtn").addEventListener("click", () => refreshAll().catch(e => setError(e.message)));

  el("gatherBtn").addEventListener("click", onGather);
  el("sellBtn").addEventListener("click", onSell);

  el("reloadItemsBtn").addEventListener("click", () => loadItems().catch(e => setError(e.message)));

  el("craftT1Btn").addEventListener("click", async () => {
    try {
      const out = await api("/economy/items/craft-tool", {
        method: "POST",
        body: JSON.stringify({ userId: state.userId, tier: 1 }),
      });
      el("itemsOut").textContent = JSON.stringify(out, null, 2);
      await loadItems();
    } catch (e) { setError(e.message || String(e)); }
  });
  el("craftT2Btn").addEventListener("click", async () => {
    try {
      const out = await api("/economy/items/craft-tool", {
        method: "POST",
        body: JSON.stringify({ userId: state.userId, tier: 2 }),
      });
      el("itemsOut").textContent = JSON.stringify(out, null, 2);
      await loadItems();
    } catch (e) { setError(e.message || String(e)); }
  });
  el("craftT3Btn").addEventListener("click", async () => {
    try {
      const out = await api("/economy/items/craft-tool", {
        method: "POST",
        body: JSON.stringify({ userId: state.userId, tier: 3 }),
      });
      el("itemsOut").textContent = JSON.stringify(out, null, 2);
      await loadItems();
    } catch (e) { setError(e.message || String(e)); }
  });
  el("reloadItemOrdersBtn").addEventListener("click", () =>
  loadItemOrders().catch(e => setError(e.message || String(e)))
);

  (async () => {
    try {
      const me = await api("/auth/me");
      state.userId = me.userId;
      state.kingdomId = me.kingdomId;
      saveSession();
      showApp();
      await refreshAll();
    } catch {
      showLogin();
    }
  })();
}

async function loadItems() {
  const data = await api(`/economy/items?userId=${encodeURIComponent(state.userId)}`);
  renderItems(data.items || []);
}

function renderItems(items) {
  const root = el("items");
  root.innerHTML = "";
  if (!items.length) {
    root.innerHTML = `<div class="muted">No items</div>`;
    return;
  }

  for (const it of items) {
    const div = document.createElement("div");
    div.className = "order";
    div.innerHTML = `
      <div class="top">
        <div class="id">${it.id}</div>
        <div class="muted small">${new Date(it.createdAt).toLocaleString()}</div>
      </div>
      <div class="meta">
        <div><b>type</b> ${it.itemType}</div>
        <div><b>tier</b> ${it.tier}</div>
        <div><b>dur</b> ${it.durability}/${it.maxDurability}</div>
        <div><b>bonus</b> +${it.bonusPct}%</div>
        <div><b>equipped</b> ${it.equipped ? "✅" : "—"}</div>
        <div><b>listed</b> ${it.listed ? "🟠" : "—"}</div>
      </div>
      <div class="actions">
        <button class="secondary" data-equip="${it.id}" ${it.listed ? "disabled" : ""}>Equip</button>
        <button class="danger" data-sell="${it.id}" ${it.equipped || it.listed ? "disabled" : ""}>Sell</button>
      </div>
    `;
    root.appendChild(div);
  }

  root.querySelectorAll("[data-equip]").forEach(btn => {
    btn.addEventListener("click", async () => {
      try {
        await api("/economy/items/equip", {
          method: "POST",
          body: JSON.stringify({ userId: state.userId, itemId: btn.dataset.equip }),
        });
        await refreshAll();
        await loadItems();
      } catch (e) { setError(e.message || String(e)); }
    });
  });

  root.querySelectorAll("[data-sell]").forEach(btn => {
    btn.addEventListener("click", async () => {
      const priceStr = prompt("Price?");
      const price = Number(priceStr || 0);
      if (!price || price <= 0) return;

      try {
        await api("/economy/market/items/sell", {
          method: "POST",
          body: JSON.stringify({
            kingdomId: state.kingdomId,
            sellerId: state.userId,
            itemId: btn.dataset.sell,
            price,
          }),
        });
        await loadItems();
        await refreshAll(); // later we’ll load item orders too
      } catch (e) { setError(e.message || String(e)); }
    });
  });

  el("itemsOut").textContent = JSON.stringify(items, null, 2);
}

async function loadItemOrders() {
  const data = await api(`/economy/market/items/orders?kingdomId=${encodeURIComponent(state.kingdomId)}`);
  renderItemOrders(data.orders || []);
}

function renderItemOrders(orders) {
  const root = el("itemOrders");
  root.innerHTML = "";

  if (!orders.length) {
    root.innerHTML = `<div class="muted">No item orders</div>`;
    el("itemOrdersOut").textContent = "[]";
    return;
  }

  for (const o of orders) {
    const isSeller = o.sellerId === state.userId;

    const div = document.createElement("div");
    div.className = "order";
    div.innerHTML = `
      <div class="top">
        <div class="id">order: ${o.orderId}</div>
        <div class="muted small">${new Date(o.createdAt).toLocaleString()}</div>
      </div>

      <div class="meta">
        <div><b>seller</b> ${o.sellerId}</div>
        <div><b>price</b> ${o.price}</div>
        <div><b>item</b> ${o.itemId}</div>
        <div><b>tier</b> ${o.tier}</div>
        <div><b>dur</b> ${o.durability}/${o.maxDurability}</div>
        <div><b>bonus</b> +${o.bonusPct}%</div>
      </div>

      <div class="actions">
        <button class="primary" data-buy="${o.orderId}" ${isSeller ? "disabled" : ""}>Buy</button>
        <button class="danger" data-cancel="${o.orderId}" ${isSeller ? "" : "disabled"}>Cancel</button>
      </div>
    `;

    root.appendChild(div);
  }

  // bind buy
  root.querySelectorAll("[data-buy]").forEach(btn => {
    btn.addEventListener("click", async () => {
      try {
        await api("/economy/market/items/buy", {
          method: "POST",
          body: JSON.stringify({
            orderId: btn.dataset.buy,
            buyerId: state.userId,
          }),
        });
        await loadItemOrders();
        await loadItems(); // чтобы увидеть item у покупателя/исчезновение listed у продавца
      } catch (e) {
        setError(e.message || String(e));
      }
    });
  });

  // bind cancel
  root.querySelectorAll("[data-cancel]").forEach(btn => {
    btn.addEventListener("click", async () => {
      try {
        await api("/economy/market/items/cancel", {
          method: "POST",
          body: JSON.stringify({
            orderId: btn.dataset.cancel,
            sellerId: state.userId,
          }),
        });
        await loadItemOrders();
        await loadItems(); // item снова станет listed=false
      } catch (e) {
        setError(e.message || String(e));
      }
    });
  });

  el("itemOrdersOut").textContent = JSON.stringify(orders, null, 2);
}
async function loadItemOrders() {
  const data = await api(`/economy/market/items/orders?kingdomId=${encodeURIComponent(state.kingdomId)}`);
  console.log("item orders data:", data);
  renderItemOrders(data.orders || []);
}


init();

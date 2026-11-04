"use strict";

/* Texts come from i18n.js: t(key, ...args). */
const LISTENERS = ["inventory", "notification"];
const DONE = new Set(["COMPLETED", "FAILED"]);
const SVG_NS = "http://www.w3.org/2000/svg";

/* Two layouts of the same board: side by side on wide screens, top to bottom on phones. */
const LAYOUTS = {
  wide: {
    viewBox: "0 0 1000 280",
    roles: "roles",
    nodes: {
      order: { x: 20, y: 100, w: 220, h: 84 },
      inventory: { x: 760, y: 30, w: 220, h: 84 },
      notification: { x: 760, y: 166, w: 220, h: 84 },
    },
    trunk: "M240 142 H470",
    branches: {
      inventory: "M470 142 L540 72 H760",
      notification: "M470 142 L540 208 H760",
    },
    label: { x: 355, y: 122, anchor: "middle" },
  },
  narrow: {
    viewBox: "0 0 360 540",
    roles: "rolesShort",
    nodes: {
      order: { x: 70, y: 10, w: 220, h: 84 },
      inventory: { x: 4, y: 430, w: 168, h: 84 },
      notification: { x: 188, y: 430, w: 168, h: 84 },
    },
    trunk: "M180 94 V230",
    branches: {
      inventory: "M180 230 L88 312 V430",
      notification: "M180 230 L272 312 V430",
    },
    label: { x: 196, y: 168, anchor: "start" },
  },
};

const $ = (id) => document.getElementById(id);
const reduceMotion = matchMedia("(prefers-reduced-motion: reduce)");
const narrow = matchMedia("(max-width: 640px)");

const state = {
  stock: {},
  quantities: {},
  nodes: { order: "idle", inventory: "idle", notification: "idle" },
  publications: [],
  seen: new Set(),
  logShown: false,
  busy: false,
};

/* ---------- API ---------- */

async function api(path, options = {}) {
  let response;
  try {
    response = await fetch(path, {
      headers: { "Content-Type": "application/json" },
      ...options,
    });
  } catch {
    throw new Error(t("noServer"));
  }
  if (!response.ok) {
    throw new Error(response.status === 400 ? t("badQuantity") : t("serverError", response.status));
  }
  const text = await response.text();
  return text ? JSON.parse(text) : null;
}

const loadStock = async () => renderProducts((state.stock = await api("/inventory")));
const loadEvents = () => api("/events");

/** Polls the event log until `isSettled` is true or about six seconds have passed. */
async function waitForEvents(isSettled) {
  let publications = [];
  for (let i = 0; i < 20; i++) {
    publications = await loadEvents();
    if (isSettled(publications)) break;
    await new Promise((resolve) => setTimeout(resolve, 300));
  }
  return publications;
}

/* ---------- Board ---------- */

function svg(tag, attributes = {}, parent) {
  const element = document.createElementNS(SVG_NS, tag);
  for (const [key, value] of Object.entries(attributes)) element.setAttribute(key, value);
  if (parent) parent.appendChild(element);
  return element;
}

/** Draws a track: a grey bed with a chain of lamps along it. */
function drawTrack(root, d, name) {
  const group = svg("g", { "data-track": name }, root);
  const path = svg("path", { d, class: "track-bed" }, group);
  const length = path.getTotalLength();
  const spacing = 22;
  for (let at = spacing / 2; at < length; at += spacing) {
    const point = path.getPointAtLength(at);
    svg("circle", { class: "lamp", cx: point.x.toFixed(1), cy: point.y.toFixed(1), r: 3.6 }, group);
  }
}

function drawNode(root, module, box, role) {
  const group = svg("g", { class: "node", "data-module": module, "data-state": state.nodes[module] }, root);
  svg("rect", { x: box.x, y: box.y, width: box.w, height: box.h, rx: 9 }, group);
  const name = svg("text", { class: "name", x: box.x + 18, y: box.y + 38 }, group);
  name.textContent = module;
  const roleText = svg("text", { class: "role", x: box.x + 18, y: box.y + 64 }, group);
  roleText.textContent = role;
  const cx = box.x + box.w - 22;
  const cy = box.y + 24;
  svg("circle", { class: "signal-glow", cx, cy, r: 11 }, group);
  svg("circle", { class: "signal", cx, cy, r: 8 }, group);
}

function renderBoard() {
  const layout = narrow.matches ? LAYOUTS.narrow : LAYOUTS.wide;
  const board = $("board");
  board.replaceChildren();
  const root = svg("svg", { viewBox: layout.viewBox, class: narrow.matches ? "is-narrow" : "", "aria-hidden": "true", focusable: "false" }, board);

  drawTrack(root, layout.trunk, "trunk");
  for (const module of LISTENERS) drawTrack(root, layout.branches[module], module);

  const label = svg("text", { class: "event-label", x: layout.label.x, y: layout.label.y, "text-anchor": layout.label.anchor }, root);
  label.textContent = "OrderPlaced";

  for (const [module, box] of Object.entries(layout.nodes)) drawNode(root, module, box, t(layout.roles)[module]);
}

function setNode(module, nodeState) {
  state.nodes[module] = nodeState;
  const node = document.querySelector(`.node[data-module="${module}"]`);
  if (node) node.dataset.state = nodeState;
}

/** Runs a light along the tracks, from order to the listeners that get the event. */
function runLight(targets = LISTENERS, fromTrunk = true) {
  if (reduceMotion.matches) return Promise.resolve();
  const css = getComputedStyle(document.documentElement);
  const off = css.getPropertyValue("--lamp-off").trim();
  const on = css.getPropertyValue("--lamp").trim();
  const step = 26;
  const flash = (lamp, delay) =>
    lamp.animate(
      [{ fill: off }, { fill: on, offset: 0.2 }, { fill: on, offset: 0.5 }, { fill: off }],
      { duration: 560, delay, easing: "ease-out" },
    ).finished;

  const runs = [];
  let offset = 0;
  if (fromTrunk) {
    const trunk = document.querySelectorAll('[data-track="trunk"] .lamp');
    trunk.forEach((lamp, i) => runs.push(flash(lamp, i * step)));
    offset = trunk.length * step;
  }
  for (const module of targets) {
    document.querySelectorAll(`[data-track="${module}"] .lamp`).forEach((lamp, i) => runs.push(flash(lamp, offset + i * step)));
  }
  return Promise.all(runs).catch(() => {}); // a redraw during the run cancels it, that is fine
}

/* ---------- Products ---------- */

const productName = (product) => t("productNames")[product] ?? product;
const BLOCKS_PER_GROUP = 10;

/** The order as the user has picked it so far: one line per product with a quantity above zero. */
const chosenLines = () =>
  Object.entries(state.quantities)
    .filter(([, quantity]) => quantity > 0)
    .map(([product, quantity]) => ({ product, quantity }));

const describeLines = (lines) => lines.map((line) => `${line.quantity} × ${productName(line.product)}`).join(", ");

/** Stock as blocks, one per unit, in groups of ten. The last `added` blocks are highlighted. */
function stockBlocks(count, added = 0) {
  const units = document.createElement("div");
  units.className = "units";
  units.setAttribute("aria-hidden", "true");
  if (count <= 0) {
    const none = document.createElement("span");
    none.className = "sold-out";
    none.textContent = t("soldOut");
    units.append(none);
    return units;
  }
  let group;
  for (let i = 0; i < count; i++) {
    if (i % BLOCKS_PER_GROUP === 0) {
      group = document.createElement("span");
      group.className = "unit-group";
      units.append(group);
    }
    const block = document.createElement("span");
    block.className = i >= count - added ? "unit is-new" : "unit";
    group.append(block);
  }
  return units;
}

function stepper(product) {
  const wrapper = document.createElement("div");
  wrapper.className = "stepper";
  const input = document.createElement("input");
  Object.assign(input, { type: "number", inputmode: "numeric", min: 0, max: 100, value: state.quantities[product] ?? 0 });
  input.setAttribute("aria-label", t("quantityOf", productName(product)));
  input.addEventListener("input", () => setQuantity(product, Number(input.value), false));

  const button = (delta, label, key) => {
    const b = document.createElement("button");
    b.type = "button";
    b.className = "step";
    b.textContent = label;
    b.setAttribute("aria-label", t(key, productName(product)));
    b.addEventListener("click", () => setQuantity(product, (state.quantities[product] ?? 0) + delta));
    return b;
  };
  wrapper.append(button(-1, "−", "less"), input, button(1, "+", "more"));
  return wrapper;
}

function setQuantity(product, value, redraw = true) {
  state.quantities[product] = Number.isFinite(value) ? Math.min(100, Math.max(0, Math.round(value))) : 0;
  if (redraw) renderProducts(state.stock);
  else document.querySelector(`.product[data-product="${product}"]`)?.classList.toggle("is-chosen", state.quantities[product] > 0);
  renderSummary();
}

function renderProducts(stock, restocked = {}) {
  const list = $("products");
  const focused = document.activeElement?.closest?.(".product") ? document.activeElement : null;
  const focusKey = focused && [focused.closest(".product").dataset.product, focused.getAttribute("aria-label")];
  list.replaceChildren();

  for (const [product, count] of Object.entries(stock)) {
    const item = document.createElement("div");
    item.className = "product";
    item.dataset.product = product;
    item.classList.toggle("is-chosen", (state.quantities[product] ?? 0) > 0);

    const heading = document.createElement("div");
    heading.className = "product-heading";
    const name = document.createElement("span");
    name.className = "product-name";
    name.textContent = productName(product);
    const amount = document.createElement("span");
    amount.className = "product-count";
    const strong = document.createElement("strong");
    strong.textContent = count;
    amount.append(strong, t("inStock"));
    heading.append(name, amount);

    const restock = document.createElement("button");
    restock.type = "button";
    restock.className = "restock";
    restock.textContent = t("restock");
    restock.setAttribute("aria-label", t("restockAria", productName(product)));
    restock.addEventListener("click", () => restockProduct(product));

    item.append(heading, stepper(product), stockBlocks(count, restocked.product === product ? restocked.added : 0), restock);
    list.append(item);
  }

  // Redrawing replaces the buttons, so give focus back to the same control.
  if (focusKey) {
    list.querySelector(`.product[data-product="${focusKey[0]}"] [aria-label="${CSS.escape(focusKey[1])}"]`)?.focus();
  }
}

function renderSummary() {
  const lines = chosenLines();
  $("summary").textContent = lines.length ? describeLines(lines) : t("summaryEmpty");
  $("summary").classList.toggle("is-empty", lines.length === 0);
}

async function restockProduct(product) {
  try {
    state.stock = await api("/inventory/restock", { method: "POST", body: JSON.stringify({ product, quantity: 10 }) });
    renderProducts(state.stock, { product, added: 10 });
    say(t("restocked", productName(product), state.stock[product]));
  } catch (error) {
    say(error.message, true);
  }
}

/* ---------- Event log ---------- */

/** Java and Go write up to nine fractional digits. Date only needs three. */
const toDate = (iso) => new Date(iso.replace(/(\.\d{3})\d+/, "$1"));
const shortId = (id) => id.slice(0, 8);
const outcome = (status) => (status === "COMPLETED" ? "ok" : status === "FAILED" ? "fail" : "wait");

function describe(publication) {
  if (publication.status === "COMPLETED") {
    if (publication.attempts > 1) return t("processedOnAttempt", publication.attempts);
    const ms = toDate(publication.completedAt) - toDate(publication.publishedAt);
    return t("processedAfter", Math.max(ms, 1));
  }
  if (publication.status === "FAILED") {
    return publication.attempts > 1 ? t("failedAttempts", publication.attempts) : t("failed");
  }
  return t("processing");
}

/** One entry per event, with one line per listener that received it. */
function groupByOrder(publications) {
  const groups = new Map();
  for (const publication of publications) {
    const key = publication.event.orderId ?? publication.id;
    if (!groups.has(key)) groups.set(key, { event: publication.event, type: publication.eventType, publishedAt: publication.publishedAt, deliveries: [] });
    groups.get(key).deliveries.push(publication);
  }
  for (const group of groups.values()) {
    group.deliveries.sort((a, b) => LISTENERS.indexOf(a.module) - LISTENERS.indexOf(b.module));
  }
  return [...groups.entries()];
}

function renderLog(publications) {
  state.publications = publications;
  const list = $("entries");
  list.replaceChildren();
  $("empty").hidden = publications.length > 0;

  for (const [key, group] of groupByOrder(publications)) {
    const entry = document.createElement("li");
    entry.className = "entry";
    if (state.logShown && !state.seen.has(key)) entry.classList.add("is-new");

    const what = document.createElement("span");
    what.className = "entry-what";
    what.textContent = describeLines(group.event.lines ?? []);

    const when = document.createElement("time");
    when.className = "entry-when";
    when.dateTime = group.publishedAt;
    when.textContent = toDate(group.publishedAt).toLocaleTimeString(t("locale"));

    const id = document.createElement("span");
    id.className = "entry-id";
    id.textContent = t("entryId", group.type, shortId(key));

    const deliveries = document.createElement("ul");
    deliveries.className = "deliveries";
    for (const publication of group.deliveries) {
      const line = document.createElement("li");
      line.className = `delivery ${outcome(publication.status)}`;
      const module = document.createElement("code");
      module.textContent = publication.module;
      line.append(module, describe(publication));
      deliveries.append(line);
    }

    entry.append(what, when, id, deliveries);
    list.append(entry);
  }
  for (const [key] of groupByOrder(publications)) state.seen.add(key);
  state.logShown = true;

  renderFailures(publications.filter((p) => p.status === "FAILED"));
}

function renderFailures(failed) {
  $("failures").hidden = failed.length === 0;
  if (failed.length === 0) return;
  const modules = [...new Set(failed.map((p) => p.module))].join(t("and"));
  $("failures-text").textContent = t("failures", failed.length, modules);
}

/* ---------- Actions ---------- */

function say(text, isError = false) {
  const message = $("message");
  message.textContent = text;
  message.classList.toggle("error", isError);
}

function setBusy(busy) {
  state.busy = busy;
  $("submit").disabled = busy;
  $("resubmit").disabled = busy;
}

async function placeOrder(event) {
  event.preventDefault();
  if (state.busy) return;

  const lines = chosenLines();
  if (lines.length === 0) {
    say(t("nothingChosen"), true);
    return;
  }

  setBusy(true);
  LISTENERS.forEach((module) => setNode(module, "idle"));
  setNode("order", "wait");

  try {
    const { orderId } = await api("/orders", {
      method: "POST",
      body: JSON.stringify({ lines }),
    });
    setNode("order", "ok");
    say(t("placed", shortId(orderId), describeLines(lines)));
    state.quantities = {};
    renderProducts(state.stock);
    renderSummary();

    const ours = (publications) => publications.filter((p) => p.event.orderId === orderId);
    const [, publications] = await Promise.all([
      runLight(),
      waitForEvents((all) => ours(all).length >= LISTENERS.length && ours(all).every((p) => DONE.has(p.status))),
    ]);

    for (const publication of ours(publications)) setNode(publication.module, outcome(publication.status));
    renderLog(publications);
    await loadStock();
  } catch (error) {
    setNode("order", "fail");
    say(error.message, true);
  } finally {
    setBusy(false);
  }
}

async function resubmit() {
  if (state.busy) return;
  setBusy(true);
  try {
    const before = await loadEvents();
    const failed = before.filter((p) => p.status === "FAILED");
    const modules = [...new Set(failed.map((p) => p.module))];
    modules.forEach((module) => setNode(module, "wait"));

    await api("/events/resubmit", { method: "POST" });
    const ids = new Set(failed.map((p) => p.id));
    const [, publications] = await Promise.all([
      runLight(modules, false),
      waitForEvents((all) => all.filter((p) => ids.has(p.id)).every((p) => DONE.has(p.status) && p.attempts > failedAttempts(failed, p.id))),
    ]);

    const retried = publications.filter((p) => ids.has(p.id));
    for (const publication of retried) setNode(publication.module, outcome(publication.status));
    const stillFailed = retried.filter((p) => p.status === "FAILED").length;
    say(stillFailed === 0 ? t("resubmitted") : t("resubmittedFailed"), stillFailed > 0);

    renderLog(publications);
    await loadStock();
  } catch (error) {
    say(error.message, true);
  } finally {
    setBusy(false);
  }
}

const failedAttempts = (failed, id) => failed.find((p) => p.id === id)?.attempts ?? 0;

/* ---------- Start ---------- */

/** Redraws everything that was built from texts in JavaScript. */
function onLanguageChange() {
  renderBoard();
  renderProducts(state.stock);
  renderSummary();
  renderLog(state.publications);
  say("");
}

async function start() {
  applyStaticTexts();
  renderBoard();
  document.addEventListener("languagechange", onLanguageChange);
  document.querySelectorAll("[data-language]").forEach((button) => button.addEventListener("click", () => setLanguage(button.dataset.language)));
  narrow.addEventListener("change", renderBoard);
  $("order-form").addEventListener("submit", placeOrder);
  $("resubmit").addEventListener("click", resubmit);

  try {
    await loadStock();
    renderSummary();
    renderLog(await loadEvents());
  } catch (error) {
    say(error.message, true);
  }
}

start();

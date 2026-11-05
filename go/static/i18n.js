"use strict";

/*
 * All visible text of the page, in German and English.
 * Keys used in index.html: data-i18n (text), data-i18n-html (text with <code>), data-i18n-aria (aria-label).
 * The HTML strings below are fixed text from this file, never data from the server.
 */
const TEXTS = {
  de: {
    locale: "de-DE",
    title: "Northbean Rösterei",
    language: "Sprache",
    intro:
      "Der Onlineshop von Northbean, einer kleinen Kaffeerösterei, aufgeteilt in drei Module. <code>order</code> kennt die anderen nicht, es meldet nur das Event " +
      "<code>OrderPlaced</code>. <code>inventory</code> und <code>notification</code> hören darauf, jedes für sich.",
    boardTitle: "Weg des Events durch die Module",
    boardAria: "order sendet OrderPlaced an inventory und notification",
    legendAria: "Bedeutung der Lampen",
    legendWait: "wird verarbeitet",
    legendOk: "verarbeitet",
    legendFail: "fehlgeschlagen, Event bleibt gespeichert",
    counterTitle: "Bestellen",
    products: "Produkte",
    quantityOf: (product) => `Menge ${product}`,
    less: (product) => `Eins weniger ${product}`,
    more: (product) => `Eins mehr ${product}`,
    summaryEmpty: "Noch nichts ausgewählt. Wähle bei jedem Produkt die Menge, die du bestellen willst.",
    nothingChosen: "Wähle mindestens ein Produkt.",
    submit: "Bestellung aufgeben",
    hint:
      "Bestellst du mehr, als im Lager ist, schlägt <code>inventory</code> für die ganze Bestellung fehl. <code>notification</code> " +
      "bestätigt trotzdem, denn jedes Modul verarbeitet das Event für sich.",
    logTitle: "Event-Protokoll",
    logIntro:
      "Jedes Event wird einmal pro Listener gespeichert und abgehakt, sobald der Listener fertig ist. " +
      "In Kotlin macht das Spring Modulith (Tabelle <code>EVENT_PUBLICATION</code>), in Go <code>internal/events</code>.",
    resubmit: "Erneut zustellen",
    empty: "Noch keine Events. Gib eine Bestellung auf, dann siehst du hier, welche Module sie verarbeitet haben.",

    productNames: { "crema-beans": "Crema-Bohnen", "espresso-beans": "Espresso-Bohnen", "espresso-cups": "Espressotassen" },
    roles: { order: "nimmt Bestellungen an", inventory: "zieht den Bestand ab", notification: "bestätigt die Bestellung" },
    rolesShort: { order: "nimmt Bestellungen an", inventory: "zieht Bestand ab", notification: "bestätigt" },
    inStock: " im Lager",
    restock: "Auffüllen (+10)",
    restockAria: (product) => `${product} auffüllen, plus 10`,
    soldOut: "ausverkauft",
    restocked: (product, count) => `${product} aufgefüllt, jetzt ${count} im Lager.`,
    processedOnAttempt: (n) => `verarbeitet im ${n}. Versuch`,
    processedAfter: (ms) => `verarbeitet nach ${ms} ms`,
    failed: "fehlgeschlagen",
    failedAttempts: (n) => `fehlgeschlagen, ${n} Versuche`,
    processing: "wird verarbeitet",
    entryId: (type, id) => `${type} für Bestellung ${id}`,
    and: " und ",
    failures: (n, modules) =>
      `${n === 1 ? "Eine Zustellung" : `${n} Zustellungen`} an ${modules} ${n === 1 ? "ist" : "sind"} fehlgeschlagen. ` +
      "Das Event bleibt gespeichert, bis der Listener es verarbeitet hat. Fülle das Lager auf und stelle es dann erneut zu.",
    noServer: "Der Server antwortet nicht. Läuft ./gradlew bootRun (Kotlin) oder go run . (Go)?",
    badQuantity: "Die Menge muss zwischen 1 und 100 liegen.",
    serverError: (status) => `Der Server meldet Fehler ${status}.`,
    placed: (id, summary) => `Bestellung ${id} aufgegeben: ${summary}.`,
    resubmitted: "Erneut zugestellt. Alle Listener haben das Event verarbeitet.",
    resubmittedFailed: "Erneut zugestellt, aber der Listener schlägt wieder fehl. Ist genug im Lager?",
  },

  en: {
    locale: "en-GB",
    title: "Northbean Roastery",
    language: "Language",
    intro:
      "The online shop of Northbean, a small coffee roastery, split into three modules. <code>order</code> does not know the others, it only announces the event " +
      "<code>OrderPlaced</code>. <code>inventory</code> and <code>notification</code> listen for it, each on its own.",
    boardTitle: "How the event travels through the modules",
    boardAria: "order sends OrderPlaced to inventory and notification",
    legendAria: "What the lamps mean",
    legendWait: "processing",
    legendOk: "processed",
    legendFail: "failed, the event stays stored",
    counterTitle: "Order",
    products: "Products",
    quantityOf: (product) => `Quantity of ${product}`,
    less: (product) => `One less ${product}`,
    more: (product) => `One more ${product}`,
    summaryEmpty: "Nothing chosen yet. Set how many of each product you want to order.",
    nothingChosen: "Choose at least one product.",
    submit: "Place order",
    hint:
      "Order more than is in stock and <code>inventory</code> fails for the whole order. <code>notification</code> still confirms, " +
      "because each module handles the event on its own.",
    logTitle: "Event log",
    logIntro:
      "Each event is stored once per listener and marked done as soon as the listener has finished. " +
      "In Kotlin, Spring Modulith does this (table <code>EVENT_PUBLICATION</code>), in Go <code>internal/events</code>.",
    resubmit: "Resubmit",
    empty: "No events yet. Place an order to see which modules handled it.",

    productNames: { "crema-beans": "Crema beans", "espresso-beans": "Espresso beans", "espresso-cups": "Espresso cups" },
    roles: { order: "accepts orders", inventory: "reduces the stock", notification: "confirms the order" },
    rolesShort: { order: "accepts orders", inventory: "reduces stock", notification: "confirms" },
    inStock: " in stock",
    restock: "Restock (+10)",
    restockAria: (product) => `Restock ${product}, plus 10`,
    soldOut: "sold out",
    restocked: (product, count) => `${product} restocked, now ${count} in stock.`,
    processedOnAttempt: (n) => `processed on attempt ${n}`,
    processedAfter: (ms) => `processed after ${ms} ms`,
    failed: "failed",
    failedAttempts: (n) => `failed, ${n} attempts`,
    processing: "processing",
    entryId: (type, id) => `${type} for order ${id}`,
    and: " and ",
    failures: (n, modules) =>
      `${n === 1 ? "One delivery" : `${n} deliveries`} to ${modules} failed. ` +
      "The event stays stored until the listener has handled it. Restock the product, then resubmit the event.",
    noServer: "The server does not respond. Is ./gradlew bootRun (Kotlin) or go run . (Go) running?",
    badQuantity: "Quantity must be between 1 and 100.",
    serverError: (status) => `The server reports error ${status}.`,
    placed: (id, summary) => `Order ${id} placed: ${summary}.`,
    resubmitted: "Resubmitted. All listeners have processed the event.",
    resubmittedFailed: "Resubmitted, but the listener failed again. Is there enough in stock?",
  },
};

const LANGUAGES = Object.keys(TEXTS);
const STORAGE_KEY = "modulith-roastery-language";

function initialLanguage() {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (LANGUAGES.includes(stored)) return stored;
  return "en";
}

let language = initialLanguage();

/** Looks up a text. Functions are called with the given arguments. */
function t(key, ...args) {
  const value = TEXTS[language][key];
  return typeof value === "function" ? value(...args) : value;
}

/** Writes the texts of the current language into the static parts of the page. */
function applyStaticTexts() {
  document.documentElement.lang = language;
  document.title = t("title");
  document.querySelectorAll("[data-i18n]").forEach((el) => (el.textContent = t(el.dataset.i18n)));
  document.querySelectorAll("[data-i18n-html]").forEach((el) => (el.innerHTML = t(el.dataset.i18nHtml)));
  document.querySelectorAll("[data-i18n-aria]").forEach((el) => el.setAttribute("aria-label", t(el.dataset.i18nAria)));
  document.querySelectorAll("[data-language]").forEach((button) => button.setAttribute("aria-pressed", String(button.dataset.language === language)));
}

function setLanguage(next) {
  if (!LANGUAGES.includes(next) || next === language) return;
  language = next;
  localStorage.setItem(STORAGE_KEY, next);
  applyStaticTexts();
  document.dispatchEvent(new CustomEvent("languagechange"));
}

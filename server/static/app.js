function licenseplate(el) {
  const input = el.querySelector("input");
  if (!input) return;

  const parentForm = (() => {
    for (var form of document.querySelectorAll("form")) {
      let current = el;
      while (current != null && form != current) {
        current = current.parentNode;
      }
    }
    return form;
  })();

  const refresh = () => {
    // Make the width of the input match the text
    if (input.value == "" && input.hasAttribute("placeholder")) {
      input.setAttribute("size", input.getAttribute("placeholder").length);
    } else {
      input.setAttribute("size", input.value.length);
    }

    if (parentForm)
      parentForm.action = "/plates/" + input.value.trim().toUpperCase();
  };

  const maxlen = parseInt(input.getAttribute("maxlength"));
  const valid = /^([0-9a-zA-Z ]|Enter)$/;
  const validate = (e) => {
    // Only let in chars we want
    if (!e.key.match(valid)) {
      e.preventDefault();
      return;
    }

    // Don't allow leading spaces
    if (e.target.value == "" && e.key == " ") {
      e.preventDefault();
      return;
    }

    // Don't allow double spaces
    if (e.target.value[e.target.value.length - 1] == " " && e.key == " ") {
      e.preventDefault();
      return;
    }

    // Fix the length insertion we break below
    if (e.target.value.length >= maxlen) return;
  };

  const betterSubmit = (e) => {
    console.log("betterSubmit", { parentForm });
    if (e) e.preventDefault();

    if (parentForm.checkValidity()) {
      visit(parentForm.action);
    } else {
      parentForm.reportValidity();
    }
  };

  refresh();
  input.addEventListener("input", refresh);
  input.addEventListener("keypress", validate);
  if (parentForm) {
    parentForm.addEventListener("submit", betterSubmit);
  }
}

function autoexpand(el) {
  let textarea = el.querySelector("textarea");
  if (!textarea) return;

  let refresh = () => {
    textarea.style.height = 0;
    if (textarea.scrollHeight > 50) {
      textarea.style.height = textarea.scrollHeight + "px";
    }
  };

  refresh();
  textarea.addEventListener("input", refresh);
}

// Given a string of HTML, replace the current page with it
function replacePage(nextPage) {
  const container = document.querySelector("#container");
  const newDoc = new DOMParser().parseFromString(nextPage, "text/html");
  const newContainer = newDoc.querySelector("#container");
  container.innerHTML = newContainer.innerHTML;
  document.title = newDoc.title;
  onPageLoad();
}

// We do have a service worker for an offline cache, but we should also keep an
// in-memory cache of visited pages.
const pageCache = {};

async function visit(url, { addHistory, resetScroll } = {}) {
  if (url == "") url = "/";
  if (addHistory == undefined) addHistory = true;
  if (resetScroll == undefined) resetScroll = true;

  console.log("Visiting", {
    url,
    addHistory,
    resetScroll,
    history: history.length,
  });

  // If we don't have the page in our cache, add it.
  if (!(url in pageCache)) {
    pageCache[url] = await fetch(url).then((r) => r.text());
  }

  // If we're here, we the page. Use it.
  if (addHistory) {
    history.pushState({}, "", url);
  }
  replacePage(pageCache[url]);

  // Scroll to the top of the page
  if (resetScroll) window.scrollTo(0, 0);
}

window.visit = visit;

function ajaxLink(el) {
  // Only handle links on this domain
  if (el.host != window.location.host) return;

  el.addEventListener("click", (e) => {
    e.preventDefault();
    visit(el.href);
  });
}

function backButton() {
  const button = document.querySelector("#back");

  // If the button is missing, nothing to do
  if (!button) return;

  // Nowhere to go back from the homepage
  if (location.pathname == "/") {
    button.parentNode.removeChild(button);
  }

  button.addEventListener("click", (e) => {
    e.preventDefault();

    if (history.length > 1) {
      // Go back in the linear history
      history.back();
    } else {
      // Go up the URL tree
      console.log("Going up the URL tree");
      visit(window.location.pathname.split("/").slice(0, -1).join("/"));
    }
  });
}

function onPageLoad() {
  backButton();
  licenseplate(document.querySelector(".license-plate"));

  for (var input of document.querySelectorAll(".input")) {
    autoexpand(input);
  }

  for (var link of document.querySelectorAll("a")) {
    ajaxLink(link);
  }
}

function startup() {
  // Remove debug statements outside dev
  if (window.location.hostname !== "localhost") {
    console.log = () => {};
  }

  // Turn on the service worker
  if (navigator && navigator.serviceWorker) {
    navigator.serviceWorker.register("/service.js");
  }

  window.addEventListener("popstate", () => {
    visit(window.location.href, { addHistory: false, resetScroll: false });
  });

  onPageLoad();
}

startup();

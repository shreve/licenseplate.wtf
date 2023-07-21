function licenseplate(el) {
  let input = el.querySelector("input");
  if (!input) return;

  let parentForm = (() => {
    for (var form of document.querySelectorAll("form")) {
      let current = el;
      while (current != null && form != current) {
        current = current.parentNode;
      }
    }
    return form;
  })();

  let refresh = () => {
    // Make the width of the input match the text
    if (input.value == "" && input.hasAttribute("placeholder")) {
      input.setAttribute("size", input.getAttribute("placeholder").length);
    } else {
      input.setAttribute("size", input.value.length);
    }

    if (parentForm)
      parentForm.action = "/plates/" + input.value.trim().toUpperCase();
  };

  let maxlen = parseInt(input.getAttribute("maxlength"));
  let valid = /^[0-9a-zA-Z ]$/;
  let validate = (e) => {
    if (e.key == "Enter") {
      betterSubmit();
      return;
    }

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

  let betterSubmit = (e) => {
    if (e) e.preventDefault();

    if (parentForm.checkValidity()) window.location = parentForm.action;
    else parentForm.reportValidity();
  };

  refresh();
  input.addEventListener("input", refresh);
  input.addEventListener("keypress", validate);
  if (parentForm) parentForm.addEventListener("submit", betterSubmit);
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
  input.addEventListener("input", refresh);
}

licenseplate(document.querySelector(".license-plate"));
for (var input of document.querySelectorAll(".input")) {
  autoexpand(input);
}

function licenseplate(el) {
  let input = el.querySelector("input");
  if (!input)
    return;
  let initial = input.scrollWidth + 5;

  let parentForm = (() => {
    for (var form of document.querySelectorAll('form')) {
      let current = el;
      while (current != null && form != current) {
        current = current.parentNode
      }
    }
    return form
  })();

  let refresh = () => {
    input.style.width = 0;
    input.style.width = Math.max(initial, input.scrollWidth) + 'px';
    if (parentForm)
      parentForm.action = "/plates/" + input.value.trim().toUpperCase();
  }

  let maxlen = parseInt(input.getAttribute('maxlength'));
  let valid = /^[0-9a-zA-Z ]$/;
  let validate = (e) => {
    if (e.key == "Enter") {
      betterSubmit();
      return
    }

    // Only let in chars we want
    if (!e.key.match(valid)) {
      e.preventDefault();
      return
    }

    // Don't allow leading spaces
    if (e.target.value == "" && e.key == " ") {
      e.preventDefault();
      return
    }

    // Don't allow double spaces
    if (e.target.value[e.target.value.length-1] == " " && e.key == " ") {
      e.preventDefault();
      return
    }

    // Fix the length insertion we break below
    if (e.target.value.length >= maxlen)
      return
  }

  let betterSubmit = (e) => {
    if (e)
      e.preventDefault()

    if (parentForm.checkValidity())
      window.location = parentForm.action;
    else
      parentForm.reportValidity()
  }

  refresh();
  input.addEventListener('input', refresh)
  input.addEventListener('keypress', validate)
  if (parentForm)
    parentForm.addEventListener('submit', betterSubmit)
}

function autoexpand(el) {
  el.setAttribute('data-baseheight', el.offsetHeight);
  let textarea = el.querySelector('textarea');
  if (!textarea) return;

  let refresh = () => {
    textarea.style.height = 0
    textarea.style.height = Math.max(50, textarea.scrollHeight) + 'px'
  }

  refresh()
  input.addEventListener('input', refresh)
}

licenseplate(document.querySelector(".license-plate"));
for (var input of document.querySelectorAll('.input')) {
  autoexpand(input);
}

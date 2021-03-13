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
    // input.size = size(input);
    input.style.width = 0;
    input.style.width = Math.max(initial, input.scrollWidth) + 'px';
    if (parentForm)
      parentForm.action = "/plates/" + input.value;
  }

  refresh();
  input.addEventListener('input', refresh)
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

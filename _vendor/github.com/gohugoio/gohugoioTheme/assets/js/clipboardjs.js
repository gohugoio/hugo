var Clipboard = require('clipboard/dist/clipboard.js');
new Clipboard('.copy', {
  target: function(trigger) {
    if(trigger.classList.contains('copy-toggle')){
      return trigger.previousElementSibling;
    }
    return trigger.nextElementSibling;
  }
  }).on('success', function(e) {
    successMessage(e.trigger, 'Copied!');
    e.clearSelection();
  }).on('error', function(e) {
    successMessage(e.trigger, fallbackMessage(e.action));
});

function successMessage(elem, msg) {
  elem.setAttribute('class', 'copied bg-primary-color-dark f6 absolute top-0 right-0 lh-solid hover-bg-primary-color-dark bn white ph3 pv2');
  elem.setAttribute('aria-label', msg);
}

function fallbackMessage(elem, action) {
  var actionMsg = '';
  var actionKey = (action === 'cut' ? 'X' : 'C');
  if (isMac) {
      actionMsg = 'Press ⌘-' + actionKey;
  } else {
      actionMsg = 'Press Ctrl-' + actionKey;
  }
  return actionMsg;
}

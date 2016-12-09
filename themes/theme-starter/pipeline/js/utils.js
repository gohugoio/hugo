//Utility functions (when not using jQuery)

let byId = function(id) {
  return document.getElementById(id);
};

let qs = function(sel) {
  return document.querySelector(sel);
};

let qsa = function(sel) {
  return document.querySelectorAll(sel);
}

//browser sniff
function getMobileOperatingSystem() {
  var userAgent = navigator.userAgent || navigator.vendor || window.opera;
  // Windows Phone must come first because its UA also contains "Android"
  if (/windows phone/i.test(userAgent)) {
    return "Windows Phone";
  }
  if (/android/i.test(userAgent)) {
    return "Android";
  }
  // iOS detection from: http://stackoverflow.com/a/9039885/177710
  if (/iPad|iPhone|iPod/.test(userAgent) && !window.MSStream) {
    return "iOS";
    return "notMobile";
  }
}

function urlize(item) {
  if (typeof item == "string") {
    return item.replace(/[^\w\s\-]/gi, '').toLowerCase().replace(' ', '-');
  } else if (item instanceof HTMLElement) {
    return item.textContent.replace(/[^\w\s\-]/gi, '').toLowerCase().replace(' ', '-');
  }
}

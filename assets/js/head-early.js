import { scrollToActive } from 'js/helpers/index';

(function () {
	// Now we know that the browser has JS enabled.
	document.documentElement.classList.remove('no-js');

	// Add os-macos class to body if user is using macOS.
	if (navigator.userAgent.indexOf('Mac') > -1) {
		document.documentElement.classList.add('os-macos');
	}

	// Wait for the DOM to be ready.
	document.addEventListener('DOMContentLoaded', function () {
		scrollToActive('DOMContentLoaded');
	});
})();

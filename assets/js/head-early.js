import { scrollToActive } from 'js/helpers/index';

(function () {
	// Now we know that the browser has JS enabled.
	document.documentElement.classList.remove('no-js');

	// Wait for the DOM to be ready.
	document.addEventListener('DOMContentLoaded', function () {
		scrollToActive('DOMContentLoaded');
	});
})();

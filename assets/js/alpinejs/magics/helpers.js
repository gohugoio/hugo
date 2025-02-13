'use strict';

export function registerMagics(Alpine) {
	Alpine.magic('copy', (currentEl) => {
		return function (el) {
			if (!el) {
				el = currentEl;
			}

			// Select the element to copy.
			let range = document.createRange();
			range.selectNode(el);
			window.getSelection().removeAllRanges();
			window.getSelection().addRange(range);

			// Remove the selection after some time.
			setTimeout(() => {
				window.getSelection().removeAllRanges();
			}, 500);

			// Trim whitespace.
			let text = el.textContent.trim();

			navigator.clipboard.writeText(text);
		};
	});

	Alpine.magic('isScrollX', (currentEl) => {
		return function (el) {
			if (!el) {
				el = currentEl;
			}
			return el.clientWidth < el.scrollWidth;
		};
	});
}

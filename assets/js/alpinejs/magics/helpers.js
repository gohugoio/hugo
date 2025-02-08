'use strict';

export function registerMagics(Alpine) {
	Alpine.magic('copy', (currentEl) => {
		return function (el) {
			if (!el) {
				el = currentEl;
			}
			let lntds = el.querySelectorAll('.lntable .lntd');
			if (lntds && lntds.length === 2) {
				el = lntds[1];
			}

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

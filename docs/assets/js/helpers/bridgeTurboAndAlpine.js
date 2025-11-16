export function bridgeTurboAndAlpine(Alpine) {
	document.addEventListener('turbo:before-render', (event) => {
		event.detail.newBody.querySelectorAll('[data-alpine-generated]').forEach((el) => {
			if (el.hasAttribute('data-alpine-generated')) {
				el.removeAttribute('data-alpine-generated');
				el.remove();
			}
		});
	});

	document.addEventListener('turbo:render', () => {
		if (document.documentElement.hasAttribute('data-turbo-preview')) {
			return;
		}

		document.querySelectorAll('[data-alpine-ignored]').forEach((el) => {
			el.removeAttribute('x-ignore');
			el.removeAttribute('data-alpine-ignored');
		});

		document.body.querySelectorAll('[x-data]').forEach((el) => {
			if (el.hasAttribute('data-turbo-permanent')) {
				return;
			}
			Alpine.initTree(el);
		});

		Alpine.startObservingMutations();
	});

	// Cleanup Alpine state on navigation.
	document.addEventListener('turbo:before-cache', () => {
		// This will be restarted in turbo:render.
		Alpine.stopObservingMutations();

		document.body.querySelectorAll('[data-turbo-permanent]').forEach((el) => {
			if (!el.hasAttribute('x-ignore')) {
				el.setAttribute('x-ignore', true);
				el.setAttribute('data-alpine-ignored', true);
			}
		});

		document.body.querySelectorAll('[x-for],[x-if],[x-teleport]').forEach((el) => {
			if (el.hasAttribute('x-for') && el._x_lookup) {
				Object.values(el._x_lookup).forEach((el) => el.setAttribute('data-alpine-generated', true));
			}

			if (el.hasAttribute('x-if') && el._x_currentIfEl) {
				el._x_currentIfEl.setAttribute('data-alpine-generated', true);
			}

			if (el.hasAttribute('x-teleport') && el._x_teleport) {
				el._x_teleport.setAttribute('data-alpine-generated', true);
			}
		});

		document.body.querySelectorAll('[x-data]').forEach((el) => {
			if (!el.hasAttribute('data-turbo-permanent')) {
				Alpine.destroyTree(el);
				// Turbo leaks DOM elements via their data-turbo-permanent handling.
				// That needs to be fixed upstream, but until then.
				let clone = el.cloneNode(true);
				el.replaceWith(clone);
			}
		});
	});
}

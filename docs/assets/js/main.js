import Alpine from 'alpinejs';
import { registerMagics } from './alpinejs/magics/index';
import { navbar, search, toc } from './alpinejs/data/index';
import { navStore, initColorScheme } from './alpinejs/stores/index';
import { bridgeTurboAndAlpine } from './helpers/index';
import persist from '@alpinejs/persist';
import focus from '@alpinejs/focus';

// ---------------------------------------------------------
// PREVENT MULTIPLE INITIALIZATIONS (important for Turbo)
// ---------------------------------------------------------
if (!window.__ALPINE_ALREADY_INIT__) {
	window.__ALPINE_ALREADY_INIT__ = true;

	// Alpine plugins
	Alpine.plugin(persist);
	Alpine.plugin(focus);

	// Magic functions (copy, tooltip, etc.)
	registerMagics(Alpine);

	// Alpine data components
	Alpine.data('navbar', () => navbar(Alpine));

	// DON'T hardcode keys directly in source code
	const searchConfig = {
		index: 'hugodocs',
		app_id: import.meta.env.VITE_SEARCH_APP_ID,
		api_key: import.meta.env.VITE_SEARCH_API_KEY,
	};

	Alpine.data('search', () => search(Alpine, searchConfig));
	Alpine.data('toc', () => toc(Alpine));

	// Stores
	Alpine.store('nav', navStore(Alpine));

	// Connect Turbo + Alpine BEFORE start
	bridgeTurboAndAlpine(Alpine);

	Alpine.start();
}

// ---------------------------------------------------------
// COLOR SCHEME FIX FOR TURBO PAGE LOADS
// ---------------------------------------------------------
document.addEventListener('turbo:render', () => {
	initColorScheme();
});

// ---------------------------------------------------------
// ROBUST SCROLL PRESERVATION (supports multiple containers)
// ---------------------------------------------------------
(() => {
	let scrollCache = new Map();

	// Store scroll positions before navigation
	addEventListener('turbo:before-visit', () => {
		document.querySelectorAll('[data-turbo-preserve-scroll-container]').forEach((el) => {
			scrollCache.set(el.dataset.turboPreserveScrollContainer, el.scrollTop);
		});
	});

	// Restore scroll positions after navigation
	addEventListener('turbo:render', () => {
		document.querySelectorAll('[data-turbo-preserve-scroll-container]').forEach((el) => {
			const key = el.dataset.turboPreserveScrollContainer;

			if (scrollCache.has(key)) {
				el.scrollTop = scrollCache.get(key);
				return;
			}

			// Fallback: scroll to first element marked as active
			const active = el.querySelector('.scroll-active');
			if (active) {
				const activePos = active.offsetTop - el.offsetTop;
				el.scrollTop = Math.max(activePos, 0);
			}
		});

		scrollCache.clear();
	});
})();

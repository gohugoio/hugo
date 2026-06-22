import Alpine from 'alpinejs';
import { registerMagics } from './alpinejs/magics/index';
import { navbar, search, toc } from './alpinejs/data/index';
import { navStore } from './alpinejs/stores/index';
import persist from '@alpinejs/persist';
import focus from '@alpinejs/focus';
import * as params from '@params';

var debug = 0 ? console.log.bind(console, '[index]') : function () {};

// Set up and start Alpine.
(function () {
	// Register AlpineJS plugins.
	{
		Alpine.plugin(focus);
		Alpine.plugin(persist);
	}
	// Register AlpineJS magics and directives.
	{
		// Handles copy to clipboard etc.
		registerMagics(Alpine);
	}

	// Register AlpineJS controllers.
	{
		// Register AlpineJS data controllers.
		let searchConfig = {
			index: 'hugodocs',
			app_id: 'D1BPLZHGYQ',
			api_key: '6df94e1e5d55d258c56f60d974d10314',
			params: params,
		};

		Alpine.data('navbar', () => navbar(Alpine));
		Alpine.data('search', () => search(Alpine, searchConfig));
		Alpine.data('toc', () => toc(Alpine));
	}

	// Register AlpineJS stores.
	{
		Alpine.store('nav', navStore(Alpine));
	}

	// Start AlpineJS.
	Alpine.start();

	// On cross-document navigation the browser snapshots the current page for
	// the view transition. An open overlay (e.g. the search modal) would
	// otherwise linger in that outgoing snapshot while the page crossfades.
	// `pageswap` runs right before the snapshot is taken, so hide such
	// elements here to make them disappear instantly on navigation.
	window.addEventListener('pageswap', () => {
		document.querySelectorAll('[data-hide-on-navigate]').forEach((el) => {
			el.classList.add('hidden');
		});
	});
})();

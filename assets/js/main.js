import Alpine from 'alpinejs';
import { registerMagics } from './alpinejs/magics/index';
import { navbar, search, toc } from './alpinejs/data/index';
import { navStore, initColorScheme } from './alpinejs/stores/index';
import { bridgeTurboAndAlpine } from './helpers/index';
import persist from '@alpinejs/persist';
import focus from '@alpinejs/focus';

var debug = 0 ? console.log.bind(console, '[index]') : function () {};

// Turbolinks init.
(function () {
	document.addEventListener('turbo:render', function (e) {
		// This is also called right after the body start. This is added to prevent flicker on navigation.
		initColorScheme();
	});
})();

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

	// Start the Turbo-Alpine bridge.
	bridgeTurboAndAlpine(Alpine);

	{
		let containerScrollTops = {};

		// To preserve scroll position in scrolling elements on navigation add data-turbo-preserve-scroll-container="somename" to the scrolling container.
		addEventListener('turbo:click', () => {
			document.querySelectorAll('[data-turbo-preserve-scroll-container]').forEach((el2) => {
				containerScrollTops[el2.dataset.turboPreserveScrollContainer] = el2.scrollTop;
			});
		});

		addEventListener('turbo:render', () => {
			document.querySelectorAll('[data-turbo-preserve-scroll-container]').forEach((ele) => {
				const containerScrollTop = containerScrollTops[ele.dataset.turboPreserveScrollContainer];
				if (containerScrollTop) {
					ele.scrollTop = containerScrollTop;
				} else {
					let els = ele.querySelectorAll('.scroll-active');
					if (els.length) {
						els.forEach((el) => {
							// Avoid scrolling if el is already in view.
							if (el.offsetTop >= ele.scrollTop && el.offsetTop <= ele.scrollTop + ele.clientHeight) {
								return;
							}
							ele.scrollTop = el.offsetTop - ele.offsetTop;
						});
					}
				}
			});

			containerScrollTops = {};
		});
	}
})();

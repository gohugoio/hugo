import { initColorScheme } from './alpinejs/stores/index';

(function () {
	// This allows us to initialize the color scheme before AlpineJS etc. is loaded.
	initColorScheme();
})();

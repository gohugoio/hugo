export const navbar = (Alpine) => ({
	init: function () {
		Alpine.bind(this.$root, this.root);
	},
	root: {
		['@scroll.window.debounce.10ms'](event) {
			this.$store.nav.scroll.atTop = window.scrollY < 40 ? true : false;
		},
	},
});

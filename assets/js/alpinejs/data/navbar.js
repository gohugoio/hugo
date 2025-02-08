export const navbar = (Alpine) => ({
	atTop: true,

	init: function () {
		Alpine.bind(this.$root, this.root);
	},
	root: {
		['@scroll.window.debounce.10ms'](event) {
			this.atTop = window.scrollY < 40 ? true : false;
		},
	},
});

export const navbar = (Alpine) => ({
	init: function () {
		Alpine.bind(this.$root, this.root);

		return this.$nextTick(() => {
			let contentEl = document.querySelector('.content:not(.content--ready)');
			if (contentEl) {
				contentEl.classList.add('content--ready');
				let anchorTemplate = document.getElementById('anchor-heading');
				if (anchorTemplate) {
					let els = contentEl.querySelectorAll('h2[id], h3[id], h4[id], h5[id], h6[id], dt[id]');
					for (let i = 0; i < els.length; i++) {
						let el = els[i];
						el.classList.add('group');
						let a = anchorTemplate.content.cloneNode(true).firstElementChild;
						a.href = '#' + el.id;
						el.appendChild(a);
					}
				}
			}
		});
	},
	root: {
		['@scroll.window.debounce.10ms'](event) {
			this.$store.nav.scroll.atTop = window.scrollY < 40 ? true : false;
		},
	},
});

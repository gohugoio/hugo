var debug = 0 ? console.log.bind(console, '[toc]') : function () {};

export const toc = (Alpine) => ({
	contentScrollSpy: null,
	activeHeading: '',
	justClicked: false,

	setActive(id) {
		debug('setActive', id);
		this.activeHeading = id;
		// Prevent the intersection observer from changing the active heading right away.
		this.justClicked = true;
		setTimeout(() => {
			this.justClicked = false;
		}, 200);
	},

	init() {
		this.$watch('$store.nav.scroll.atTop', (value) => {
			if (!value) return;
			this.activeHeading = '';
			this.$root.scrollTop = 0;
		});

		return this.$nextTick(() => {
			let contentEl = document.getElementById('content');
			if (contentEl) {
				const handleIntersect = (entries) => {
					if (this.justClicked) {
						return;
					}
					for (let entry of entries) {
						if (entry.isIntersecting) {
							let id = entry.target.id;
							this.activeHeading = id;
							let liEl = this.$refs[id];
							if (liEl) {
								// If liEl is not in the viewport, scroll it into view.
								let bounding = liEl.getBoundingClientRect();
								if (bounding.top < 0 || bounding.bottom > window.innerHeight) {
									this.$root.scrollTop = liEl.offsetTop - 100;
								}
							}
							debug('intersecting', id);
							break;
						}
					}
				};

				let opts = {
					rootMargin: '0px 0px -75%',
					threshold: 0.75,
				};

				this.contentScrollSpy = new IntersectionObserver(handleIntersect, opts);
				// Observe all headings.
				let headings = contentEl.querySelectorAll('h2, h3, h4, h5, h6');
				for (let heading of headings) {
					this.contentScrollSpy.observe(heading);
				}
			}
		});
	},

	destroy() {
		if (this.contentScrollSpy) {
			debug('disconnecting');
			this.contentScrollSpy.disconnect();
		}
	},
});

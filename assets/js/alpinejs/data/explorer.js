var debug = 0 ? console.log.bind(console, '[explorer]') : function () {};

// This is cureently not used, but kept in case I change my mind.
export const explorer = (Alpine) => ({
	uiState: {
		containerScrollTop: -1,
		lastActiveRef: '',
	},
	treeState: {
		// The href of the current page.
		currentNode: '',
		// The state of each node in the tree.
		nodes: {},

		// We currently only list the sections, not regular pages, in the side bar.
		// This strikes me as the right balance. The pages gets listed on the section pages.
		// This array is sorted by length, so we can find the longest prefix of the current page
		// without having to iterate over all the keys.
		nodeRefsByLength: [],
	},
	async init() {
		let keys = Reflect.ownKeys(this.$refs);
		for (let key of keys) {
			let n = {
				open: false,
				active: false,
			};
			this.treeState.nodes[key] = n;
			this.treeState.nodeRefsByLength.push(key);
		}

		this.treeState.nodeRefsByLength.sort((a, b) => b.length - a.length);

		this.setCurrentActive();
	},

	longestPrefix(ref) {
		let longestPrefix = '';
		for (let key of this.treeState.nodeRefsByLength) {
			if (ref.startsWith(key)) {
				longestPrefix = key;
				break;
			}
		}
		return longestPrefix;
	},

	setCurrentActive() {
		let ref = this.longestPrefix(window.location.pathname);
		let activeChanged = this.uiState.lastActiveRef !== ref;
		debug('setCurrentActive', this.uiState.lastActiveRef, window.location.pathname, '=>', ref, activeChanged);
		this.uiState.lastActiveRef = ref;
		if (this.uiState.containerScrollTop === -1 && activeChanged) {
			// Navigation outside of the explorer menu.
			let el = document.querySelector(`[x-ref="${ref}"]`);
			if (el) {
				this.$nextTick(() => {
					debug('scrolling to', ref);
					el.scrollIntoView({ behavior: 'smooth', block: 'center' });
				});
			}
		}
		this.treeState.currentNode = ref;
		for (let key in this.treeState.nodes) {
			let n = this.treeState.nodes[key];
			n.active = false;
			n.open = ref == key || ref.startsWith(key);
			if (n.open) {
				debug('open', key);
			}
		}

		let n = this.treeState.nodes[this.longestPrefix(ref)];
		if (n) {
			n.active = true;
		}
	},

	getScrollingContainer() {
		return document.getElementById('leftsidebar');
	},

	onLoad() {
		debug('onLoad', this.uiState.containerScrollTop);
		if (this.uiState.containerScrollTop >= 0) {
			debug('onLoad: scrolling to', this.uiState.containerScrollTop);
			this.getScrollingContainer().scrollTo(0, this.uiState.containerScrollTop);
		}
		this.uiState.containerScrollTop = -1;
	},

	onBeforeRender() {
		debug('onBeforeRender', this.uiState.containerScrollTop);
		this.setCurrentActive();
	},

	toggleNode(ref) {
		this.uiState.containerScrollTop = this.getScrollingContainer().scrollTop;
		this.uiState.lastActiveRef = '';
		debug('toggleNode', ref, this.uiState.containerScrollTop);

		let node = this.treeState.nodes[ref];
		if (!node) {
			debug('node not found', ref);
			return;
		}
		let wasOpen = node.open;
	},

	isCurrent(ref) {
		let n = this.treeState.nodes[ref];
		return n && n.active;
	},

	isOpen(ref) {
		let node = this.treeState.nodes[ref];
		if (!node) return false;
		if (node.open) {
			debug('isOpen', ref);
		}
		return node.open;
	},
});

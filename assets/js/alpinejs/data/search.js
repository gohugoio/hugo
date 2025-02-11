const designMode = false;

const groupByLvl0 = (array) => {
	if (!array) return [];
	return array.reduce((result, currentValue) => {
		(result[currentValue.hierarchy.lvl0] = result[currentValue.hierarchy.lvl0] || []).push(currentValue);
		return result;
	}, {});
};

const applyHelperFuncs = (array) => {
	if (!array) return [];
	return array.map((item) => {
		item.getHeadingHTML = function () {
			let lvl2 = this._highlightResult.hierarchy.lvl2;
			let lvl3 = this._highlightResult.hierarchy.lvl3;

			if (!lvl3) {
				if (lvl2) {
					return lvl2.value;
				}
				return '';
			}

			if (!lvl2) {
				return lvl3.value;
			}

			return `${lvl2.value} <span class="text-gray-500">&nbsp;>&nbsp;</span> ${lvl3.value}`;
		};
		return item;
	});
};

export const search = (Alpine, cfg) => ({
	query: designMode ? 'shortcodes' : '',
	open: designMode,
	result: {},

	init() {
		Alpine.bind(this.$root, this.root);

		this.checkOpen();
		return this.$nextTick(() => {
			this.$watch('query', () => {
				this.search();
			});
		});
	},
	toggleOpen: function () {
		this.open = !this.open;
		this.checkOpen();
	},
	checkOpen: function () {
		if (!this.open) {
			return;
		}
		this.search();
		this.$nextTick(() => {
			this.$refs.input.focus();
		});
	},

	search: function () {
		if (!this.query) {
			this.result = {};
			return;
		}
		var queries = {
			requests: [
				{
					indexName: cfg.index,
					params: `query=${encodeURIComponent(this.query)}`,
					attributesToHighlight: ['hierarchy', 'content'],
					attributesToRetrieve: ['hierarchy', 'url', 'content'],
				},
			],
		};

		const host = `https://${cfg.app_id}-dsn.algolia.net`;
		const url = `${host}/1/indexes/*/queries`;

		fetch(url, {
			method: 'POST',
			headers: {
				'X-Algolia-Application-Id': cfg.app_id,
				'X-Algolia-API-Key': cfg.api_key,
			},
			body: JSON.stringify(queries),
		})
			.then((response) => response.json())
			.then((data) => {
				this.result = groupByLvl0(applyHelperFuncs(data.results[0].hits));
			});
	},
	root: {
		['@click']() {
			if (!this.open) {
				this.toggleOpen();
			}
		},
		['@search-toggle.window']() {
			this.toggleOpen();
		},
		['@keydown.meta.k.window']() {
			this.toggleOpen();
		},
	},
});

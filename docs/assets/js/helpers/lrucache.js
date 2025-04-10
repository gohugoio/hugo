// A simple LRU cache implementation backed by a map.
export class LRUCache {
	constructor(maxSize) {
		this.maxSize = maxSize;
		this.cache = new Map();
	}

	get(key) {
		return this.cache.get(key);
	}

	put(key, value) {
		if (this.cache.size >= this.maxSize) {
			const firstKey = this.cache.keys().next().value;
			this.cache.delete(firstKey);
		}
		this.cache.set(key, value);
	}
}

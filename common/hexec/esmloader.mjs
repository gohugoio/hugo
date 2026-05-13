// Node.js ESM resolver hook installed by Hugo.
//
// Node's ESM resolver does not consult NODE_PATH, unlike CJS require().
// That breaks postcss.config.js / babel.config.js / etc. files written in
// ESM and loaded from outside the project tree (typically the Hugo module
// cache): bare imports like `import x from "postcss-import"` cannot be
// resolved by walking up from the file's location.
//
// This hook makes the ESM resolver fall back to NODE_PATH for bare
// specifiers when Node's normal resolution fails. It is a no-op for
// relative/absolute paths and URL-scheme specifiers, and it never fires
// unless Node would itself have thrown ERR_MODULE_NOT_FOUND.
//
// Uses the synchronous registerHooks API so it runs on the main thread and
// does not require --allow-worker under the Node permission model.

import { registerHooks, createRequire } from 'node:module';
import { pathToFileURL } from 'node:url';

const resolvers = [];
const np = process.env.NODE_PATH;
if (np) {
	const sep = process.platform === 'win32' ? ';' : ':';
	for (const p of np.split(sep)) {
		if (p) resolvers.push(createRequire(p + '/_'));
	}
}

function isBareSpecifier(s) {
	if (!s) return false;
	if (s.startsWith('.') || s.startsWith('/') || s.startsWith('#')) return false;
	if (/^[a-z][a-z0-9+.-]*:/i.test(s)) return false;
	return true;
}

registerHooks({
	resolve(specifier, context, nextResolve) {
		try {
			return nextResolve(specifier, context);
		} catch (err) {
			if (err?.code !== 'ERR_MODULE_NOT_FOUND') throw err;
			if (!isBareSpecifier(specifier)) throw err;
			for (const r of resolvers) {
				try {
					const resolved = r.resolve(specifier);
					return { url: pathToFileURL(resolved).href, shortCircuit: true, format: null };
				} catch (_) { /* try next */ }
			}
			throw err;
		}
	},
});

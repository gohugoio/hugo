/*
Hugo adds a specific prefix, "__hugo_navigate", to the path in certain situations to signal
navigation to another content page.
*/
function HugoReload() {}

HugoReload.identifier = 'hugoReloader';
HugoReload.version = '0.9';

HugoReload.prototype.reload = function (path, options) {
	var prefix = '__hugo_navigate';

	if (path.lastIndexOf(prefix, 0) !== 0) {
		return false;
	}

	path = path.substring(prefix.length);

	var portChanged = options.overrideURL && options.overrideURL != window.location.port;

	if (!portChanged && window.location.pathname === path) {
		window.location.reload();
	} else {
		if (portChanged) {
			window.location = location.protocol + '//' + location.hostname + ':' + options.overrideURL + path;
		} else {
			window.location.pathname = path;
		}
	}

	return true;
};

LiveReload.addPlugin(HugoReload);

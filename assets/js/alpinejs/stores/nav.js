var debug = 1 ? console.log.bind(console, '[navStore]') : function () {};

var ColorScheme = {
	System: 1,
	Light: 2,
	Dark: 3,
};

const localStorageUserSettingsKey = 'hugoDocsUserSettings';

export const navStore = (Alpine) => ({
	init() {
		// There is no $watch available in Alpine stores,
		// but this has the same effect.
		this.userSettings.onColorSchemeChanged = Alpine.effect(() => {
			if (this.userSettings.settings.colorScheme) {
				this.userSettings.isDark = isDark(this.userSettings.settings.colorScheme);
				toggleDarkMode(this.userSettings.isDark);
			}
		});

		// Also react to changes in system settings.
		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
			this.userSettings.setColorScheme(ColorScheme.System);
		});
	},

	destroy() {},

	userSettings: {
		// settings gets persisted between page navigations.
		settings: Alpine.$persist({
			// light, dark or system mode.
			// If not set, we use the OS setting.
			colorScheme: ColorScheme.System,
			// Used to show the most relevant tab in config listings etc.
			configFileType: 'toml',
		}).as(localStorageUserSettingsKey),

		isDark: false,

		setColorScheme(colorScheme) {
			this.settings.colorScheme = colorScheme;
			this.isDark = isDark(colorScheme);
		},

		toggleColorScheme() {
			let next = this.settings.colorScheme + 1;
			if (next > ColorScheme.Dark) {
				next = ColorScheme.System;
			}
			this.setColorScheme(next);
		},
		colorScheme() {
			return this.settings.colorScheme ? this.settings.colorScheme : ColorScheme.System;
		},
	},
});

function isMediaDark() {
	return window.matchMedia('(prefers-color-scheme: dark)').matches;
}

function isDark(colorScheme) {
	if (!colorScheme || colorScheme == ColorScheme.System) {
		return isMediaDark();
	}

	return colorScheme == ColorScheme.Dark;
}

export function initColorScheme() {
	// The AlpineJS store has not have been initialized yet, so access the
	// localStorage directly.
	let settingsJSON = localStorage[localStorageUserSettingsKey];
	if (settingsJSON) {
		let settings = JSON.parse(settingsJSON);
		toggleDarkMode(isDark(settings.colorScheme));
		return;
	}
	toggleDarkMode(isDark(null));
}

const toggleDarkMode = function (dark) {
	if (dark) {
		document.body.classList.add('dark');
	} else {
		document.body.classList.remove('dark');
	}
};

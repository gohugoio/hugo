export const scrollToActive = (when) => {
	let els = document.querySelectorAll('.scroll-active');
	if (!els.length) {
		return;
	}
	els.forEach((el) => {
		// Find scrolling container.
		let container = el.closest('[data-turbo-preserve-scroll-container]');
		if (container) {
			// Avoid scrolling if el is already in view.
			if (el.offsetTop >= container.scrollTop && el.offsetTop <= container.scrollTop + container.clientHeight) {
				return;
			}
			container.scrollTop = el.offsetTop - container.offsetTop;
		}
	});
};

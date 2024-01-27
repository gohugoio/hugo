/**
 * Scripts which manages Code Toggle tabs.
 */
var i;
// store tabs variable
var allTabs = document.querySelectorAll("[data-toggle-tab]");
var allPanes = document.querySelectorAll("[data-pane]");

function toggleTabs(event) {

	if(event.target){
		event.preventDefault();
		var clickedTab = event.currentTarget;
		var targetKey = clickedTab.getAttribute("data-toggle-tab")
	}else {
		var targetKey = event
	}
	// We store the config language selected in users' localStorage
	if(window.localStorage){
		window.localStorage.setItem("configLangPref", targetKey)
	}
	var selectedTabs = document.querySelectorAll("[data-toggle-tab='" + targetKey + "']");
	var selectedPanes = document.querySelectorAll("[data-pane='" + targetKey + "']");

	for (var i = 0; i < allTabs.length; i++) {
		allTabs[i].classList.remove("active");
		allPanes[i].classList.remove("active");
	}

	for (var i = 0; i < selectedTabs.length; i++) {
		selectedTabs[i].classList.add("active");
		selectedPanes[i].classList.add("active");
	}

}

for (i = 0; i < allTabs.length; i++) {
	allTabs[i].addEventListener("click", toggleTabs)
}
// Upon page load, if user has a preferred language in its localStorage, tabs are set to it.
if(window.localStorage.getItem('configLangPref')) {
	toggleTabs(window.localStorage.getItem('configLangPref'))
}

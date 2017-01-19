//If this script is the last tag in the body, then the DOM has been loaded --
//i.e., the document is ready -- per:
//https://www.sitepoint.com/jquery-document-ready-plain-javascript/

//initialize
(function () {
	'use strict';

/* ** We should use this, later:

	//accept multiple window onload handlers
	//See: http://www.htmlgoodies.com/beyond/javascript/article.php/3724571/Using-Multiple-JavaScript-Onload-Functions.html

	function loadHandlerAdd(functionAdd) {
		var func;
		//assign any pre-defined function on 'window.onload' to a variable
		var functionPrevious = window.onload;
		//if no function already is hooked to it
		if (typeof functionPrevious != 'function') {
			//you can hook your function to it directly
			func = functionAdd;
		} else {
			//hook this new function instead
			func = function() {
				//call the pre-defined function
				functionPrevious();
				//call your function
				functionAdd();
			}
		}
		window.onload = func;
	}
** */

	//transform an element (by adding or removing some class)
	function elementTransformClass(elementId, action, thing) {
		var element = document.getElementById(elementId);
		if ('add' == action) {
			element.className += ' ' + thing;
		} else if ('remove' == action) {
			element.className = element.className.replace(thing, ' ');
		}
		element.className = element.className.trim();
	}

	var toggleAlternateClass = 'toggle-alternate';
	var toggleAlternateClassRegexp = new RegExp('(?:^|\\s)' + toggleAlternateClass + '(?:$|\\s)');

	//toggle an element
	function elementToggle(elementId) {
		var element = document.getElementById(elementId);
		if((element === null) || (element === undefined)) {
			return;
		}
		if (toggleAlternateClassRegexp.test(element.className)) {
			elementTransformClass(elementId, 'remove', toggleAlternateClassRegexp);
		} else {
			elementTransformClass(elementId, 'add',    toggleAlternateClass);
		}
	}

	//initialize top hamburger;
	//keep before its onclick initialization
	var hamburgerTop = document.getElementById('hamburger-top');
	hamburgerTop.setAttribute('title', 'toggle navigation sidebar');
	hamburgerTop.setAttribute('href', 'javascript:');

	//top hamburger toggles navigation menu (and other things)
	(function() {
		var manyToggle = function() {
//			elementToggle('junk-test');
			elementToggle('hamburger-top');
			elementToggle('main');
			elementToggle('nav-binary');
			elementToggle('nav-menu');
			elementToggle('TableOfContents');
		}

		//toggle on click (or keypress)
		hamburgerTop.onclick = manyToggle;
		hamburgerTop.keydown = manyToggle;
	})();

	//initialize other
	(function() {
		//initialize menu hamburger
		elementTransformClass('anchor-menu', 'add', toggleAlternateClass);

		//for some elements, allow some stylesheets to reverse the meaning of the toggle switch
		elementTransformClass('nav-menu', 'add', 'reverse');
	})();
})();

//toggle off-canvas navigation for M- screens
$('#navigation-toggle').on('click', function(evt) {
	evt.preventDefault();
	evt.stopPropagation();
  $('#site-navigation,.all-content-wrapper,#navigation-toggle,#site-footer').toggleClass('navigation-open');
});
//close navigation if body content is clicked when docs are open
$('#all-content-wrapper').on('click', function() {
  if ($('.site-navigation.navigation-open')) {
    $('.site-navigation.navigation-open,.all-content-wrapper.navigation-open,#navigation-toggle,#site-footer').removeClass('navigation-open');
  }
});

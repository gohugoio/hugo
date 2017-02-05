$('#navigation-toggle').on('click', function() {
  $('#site-navigation,.all-content-wrapper').toggleClass('navigation-open');
});

$('#all-content-wrapper').on('click',function(){
	if($('.site-navigation.navigation-open')){
		$('.site-navigation.navigation-open,.all-content-wrapper.navigation-open').removeClass('navigation-open');
	}
});

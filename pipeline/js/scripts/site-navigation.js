$('.top-menu-item-link.true').on('click', function(evt) {
  evt.preventDefault();
  evt.stopPropagation();
  var $ul = $(this).next('ul'),
    isOpen = $ul.is(':visible'),
    slideDir = isOpen ? 'slideUp' : 'slideDown',
    dur = isOpen ? 200 : 400;
  $ul.velocity(slideDir, {
    easing: 'easeOut',
    duration: dur
  });
});

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

$('.body-copy').on('click', function() {
  if ($('.toc-toggle.toc-open')) {
    document.getElementById('toc').classList.remove('toc-open');
    document.getElementById('toc-toggle').classList.remove('toc-open');
  }
});

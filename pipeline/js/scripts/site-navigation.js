$('.top-menu-item-link.has-children').on('click', function(evt) {
  evt.preventDefault();
  evt.stopPropagation();
  var $ul = $(this).next('ul'),
    siblingLinks = $('.top-menu-item-link.has-children.open').next('ul');
  siblingLinks.velocity('slideUp', { easing: 'easeOut', duration: 200 });
  var isOpen = $(this).hasClass('open'),
    slideDir = isOpen ? 'slideUp' : 'slideDown',
    dur = isOpen ? 200 : 400;
  $ul.velocity(slideDir, {
    easing: 'easeOut',
    duration: dur
  });
  siblingLinks.velocity('slideUp', { easing: 'easeOut', duration: 200 });
  if (isOpen) {
    $(this).removeClass('open');
  } else {
    $(this).addClass('open');
  }
  // console.log(siblingLinks);
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

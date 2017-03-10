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

$(document).ready(function() {
  $(".trigger").on("click", function(evt) {
    evt.preventDefault();
    evt.stopPropagation();
    $(".navbutton__wrapper").toggleClass("nav__wrapper--active");
    $('#site-navigation,.all-content-wrapper,#navigation-toggle,#site-footer').toggleClass('navigation-open');
  });
});

// //toggle off-canvas navigation for M- screens
// $('#navigation-toggle').on('click', function(evt) {


// });
//close navigation if body content is clicked when docs are open
$('#all-content-wrapper').on('click', function() {
  if ($('.site-navigation.navigation-open')) {
    $('.site-navigation.navigation-open,.all-content-wrapper.navigation-open,#site-footer').removeClass('navigation-open');
    $('.navbutton__wrapper.nav__wrapper--active').removeClass('nav__wrapper--active');
  }
});

$('.body-copy').on('click', function() {
  if ($('.toc-toggle.toc-open')) {
    document.getElementById('toc').classList.remove('toc-open');
    document.getElementById('toc-toggle').classList.remove('toc-open');
  }
});

$('.top-menu-item-link').on('click',function(evt) {
  evt.preventDefault();
  evt.stopPropagation();
  var $ul = $(this).next('ul'),
    isOpen = $ul.is(':visible'),
    slideDir = isOpen ? 'slideUp' : 'slideDown',
    dur = isOpen ? 200 : 400;
  $ul.velocity(slideDir, {
    easing: 'easeOutQuart',
    duration: dur
  });
});

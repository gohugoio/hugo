$(document).ready(function() {
  // bind click event to all internal page anchors
  $('#toc a').on('click', function(e) {
    // prevent default action and bubbling
    e.preventDefault();
    e.stopPropagation();
    // set target to anchor's "href" attribute
    var target = $(this).attr('href');
    var hashid = target.split('#')[1];
    // scroll to each target
    $(target).velocity('scroll', {
      duration: 500,
      offset: -80,
      easing: 'ease-in-out'
    });
    location.hash = hashid;
  });
});
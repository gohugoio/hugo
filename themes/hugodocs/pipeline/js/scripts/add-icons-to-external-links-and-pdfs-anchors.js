$('.body-copy a[href$=".pdf"]').append('<i class="icon-pdf"></i>');
$('.body-copy > h2,.body-copy > h3').each(function() {
	var id = $(this).attr('id');
  $(this).append('<a class="smooth-scroll heading-link" href="#' + id + '"><i class="icon-link"></i></a>');
});


$(document).ready(function () {
    // bind click event to all internal page anchors
    $('a.heading-link').on('click', function (e) {
        // prevent default action and bubbling
        e.preventDefault();
        e.stopPropagation();
        // set target to anchor's "href" attribute
        var target = $(this).attr('href');
        var hashid = target.split('#')[1];
        // scroll to each target
        $(target).velocity('scroll', {
            duration: 500,
            offset: -50,
            easing: 'ease-in-out'
        });
        location.hash = hashid;
    });
});

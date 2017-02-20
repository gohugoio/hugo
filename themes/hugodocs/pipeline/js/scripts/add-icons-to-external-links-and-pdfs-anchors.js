$('document').ready(function() {
  $('.body-copy a[href$=".pdf"]').append('<i class="icon-pdf"></i>');
  $('.body-copy > h2:not(.section-heading),.body-copy > h3').each(function() {
    var id = $(this).attr('id');
      // var baseurl = window.location.origin,
      // path = window.location.pathname,
      // fullurl = `${baseurl}${path}#${id}`;
    $(this).prepend(`<a class="smooth-scroll heading-link" href="#${id}"><i class="icon-link"></i></a>`);
  });
});
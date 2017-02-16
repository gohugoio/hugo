$('document').ready(function() {
  $('.body-copy a[href$=".pdf"]').append('<i class="icon-pdf"></i>');
  $('.body-copy > h2:not(.section-heading),.body-copy > h3').each(function() {
    var id = $(this).attr('id'),
      baseurl = window.location.origin,
      path = window.location.pathname,
      fullurl = `${baseurl}${path}#${id}`;
    $(this).append(`<a class="smooth-scroll heading-link" title="Copy heading link to clipboard" data-clipboard-text="${fullurl}" href="#${id}"><i class="icon-link"></i></a>`);
  });
});
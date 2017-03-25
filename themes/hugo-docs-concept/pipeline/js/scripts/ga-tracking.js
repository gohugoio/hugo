$('.copy-button').on('click', function() {
  let id = $(this).parent().attr('id'),
    url = window.location.pathname;
  ga('send', 'event', 'Code', 'Copy', id);
});

$('.download-button').on('click', function() {
	let id = $(this).parent().attr('id'),
    url = window.location.pathname;
  ga('send', 'event', 'Code', 'Download', id);
});

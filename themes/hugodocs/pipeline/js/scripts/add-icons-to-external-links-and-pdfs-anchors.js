$('document').ready(function() {
  $('.body-copy a[href$=".pdf"]').append('<i class="icon-pdf"></i>');
  $('.body-copy > h2,.body-copy > h3').each(function() {
    var id = $(this).attr('id'),
      baseurl = window.location.origin,
      path = window.location.pathname,
      fullurl = `${baseurl}${path}#${id}`;
    $(this).append(`<a class="smooth-scroll heading-link" title="Copy heading link to clipboard" data-clip="${fullurl}" href="#${id}"><i class="icon-link"></i></a>`);
  });
});
// const hasCopy = document.execCommand('copy', false, null);
// console.log(`${hasCopy} is hasCopy value`);

window.setTimeout(function() {
  // let headLinks = document.querySelector('.heading-link');
  // for (var i = 0; i < headLinks.length; i++) {
  //   headLinks[i].addEventListener('click', copyHeaderLink, false);
  // }
}, 2000);

function copyHeaderLink(evt) {
	let link = evt.target.dataset.clip;
	console.log(link);
}

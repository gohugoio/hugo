$('document').ready(function() {
  $('.body-copy a[href$=".pdf"]').append('<i class="icon-pdf"></i>');
  $('.body-copy > h2:not(.section-heading),.body-copy > h3').each(function() {
    var id = $(this).attr('id'),
      loc = window.location.href.split('#')[0];
    $(this).prepend(`<a class="smooth-scroll heading-link" href="#${id}" copy-text="${loc}#${id}"><i class="icon-link"></i></a>`);
  });
  // $('.heading-link').on('click',function(evt){
  // evt.preventDefault();
  // evt.stopPropagation();
  // });
  // var headingLink = new Clipboard('.heading-link');
  // headingLink.on('success', function(e) {
  // 	e.trigger.classList.add('copied');
  //   var x = window.scrollX;
  //   var y = window.scrollY;
  //   setTimeout(function() {
  //     window.scrollTo(x, y);
  //   }, 0);
  //   e.clearSelection();
  // });
  $('.heading-link').on('click', function(evt) {
    evt.preventDefault();
    evt.stopPropagation();
    var targ = $(this),
      text = $(this).attr('copy-text'),
      // set target to anchor's "href" attribute
      target = $(this).attr('href'),
      hashid = target.split('#')[1],
      textField = document.createElement('textarea');
    // scroll to each target
    $(target).velocity('scroll', {
      duration: 500,
      offset: -80,
      easing: 'ease-in-out'
    });
    location.hash = hashid;
    textField.innerText = text;
    document.body.appendChild(textField);
    textField.select();
    document.execCommand('copy');
    $(textField).remove();
    $(this).addClass('copied');
    setTimeout(function() {
      targ.removeClass('copied');
    }, 600);
  })
});

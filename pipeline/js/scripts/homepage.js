$(document).ready(function() {
  let isHome = document.getElementById('homepage') ? true : false;
  if (!isHome) {
    return;
  } else {
    $('#homepage-nav > ul.animated').addClass('fadeInDown');
    $('.home-row.hero.animated').addClass('fadeIn');
    var
      scrolled = false,
      atTop = true,
      terminal = $('.homepage-terminal'),
      gopher = $('#gopher'),
      cape = $('.gopher-cape'),
      badge = $('.gopher-badge'),
      installrow = $('.home-row.install');

    $(window).scroll(function() {
      let scroll = $(this).scrollTop(),
        offset = installrow.offset().top - 120;
      if (scroll > 20 && atTop == true) {
        $('#homepage-nav').addClass('shadow');
        atTop = false;
      } else
      if (scroll < 20) {
        $('#homepage-nav').removeClass('shadow');
        atTop = true;
      } else if (scroll > offset && scrolled == false) {
        terminal.blast({ delimiter: "letter" }).velocity("transition.fadeIn", {
          display: null,
          customClass: "visible",
          duration: 0,
          stagger: 60,
          delay: 0,
          complete: function() {
            gopher.velocity('callout.tada', {
              duration: 800,
              complete: function() {
                badge.addClass('bounceIn');
                cape.addClass('fadeIn');
              }
            });
          }
        });
        scrolled = true;
      }
      // else if (scroll > offset3)
    });
  }
});
$('#homepage-nav > ul.animated').addClass('fadeInDown');
$('.homepage-icon.animated').addClass('fadeIn');


$(document).ready(function() {
  let offset1 = $('.hero').height() * .2,
    offset2 = offset1 * 2,
    scrolled1 = false,
    scrolled2 = false,
    atTop = true,
    terminal = $('.homepage-terminal'),
    firstFade = $('.first-fade'),
    gopher = $('#gopher'),
    cape = $('.gopher-cape'),
    badge = $('.gopher-badge');
  $(window).scroll(function() {
    let scroll = $(this).scrollTop();
    if (scroll > 20 && atTop == true) {
      $('#homepage-nav').addClass('shadow');
      atTop = false;
    }else if (scroll < 20) {
      $('#homepage-nav').removeClass('shadow');
      atTop = true;
    }
    if (scroll > offset1 && scrolled1 == false) {
      firstFade.addClass('fadeIn');
      scrolled1 = true;
    } else if (scroll > offset2 && scrolled2 == false) {
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
      scrolled2 = true;
    }
    // else if (scroll > offset3)
  });
});

// var mySequence = [
//     { e: $element1, p: { translateX: 100 }, o: { duration: 1000 } },
//     /* The call below will run at the same time as the first call. */
//     { e: $element2, p: { translateX: 200 }, o: { duration: 1000, sequenceQueue: false },
//     /* As normal, the call below will run once the second call is complete. */
//     { e: $element3, p: { translateX: 300 }, o: { duration: 1000 }
// ];
// $.Velocity.RunSequence(mySequence);



// terminalElements = document.querySelectorAll('span.terminal-line');

// function terminals_init() {
//   var elements = terminalElements;
//   for (var i = 0; i < elements.length; i++) {
//     var e = elements[i]
//     e.completeHTML = e.innerHTML
//     e.innerHTML = ""
//     e.i = 0
//   }

// }

// function terminals_start() {
//   //terminalIntervalFunction = function(){ terminals_next() }
//   // setInterval(terminalIntervalFunction, 1000/60);
//   terminals_next()
// }

// function terminals_next() {
//   var done = true
//   var elements = terminalElements;

//   for (var i = 0; i < elements.length && done == true; i++) {
//     var e = elements[i]

//     e.i++

//       if (e.innerHTML.length >= e.completeHTML.length) {
//         e.innerHTML = e.completeHTML
//       } else {
//         var s = e.completeHTML.substring(0, e.i)

//         if (s.slice(-1) == " ") {
//           if (e.innerHTML.length % 6 == 0) { s = s + "" } else { s = s + "|" }
//         } else {
//           s = s + "|"
//         }

//         e.innerHTML = s
//         e.parentNode.style.display = 'none'
//         e.parentNode.style.display = 'block'
//         done = false
//       }
//   }

//   if (!done) {
//     //window.clearInterval(terminalIntervalFunction)
//     setTimeout(terminals_next, 1600 / 15)
//   }
// }

// terminals_init()

// window.onload = function() { terminals_start() };

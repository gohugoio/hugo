$('#homepage-nav > ul.animated').addClass('fadeInDown');
$('.homepage-icon.animated').addClass('fadeIn');


$(document).ready(function() {
  let heroHeight = $('.hero').height() * .2,
    heroDouble = heroHeight * 2,
    scrolled1 = false,
    scrolled2 = false;
  $(window).scroll(function() {
    let scroll = $(this).scrollTop();
    if (scroll > heroHeight && scrolled1 == false) {
      $('.svgs-one.animated').addClass('fadeIn');
      scrolled1 = true;
    } else if (scroll > heroDouble && scrolled2 == false) {
      $(".homepage-terminal")
        // Blast the text apart by word.
        .blast({ delimiter: "letter" })
        // Fade the words into view using Velocity.js.
        .velocity("transition.fadeIn", {
          display: null,
          customClass: "visible",
          duration: 0,
          stagger: 60,
          delay: 0,
          complete: function() {
            $('.after-terminal').velocity('transition.fadeIn');
          }
        });
      scrolled2 = true;
    }
  });
});
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

$('#homepage-nav > ul.animated').addClass('fadeInDown');
$('.homepage-icon.animated').addClass('fadeIn');


let terminalElements = document.querySelectorAll('span.terminal-line');

function terminals_init() {
  var elements = terminalElements;
  for (var i = 0; i < elements.length; i++) {
    var e = elements[i]
    e.completeHTML = e.innerHTML
    e.innerHTML = ""
    e.i = 0
  }

}

function terminals_start() {
  //terminalIntervalFunction = function(){ terminals_next() }
  //setInterval(terminalIntervalFunction, 1000/60);
  terminals_next()
}

function terminals_next() {
  var done = true
  var elements = terminalElements;

  for (var i = 0; i < elements.length && done == true; i++) {
    var e = elements[i]

    e.i++

      if (e.innerHTML.length >= e.completeHTML.length) {
        e.innerHTML = e.completeHTML
      } else {
        var s = e.completeHTML.substring(0, e.i)

        if (s.slice(-1) == " ") {
          if (e.innerHTML.length % 6 == 0) { s = s + "" } else { s = s + "|" }
        } else {
          s = s + "|"
        }

        e.innerHTML = s
        e.parentNode.style.display = 'none'
        e.parentNode.style.display = 'block'
        done = false
      }
  }

  if (!done) {
    //window.clearInterval(terminalIntervalFunction)
    setTimeout(terminals_next, 800 / 20)
  }
}

terminals_init()

window.onload = function() { terminals_start() };

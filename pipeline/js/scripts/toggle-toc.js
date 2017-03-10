// //removes toc if not enough headings
// (function() {
//   let tocLinks = document.querySelectorAll('#TableOfContents > ul a'),
//     toc = document.getElementById('toc');
//   if (toc && (tocLinks.length < 2)) {
//     toc.remove();
//   } else if (tocLinks.length > 1) {
//     document.getElementById('toc-toggle').addEventListener('click', toggleToc, false);
//   }
// })();

// function toggleToc(evt) {
//   evt.preventDefault();
//   evt.stopPropagation();
//   document.getElementById('toc').classList.toggle('toc-open');
//   document.getElementById('toc-toggle').classList.toggle('toc-open');
// }

var kebab = document.querySelector('.kebab'),
  middle = document.querySelector('.middle'),
  cross = document.querySelector('.cross');

if (kebab) {
  kebab.addEventListener('click', function() {
    middle.classList.toggle('active');
    cross.classList.toggle('active');
    document.getElementById('toc').classList.toggle('toc-open');
  })
}

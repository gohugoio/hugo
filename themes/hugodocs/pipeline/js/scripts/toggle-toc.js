//removes toc if not enough headings
(function() {
  let tocLinks = document.querySelectorAll('#TableOfContents > ul a');
  let toc = document.getElementById('toc');
  if (toc && (tocLinks.length < 2)) {
    toc.remove();
  } else if (tocLinks.length > 1) {
    document.getElementById('toc-toggle').addEventListener('click', toggleToc, false);
  }
})();

function toggleToc(evt) {
  evt.preventDefault();
  evt.stopPropagation();
  console.log("hello");
  document.getElementById('toc').classList.toggle('toc-open');
  document.getElementById('toc-toggle').classList.toggle('toc-open');
}

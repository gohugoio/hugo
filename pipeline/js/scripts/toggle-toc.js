var kebab = document.querySelector('.kebab'),
  middle = document.querySelector('.middle'),
  cross = document.querySelector('.cross'),
  body = document.querySelector('.body-copy');

if (kebab) {
  kebab.addEventListener('click', function() {
    middle.classList.toggle('active');
    cross.classList.toggle('active');
    document.getElementById('toc').classList.toggle('toc-open');
  });
  body.addEventListener('click', function() {
    if (document.querySelector('#toc.toc-open')) {
      middle.classList.toggle('active');
      cross.classList.toggle('active');
      document.getElementById('toc').classList.toggle('toc-open');
    }
  });
}

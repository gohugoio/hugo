$('#toggle-search').on('click', function(evt) {
  evt.preventDefault();
  evt.stopPropagation();
  $('#site-search-form').toggleClass('search-open');
  window.setTimeout(function() {
  	var sInput = document.getElementById('search-input');
    sInput.focus();
  },800);
});

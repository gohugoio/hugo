$('#toggle-search').on('click',function(evt){
	evt.preventDefault();
	evt.stopPropagation();
	$('#site-search-form').toggleClass('search-open');
});
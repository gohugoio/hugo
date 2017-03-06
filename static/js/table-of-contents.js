//IIFE to remove all third-level li's from TOC so as not to negatively affect scrollspy (these are currently hidden anyway)
(function() {
  let toc = document.getElementById('toc') ? true : false;
  if (toc) {
    let levelFours = document.querySelectorAll('#TableOfContents > ul > li > ul > li > ul > li > ul');
    if (levelFours) {
      for (var i = 0; i < levelFours.length; i++) {
        levelFours[i].remove();
      }
    }
  }
})();

$(document).ready(function() {
  // Cache selectors
  var lastId,
    tocNav = $("#TableOfContents"),
    tocNavHeight = tocNav.outerHeight() + 15,
    headerHeight = $('#site-header').height(),
    // All list items
    menuItems = tocNav.find("a"),
    // Anchors corresponding to menu items
    scrollItems = menuItems.map(function() {
      var item = $($(this).attr("href"));
      if (item.length) {
        return item;
      }
    });
  // Bind to scroll
  $(window).scroll(function(evt) {
    evt.preventDefault();
    var isBottom = $(window).scrollTop() + $(window).height() == $(document).height();
    // Get container scroll position
    var fromTop = $(this).scrollTop() + headerHeight + 100;
    // Get id of current scroll item
    var cur = scrollItems.map(function() {
      if ($(this).offset().top < fromTop)
        return this;
    });
    // Get the id of the current element
    cur = cur[cur.length - 1];
    var id = cur && cur.length ? cur[0].id : "";
    if (lastId !== id) {
      lastId = id;
      // Set/remove active class
      menuItems.parent().removeClass("active").end().filter("[href='#" + id + "']").parent().addClass("active");
      history.replaceState({}, "", menuItems.filter("[href='#" + id + "']").attr("href"));
    }
    if (isBottom) {
      menuItems.parent().removeClass('active');
      if (menuItems.last().parent().attr('class') !== 'active') {
        menuItems.last().parent().addClass('active');
      }
    }
  });
});

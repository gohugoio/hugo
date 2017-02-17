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
  console.log(tocNavHeight);

  // Bind to scroll
  $(window).scroll(function() {
    // Get container scroll position
    var fromTop = $(this).scrollTop() + headerHeight + 100;
    console.log("fromTop = " + fromTop);

    // Get id of current scroll item
    var cur = scrollItems.map(function() {
      if ($(this).offset().top < fromTop)
        return this;
    });
    // Get the id of the current element
    cur = cur[cur.length - 1];
    var id = cur && cur.length ? cur[0].id : "";
    var isBottom = $(window).scrollTop() + $(window).height() == $(document).height();

    if (lastId !== id) {
      lastId = id;
      // Set/remove active class
      menuItems
        .parent().removeClass("active")
        .end().filter("[href='#" + id + "']").parent().addClass("active");
    }
    if (isBottom) {
      menuItems.parent().removeClass('active');
      if (menuItems.last().parent().attr('class') !== 'active') {
        menuItems.last().parent().addClass('active');
      }
    }
  });
});

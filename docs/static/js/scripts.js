function initializeJS() {

    //tool tips
    jQuery('.tooltips').tooltip();

    //popovers
    jQuery('.popovers').popover();

    //search box setup
    var duration = 0; // milliseconds
    var searchBoxSel = '\
div#container-page > main#main > article#article > div.cse';
    jQuery(searchBoxSel).show(duration);

    //close all inactive submenus
    var duration = 0; // milliseconds
    var allMenus = jQuery('ul#nav-menu > li.menu-dynamic:not(.active)');
    allMenus.removeClass("open");
    var allSubmenus = jQuery('ul#nav-menu > li.menu-dynamic:not(.active) > ul.menu-sub');
    allSubmenus.removeClass("open");
    allSubmenus.hide(duration);

    //sidebar dropdown menu
    jQuery('ul#nav-menu > li.menu-dynamic > a.name').click(function () {
        //close previous open submenu
        var last = jQuery('.menu-sub.open', jQuery('ul#nav-menu'));
        jQuery(last).slideUp(200);
        jQuery(last).removeClass("open");
        jQuery(last).parent().removeClass("open");
        jQuery('.menu-arrow', jQuery(last).parent()).addClass('fa-angle-right');
        jQuery('.menu-arrow', jQuery(last).parent()).removeClass('fa-angle-down');

        //toggle current submenu
        var sub = jQuery(this).next();
        if (sub.is(":visible")) {
            jQuery('.menu-arrow', this).addClass('fa-angle-right');
            jQuery('.menu-arrow', this).removeClass('fa-angle-down');
            sub.slideUp(200);
            jQuery(sub).removeClass("open");
            jQuery(this.parentNode).removeClass("open");
        } else {
            jQuery('.menu-arrow', this).addClass('fa-angle-down');
            jQuery('.menu-arrow', this).removeClass('fa-angle-right');
            sub.slideDown(200);
            jQuery(sub).addClass("open");
            jQuery(this.parentNode).addClass("open");
        }

        //center menu on screen
        var o = (jQuery(this).offset());
        diff = 200 - o.top;
        if(diff>0)
            jQuery('ul#nav-menu').scrollTo("-="+Math.abs(diff),500);

        else
            jQuery('ul#nav-menu').scrollTo("+="+Math.abs(diff),500);
    });

    //bar chart
    if (jQuery(".custom-custom-bar-chart")) {
        jQuery(".bar").each(function () {
            var i = jQuery(this).find(".value").html();
            jQuery(this).find(".value").html("");
            jQuery(this).find(".value").animate({
                height: i
            }, 2000)
        })
    }

    // GitHub buttons
    (function(jQuery) {
        (function() {
                var caches = {};
                $.fn.showGithub = function(user, repo, type, isCount) {
                        $(this).each(function() {
                                var $e = $(this);
                                var user = $e.data('user') || user;
                                var repo = $e.data('repo') || repo;
                                var type = $e.data('type') || type;
                                var isCount = ($e.data('count') == 'true') || isCount || true;
                                var $mainButton = $e.html(
'<span class="repo-btn">\
    <a class="btn-default btn btn-xs"\
    rel="noopener noreferrer" href="#">\
        <i class="icon-github"></i>\
        <span class="gh-text"></span>\
    </a>\
    <a class="gh-count"\
    rel="noopener noreferrer" href="#">\
    </a>\
</span>'
                                ).find('span.repo-btn');
                                var $button = $mainButton.find('a.btn');
                                var $text = $mainButton.find('a.btn > span.gh-text');
                                var $counter = $mainButton.find('a.gh-count');

                                function addCommas(a) {
                                        return String(a).replace(/(\d)(?=(\d{3})+$)/g, '$1,');
                                }

                                function callback(a) {
                                        if (type == 'follow') {
                                                $counter.html(addCommas(a.followers));
                                        } else if (type == 'star') {
                                                $counter.html(addCommas(a.watchers));
                                        } else if (type == 'fork') {
                                                $counter.html(addCommas(a.forks));
                                        }

                                        if (isCount) {
                                                $counter.css('display', 'inline-block');
                                        }
                                }

                                function jsonp(urlParam) {
                                        var ctx = caches[urlParam] || {};
                                        caches[urlParam] = ctx;
                                        if (ctx.onload || ctx.data){
                                                if (ctx.data){
                                                        callback(ctx.data);
                                                } else {
                                                        var duration = 5000; // milliseconds
                                                        setTimeout(jsonp, duration, urlParam);
                                                }
                                        }else{
                                                ctx.onload = true;
                                                $.getJSON(urlParam, function(a) {
                                                        ctx.onload = false;
                                                        ctx.data = a;
                                                        callback(a);
                                                });
                                        }
                                }

                                var urlUser = ['https://github.com', user].join('/');
                                var urlRepo = [urlUser, repo].join('/');
                                var title;

                                if (type == 'follow') {
                                        title = 'followers';
                                        $counter.attr('href', [urlUser, title].join('/'));
                                } else if (type == 'star') {
                                        title = 'stargazers';
                                        $counter.attr('href', [urlRepo, title].join('/'));
                                } else if (type == 'fork') {
                                        title = 'forks';
                                        $counter.attr('href', [urlRepo, 'network', 'members'].join('/'));
                                }

                                var urlBaseApi = 'https://api.github.com';
                                if (type == 'follow') {
                                        $text.html('@' + user);
                                        $button.attr('href', urlUser);
                                        jsonp([urlBaseApi, 'users', user].join('/'));
                                } else {
                                        $text.html(type);
                                        $button.attr('href', urlRepo);
                                        jsonp([urlBaseApi, 'repos', user, repo].join('/'));
                                }

                                $mainButton.addClass(['github',title].join('-'));
                                $button. attr('title', type);
                                $counter.attr('title', title);

                        });
                };

        })();

        $('[rel=show-repo]').showGithub();
    })();
}

jQuery(document).ready(function(){
    initializeJS();
});

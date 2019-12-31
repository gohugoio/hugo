/*!
 *  Sharrre.com - Make your sharing widget!
 *  Version: beta 1.3.5
 *  Author: Julien Hany
 *  License: MIT http://en.wikipedia.org/wiki/MIT_License or GPLv2 http://en.wikipedia.org/wiki/GNU_General_Public_License
 */

;
(function($, window, document, undefined) {

    /* Defaults
    ================================================== */
    var pluginName = 'sharrre',
        defaults = {
            className: 'sharrre',
            share: {
                googlePlus: false,
                facebook: false,
                twitter: false,
                digg: false,
                delicious: false,
                stumbleupon: false,
                linkedin: false,
                pinterest: false
            },
            shareTotal: 0,
            template: '',
            title: '',
            url: document.location.href,
            text: document.title,
            urlCurl: 'sharrre.php', //PHP script for google plus...
            count: {}, //counter by social network
            total: 0, //total of sharing
            shorterTotal: true, //show total by k or M when number is to big
            enableHover: true, //disable if you want to personalize hover event with callback
            enableCounter: true, //disable if you just want use buttons
            enableTracking: false, //tracking with google analitycs
            hover: function() {}, //personalize hover event with this callback function
            hide: function() {}, //personalize hide event with this callback function
            click: function() {}, //personalize click event with this callback function
            render: function() {}, //personalize render event with this callback function
            buttons: { //settings for buttons
                googlePlus: { //http://www.google.com/webmasters/+1/button/
                    url: '', //if you need to personnalize button url
                    urlCount: false, //if you want to use personnalize button url on global counter
                    size: 'medium',
                    lang: 'en-US',
                    annotation: ''
                },
                facebook: { //http://developers.facebook.com/docs/reference/plugins/like/
                    url: '', //if you need to personalize url button
                    urlCount: false, //if you want to use personnalize button url on global counter
                    action: 'like',
                    layout: 'button_count',
                    width: '',
                    send: 'false',
                    faces: 'false',
                    colorscheme: '',
                    font: '',
                    lang: 'en_US'
                },
                twitter: { //http://twitter.com/about/resources/tweetbutton
                    url: '', //if you need to personalize url button
                    urlCount: false, //if you want to use personnalize button url on global counter
                    count: 'horizontal',
                    hashtags: '',
                    via: '',
                    related: '',
                    lang: 'en'
                },
                digg: { //http://about.digg.com/downloads/button/smart
                    url: '', //if you need to personalize url button
                    urlCount: false, //if you want to use personnalize button url on global counter
                    type: 'DiggCompact'
                },
                delicious: {
                    url: '', //if you need to personalize url button
                    urlCount: false, //if you want to use personnalize button url on global counter
                    size: 'medium' //medium or tall
                },
                stumbleupon: { //http://www.stumbleupon.com/badges/
                    url: '', //if you need to personalize url button
                    urlCount: false, //if you want to use personnalize button url on global counter
                    layout: '1'
                },
                linkedin: { //http://developer.linkedin.com/plugins/share-button
                    url: '', //if you need to personalize url button
                    urlCount: false, //if you want to use personnalize button url on global counter
                    counter: ''
                },
                pinterest: { //http://pinterest.com/about/goodies/
                    url: '', //if you need to personalize url button
                    media: '',
                    description: '',
                    layout: 'horizontal'
                }
            }
        },
        /* Json URL to get count number
        ================================================== */
        urlJson = {
            googlePlus: "",

            //new FQL method by Sire
            facebook: "https://graph.facebook.com/fql?q=SELECT%20url,%20normalized_url,%20share_count,%20like_count,%20comment_count,%20total_count,commentsbox_count,%20comments_fbid,%20click_count%20FROM%20link_stat%20WHERE%20url=%27{url}%27&callback=?",
            //old method facebook: "http://graph.facebook.com/?id={url}&callback=?",
            //facebook : "http://api.ak.facebook.com/restserver.php?v=1.0&method=links.getStats&urls={url}&format=json"

            twitter: "http://cdn.api.twitter.com/1/urls/count.json?url={url}&callback=?",
            digg: "http://services.digg.com/2.0/story.getInfo?links={url}&type=javascript&callback=?",
            delicious: 'http://feeds.delicious.com/v2/json/urlinfo/data?url={url}&callback=?',
            //stumbleupon: "http://www.stumbleupon.com/services/1.01/badge.getinfo?url={url}&format=jsonp&callback=?",
            stumbleupon: "",
            linkedin: "http://www.linkedin.com/countserv/count/share?format=jsonp&url={url}&callback=?",
            pinterest: "http://api.pinterest.com/v1/urls/count.json?url={url}&callback=?"
        },
        /* Load share buttons asynchronously
        ================================================== */
        loadButton = {
            googlePlus: function(self) {
                var sett = self.options.buttons.googlePlus;
                //$(self.element).find('.buttons').append('<div class="button googleplus"><g:plusone size="'+self.options.buttons.googlePlus.size+'" href="'+self.options.url+'"></g:plusone></div>');
                $(self.element).find('.buttons').append('<div class="button googleplus"><div class="g-plusone" data-size="' + sett.size + '" data-href="' + (sett.url !== '' ? sett.url : self.options.url) + '" data-annotation="' + sett.annotation + '"></div></div>');
                window.___gcfg = {
                    lang: self.options.buttons.googlePlus.lang
                };
                var loading = 0;
                if (typeof gapi === 'undefined' && loading == 0) {
                    loading = 1;
                    (function() {
                        var po = document.createElement('script');
                        po.type = 'text/javascript';
                        po.async = true;
                        po.src = '//apis.google.com/js/plusone.js';
                        var s = document.getElementsByTagName('script')[0];
                        s.parentNode.insertBefore(po, s);
                    })();
                } else {
                    gapi.plusone.go();
                }
            },
            facebook: function(self) {
                var sett = self.options.buttons.facebook;
                $(self.element).find('.buttons').append('<div class="button facebook"><div id="fb-root"></div><div class="fb-like" data-href="' + (sett.url !== '' ? sett.url : self.options.url) + '" data-send="' + sett.send + '" data-layout="' + sett.layout + '" data-width="' + sett.width + '" data-show-faces="' + sett.faces + '" data-action="' + sett.action + '" data-colorscheme="' + sett.colorscheme + '" data-font="' + sett.font + '" data-via="' + sett.via + '"></div></div>');
                var loading = 0;
                if (typeof FB === 'undefined' && loading == 0) {
                    loading = 1;
                    (function(d, s, id) {
                        var js, fjs = d.getElementsByTagName(s)[0];
                        if (d.getElementById(id)) {
                            return;
                        }
                        js = d.createElement(s);
                        js.id = id;
                        js.src = '//connect.facebook.net/' + sett.lang + '/all.js#xfbml=1';
                        fjs.parentNode.insertBefore(js, fjs);
                    }(document, 'script', 'facebook-jssdk'));
                } else {
                    FB.XFBML.parse();
                }
            },
            twitter: function(self) {
                var sett = self.options.buttons.twitter;
                $(self.element).find('.buttons').append('<div class="button twitter"><a href="https://twitter.com/share" class="twitter-share-button" data-url="' + (sett.url !== '' ? sett.url : self.options.url) + '" data-count="' + sett.count + '" data-text="' + self.options.text + '" data-via="' + sett.via + '" data-hashtags="' + sett.hashtags + '" data-related="' + sett.related + '" data-lang="' + sett.lang + '">Tweet</a></div>');
                var loading = 0;
                if (typeof twttr === 'undefined' && loading == 0) {
                    loading = 1;
                    (function() {
                        var twitterScriptTag = document.createElement('script');
                        twitterScriptTag.type = 'text/javascript';
                        twitterScriptTag.async = true;
                        twitterScriptTag.src = '//platform.twitter.com/widgets.js';
                        var s = document.getElementsByTagName('script')[0];
                        s.parentNode.insertBefore(twitterScriptTag, s);
                    })();
                } else {
                    $.ajax({
                        url: '//platform.twitter.com/widgets.js',
                        dataType: 'script',
                        cache: true
                    }); //http://stackoverflow.com/q/6536108
                }
            },
            digg: function(self) {
                var sett = self.options.buttons.digg;
                $(self.element).find('.buttons').append('<div class="button digg"><a class="DiggThisButton ' + sett.type + '" rel="nofollow external" href="http://digg.com/submit?url=' + encodeURIComponent((sett.url !== '' ? sett.url : self.options.url)) + '"></a></div>');
                var loading = 0;
                if (typeof __DBW === 'undefined' && loading == 0) {
                    loading = 1;
                    (function() {
                        var s = document.createElement('SCRIPT'),
                            s1 = document.getElementsByTagName('SCRIPT')[0];
                        s.type = 'text/javascript';
                        s.async = true;
                        s.src = '//widgets.digg.com/buttons.js';
                        s1.parentNode.insertBefore(s, s1);
                    })();
                }
            },
            delicious: function(self) {
                if (self.options.buttons.delicious.size == 'tall') { //tall
                    var css = 'width:50px;',
                        cssCount = 'height:35px;width:50px;font-size:15px;line-height:35px;',
                        cssShare = 'height:18px;line-height:18px;margin-top:3px;';
                } else { //medium
                    var css = 'width:93px;',
                        cssCount = 'float:right;padding:0 3px;height:20px;width:26px;line-height:20px;',
                        cssShare = 'float:left;height:20px;line-height:20px;';
                }
                var count = self.shorterTotal(self.options.count.delicious);
                if (typeof count === "undefined") {
                    count = 0;
                }
                $(self.element).find('.buttons').append(
                    '<div class="button delicious"><div style="' + css + 'font:12px Arial,Helvetica,sans-serif;cursor:pointer;color:#666666;display:inline-block;float:none;height:20px;line-height:normal;margin:0;padding:0;text-indent:0;vertical-align:baseline;">' +
                    '<div style="' + cssCount + 'background-color:#fff;margin-bottom:5px;overflow:hidden;text-align:center;border:1px solid #ccc;border-radius:3px;">' + count + '</div>' +
                    '<div style="' + cssShare + 'display:block;padding:0;text-align:center;text-decoration:none;width:50px;background-color:#7EACEE;border:1px solid #40679C;border-radius:3px;color:#fff;">' +
                    '<img src="http://www.delicious.com/static/img/delicious.small.gif" height="10" width="10" alt="Delicious" /> Add</div></div></div>');

                $(self.element).find('.delicious').on('click', function() {
                    self.openPopup('delicious');
                });
            },
            stumbleupon: function(self) {
                var sett = self.options.buttons.stumbleupon;
                $(self.element).find('.buttons').append('<div class="button stumbleupon"><su:badge layout="' + sett.layout + '" location="' + (sett.url !== '' ? sett.url : self.options.url) + '"></su:badge></div>');
                var loading = 0;
                if (typeof STMBLPN === 'undefined' && loading == 0) {
                    loading = 1;
                    (function() {
                        var li = document.createElement('script');
                        li.type = 'text/javascript';
                        li.async = true;
                        li.src = '//platform.stumbleupon.com/1/widgets.js';
                        var s = document.getElementsByTagName('script')[0];
                        s.parentNode.insertBefore(li, s);
                    })();
                    s = window.setTimeout(function() {
                        if (typeof STMBLPN !== 'undefined') {
                            STMBLPN.processWidgets();
                            clearInterval(s);
                        }
                    }, 500);
                } else {
                    STMBLPN.processWidgets();
                }
            },
            linkedin: function(self) {
                var sett = self.options.buttons.linkedin;
                $(self.element).find('.buttons').append('<div class="button linkedin"><script type="in/share" data-url="' + (sett.url !== '' ? sett.url : self.options.url) + '" data-counter="' + sett.counter + '"></script></div>');
                var loading = 0;
                if (typeof window.IN === 'undefined' && loading == 0) {
                    loading = 1;
                    (function() {
                        var li = document.createElement('script');
                        li.type = 'text/javascript';
                        li.async = true;
                        li.src = '//platform.linkedin.com/in.js';
                        var s = document.getElementsByTagName('script')[0];
                        s.parentNode.insertBefore(li, s);
                    })();
                } else {
                    window.IN.init();
                }
            },
            pinterest: function(self) {
                var sett = self.options.buttons.pinterest;
                $(self.element).find('.buttons').append('<div class="button pinterest"><a href="http://pinterest.com/pin/create/button/?url=' + (sett.url !== '' ? sett.url : self.options.url) + '&media=' + sett.media + '&description=' + sett.description + '" class="pin-it-button" count-layout="' + sett.layout + '">Pin It</a></div>');

                (function() {
                    var li = document.createElement('script');
                    li.type = 'text/javascript';
                    li.async = true;
                    li.src = '//assets.pinterest.com/js/pinit.js';
                    var s = document.getElementsByTagName('script')[0];
                    s.parentNode.insertBefore(li, s);
                })();
            }
        },
        /* Tracking for Google Analytics
        ================================================== */
        tracking = {
            googlePlus: function() {},
            facebook: function() {
                //console.log('facebook');
                fb = window.setInterval(function() {
                    if (typeof FB !== 'undefined') {
                        FB.Event.subscribe('edge.create', function(targetUrl) {
                            _gaq.push(['_trackSocial', 'facebook', 'like', targetUrl]);
                        });
                        FB.Event.subscribe('edge.remove', function(targetUrl) {
                            _gaq.push(['_trackSocial', 'facebook', 'unlike', targetUrl]);
                        });
                        FB.Event.subscribe('message.send', function(targetUrl) {
                            _gaq.push(['_trackSocial', 'facebook', 'send', targetUrl]);
                        });
                        //console.log('ok');
                        clearInterval(fb);
                    }
                }, 1000);
            },
            twitter: function() {
                //console.log('twitter');
                tw = window.setInterval(function() {
                    if (typeof twttr !== 'undefined') {
                        twttr.events.bind('tweet', function(event) {
                            if (event) {
                                _gaq.push(['_trackSocial', 'twitter', 'tweet']);
                            }
                        });
                        //console.log('ok');
                        clearInterval(tw);
                    }
                }, 1000);
            },
            digg: function() {
                //if somenone find a solution, mail me !
                /*$(this.element).find('.digg').on('click', function(){
                  _gaq.push(['_trackSocial', 'digg', 'add']);
                });*/
            },
            delicious: function() {},
            stumbleupon: function() {},
            linkedin: function() {
                function LinkedInShare() {
                    _gaq.push(['_trackSocial', 'linkedin', 'share']);
                }
            },
            pinterest: function() {
                //if somenone find a solution, mail me !
            }
        },
        /* Popup for each social network
        ================================================== */
        popup = {
            googlePlus: function(opt) {
                window.open("https://plus.google.com/share?hl=" + opt.buttons.googlePlus.lang + "&url=" + encodeURIComponent((opt.buttons.googlePlus.url !== '' ? opt.buttons.googlePlus.url : opt.url)), "", "toolbar=0, status=0, width=900, height=500");
            },
            facebook: function(opt) {
                window.open("http://www.facebook.com/sharer/sharer.php?u=" + encodeURIComponent((opt.buttons.facebook.url !== '' ? opt.buttons.facebook.url : opt.url)) + "&t=" + opt.text + "", "", "toolbar=0, status=0, width=900, height=500");
            },
            twitter: function(opt) {
                window.open("https://twitter.com/intent/tweet?text=" + encodeURIComponent(opt.text) + "&url=" + encodeURIComponent((opt.buttons.twitter.url !== '' ? opt.buttons.twitter.url : opt.url)) + (opt.buttons.twitter.via !== '' ? '&via=' + opt.buttons.twitter.via : ''), "", "toolbar=0, status=0, width=650, height=360");
            },
            digg: function(opt) {
                window.open("http://digg.com/tools/diggthis/submit?url=" + encodeURIComponent((opt.buttons.digg.url !== '' ? opt.buttons.digg.url : opt.url)) + "&title=" + opt.text + "&related=true&style=true", "", "toolbar=0, status=0, width=650, height=360");
            },
            delicious: function(opt) {
                window.open('http://www.delicious.com/save?v=5&noui&jump=close&url=' + encodeURIComponent((opt.buttons.delicious.url !== '' ? opt.buttons.delicious.url : opt.url)) + '&title=' + opt.text, 'delicious', 'toolbar=no,width=550,height=550');
            },
            stumbleupon: function(opt) {
                window.open('http://www.stumbleupon.com/badge/?url=' + encodeURIComponent((opt.buttons.stumbleupon.url !== '' ? opt.buttons.stumbleupon.url : opt.url)), 'stumbleupon', 'toolbar=no,width=550,height=550');
            },
            linkedin: function(opt) {
                window.open('https://www.linkedin.com/cws/share?url=' + encodeURIComponent((opt.buttons.linkedin.url !== '' ? opt.buttons.linkedin.url : opt.url)) + '&token=&isFramed=true', 'linkedin', 'toolbar=no,width=550,height=550');
            },
            pinterest: function(opt) {
                window.open('http://pinterest.com/pin/create/button/?url=' + encodeURIComponent((opt.buttons.pinterest.url !== '' ? opt.buttons.pinterest.url : opt.url)) + '&media=' + encodeURIComponent(opt.buttons.pinterest.media) + '&description=' + opt.buttons.pinterest.description, 'pinterest', 'toolbar=no,width=700,height=300');
            }
        };

    /* Plugin constructor
    ================================================== */
    function Plugin(element, options) {
        this.element = element;

        this.options = $.extend(true, {}, defaults, options);
        this.options.share = options.share; //simple solution to allow order of buttons

        this._defaults = defaults;
        this._name = pluginName;

        this.init();
    };

    /* Initialization method
    ================================================== */
    Plugin.prototype.init = function() {
        var self = this;
        if (this.options.urlCurl !== '') {
            urlJson.googlePlus = this.options.urlCurl + '?url={url}&type=googlePlus'; // PHP script for GooglePlus...
            urlJson.stumbleupon = this.options.urlCurl + '?url={url}&type=stumbleupon'; // PHP script for Stumbleupon...
        }
        $(this.element).addClass(this.options.className); //add class

        //HTML5 Custom data
        if (typeof $(this.element).data('title') !== 'undefined') {
            this.options.title = $(this.element).attr('data-title');
        }
        if (typeof $(this.element).data('url') !== 'undefined') {
            this.options.url = $(this.element).data('url');
        }
        if (typeof $(this.element).data('text') !== 'undefined') {
            this.options.text = $(this.element).data('text');
        }

        //how many social website have been selected
        $.each(this.options.share, function(name, val) {
            if (val === true) {
                self.options.shareTotal++;
            }
        });

        if (self.options.enableCounter === true) { //if for some reason you don't need counter
            //get count of social share that have been selected
            $.each(this.options.share, function(name, val) {
                if (val === true) {
                    //self.getSocialJson(name);
                    try {
                        self.getSocialJson(name);
                    } catch (e) {}
                }
            });
        } else if (self.options.template !== '') { //for personalized button (with template)
            this.options.render(this, this.options);
        } else { // if you want to use official button like example 3 or 5
            this.loadButtons();
        }

        //add hover event
        $(this.element).hover(function() {
            //load social button if enable and 1 time
            if ($(this).find('.buttons').length === 0 && self.options.enableHover === true) {
                self.loadButtons();
            }
            self.options.hover(self, self.options);
        }, function() {
            self.options.hide(self, self.options);
        });

        //click event
        $(this.element).click(function() {
            self.options.click(self, self.options);
            return false;
        });
    };

    /* loadButtons methode
    ================================================== */
    Plugin.prototype.loadButtons = function() {
        var self = this;
        $(this.element).append('<div class="buttons"></div>');
        $.each(self.options.share, function(name, val) {
            if (val == true) {
                loadButton[name](self);
                if (self.options.enableTracking === true) { //add tracking
                    tracking[name]();
                }
            }
        });
    };

    /* getSocialJson methode
    ================================================== */
    Plugin.prototype.getSocialJson = function(name) {
        var self = this,
            count = 0,
            url = urlJson[name].replace('{url}', encodeURIComponent(this.options.url));
        if (this.options.buttons[name].urlCount === true && this.options.buttons[name].url !== '') {
            url = urlJson[name].replace('{url}', this.options.buttons[name].url);
        }
        //console.log('name : ' + name + ' - url : '+url); //debug
        if (url != '' && self.options.urlCurl !== '') { //urlCurl = '' if you don't want to used PHP script but used social button
            $.getJSON(url, function(json) {
                    if (typeof json.count !== "undefined") { //GooglePlus, Stumbleupon, Twitter, Pinterest and Digg
                        var temp = json.count + '';
                        temp = temp.replace('\u00c2\u00a0', ''); //remove google plus special chars
                        count += parseInt(temp, 10);
                    }
                    //get the FB total count (shares, likes and more)
                    else if (json.data && json.data.length > 0 && typeof json.data[0].total_count !== "undefined") { //Facebook total count
                        count += parseInt(json.data[0].total_count, 10);
                    } else if (typeof json[0] !== "undefined") { //Delicious
                        count += parseInt(json[0].total_posts, 10);
                    } else if (typeof json[0] !== "undefined") { //Stumbleupon
                    }
                    self.options.count[name] = count;
                    self.options.total += count;
                    self.renderer();
                    self.rendererPerso();
                    //console.log(json); //debug
                })
                .error(function() {
                    self.options.count[name] = 0;
                    self.rendererPerso();
                });
        } else {
            self.renderer();
            self.options.count[name] = 0;
            self.rendererPerso();
        }
    };

    /* launch render methode
    ================================================== */
    Plugin.prototype.rendererPerso = function() {
        //check if this is the last social website to launch render
        var shareCount = 0;
        for (e in this.options.count) {
            shareCount++;
        }
        if (shareCount === this.options.shareTotal) {
            this.options.render(this, this.options);
        }
    };

    /* render methode
    ================================================== */
    Plugin.prototype.renderer = function() {
        var total = this.options.total,
            template = this.options.template;
        if (this.options.shorterTotal === true) { //format number like 1.2k or 5M
            total = this.shorterTotal(total);
        }

        if (template !== '') { //if there is a template
            template = template.replace('{total}', total);
            $(this.element).html(template);
        } else { //template by defaults
            $(this.element).html(
                '<div class="box"><a class="count" href="#">' + total + '</a>' +
                (this.options.title !== '' ? '<a class="share" href="#">' + this.options.title + '</a>' : '') +
                '</div>'
            );
        }
    };

    /* format total numbers like 1.2k or 5M
    ================================================== */
    Plugin.prototype.shorterTotal = function(num) {
        if (num >= 1e6) {
            num = (num / 1e6).toFixed(2) + "M"
        } else if (num >= 1e3) {
            num = (num / 1e3).toFixed(1) + "k"
        }
        return num;
    };

    /* Methode for open popup
    ================================================== */
    Plugin.prototype.openPopup = function(site) {
        popup[site](this.options); //open
        if (this.options.enableTracking === true) { //tracking!
            var tracking = {
                googlePlus: {
                    site: 'Google',
                    action: '+1'
                },
                facebook: {
                    site: 'facebook',
                    action: 'like'
                },
                twitter: {
                    site: 'twitter',
                    action: 'tweet'
                },
                digg: {
                    site: 'digg',
                    action: 'add'
                },
                delicious: {
                    site: 'delicious',
                    action: 'add'
                },
                stumbleupon: {
                    site: 'stumbleupon',
                    action: 'add'
                },
                linkedin: {
                    site: 'linkedin',
                    action: 'share'
                },
                pinterest: {
                    site: 'pinterest',
                    action: 'pin'
                }
            };
            _gaq.push(['_trackSocial', tracking[site].site, tracking[site].action]);
        }
    };

    /* Methode for add +1 to a counter
    ================================================== */
    Plugin.prototype.simulateClick = function() {
        var html = $(this.element).html();
        $(this.element).html(html.replace(this.options.total, this.options.total + 1));
    };

    /* Methode for add +1 to a counter
    ================================================== */
    Plugin.prototype.update = function(url, text) {
        if (url !== '') {
            this.options.url = url;
        }
        if (text !== '') {
            this.options.text = text;
        }
    };

    /* A really lightweight plugin wrapper around the constructor, preventing against multiple instantiations
    ================================================== */
    $.fn[pluginName] = function(options) {
        var args = arguments;
        if (options === undefined || typeof options === 'object') {
            return this.each(function() {
                if (!$.data(this, 'plugin_' + pluginName)) {
                    $.data(this, 'plugin_' + pluginName, new Plugin(this, options));
                }
            });
        } else if (typeof options === 'string' && options[0] !== '_' && options !== 'init') {
            return this.each(function() {
                var instance = $.data(this, 'plugin_' + pluginName);
                if (instance instanceof Plugin && typeof instance[options] === 'function') {
                    instance[options].apply(instance, Array.prototype.slice.call(args, 1));
                }
            });
        }
    };
})(jQuery, window, document);
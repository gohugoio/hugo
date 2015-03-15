---
date: 2014-03-10
linktitle: Migrating from Jekyll
menu:
  main:
    parent: tutorials
prev: /tutorials/mathjax
title: Migrate to Hugo from Jekyll
weight: 10
---

## Move static content to `static`
Jekyll has a rule that any directory not starting with `_` will be copied as-is to the `_site` output. Hugo keeps all static content under `static`. You should therefore move it all there.
With Jekyll, something that looked like

    ▾ <root>/
        ▾ images/
            logo.png

should become

    ▾ <root>/
        ▾ static/
            ▾ images/
                logo.png

Additionally, you'll want any files that should reside at the root (such as `CNAME`) to be moved to `static`.

## Create your Hugo configuration file
Hugo can read your configuration as JSON, YAML or TOML. Hugo supports parameters custom configuration too. Refer to the [Hugo configuration documentation](/overview/configuration/) for details.

## Set your configuration publish folder to `_site`
The default is for Jekyll to publish to `_site` and for Hugo to publish to `public`. If, like me, you have [`_site` mapped to a git submodule on the `gh-pages` branch](http://blog.blindgaenger.net/generate_github_pages_in_a_submodule.html), you'll want to do one of two alternatives:

1. Change your submodule to point to map `gh-pages` to public instead of `_site` (recommended).

        git submodule deinit _site
        git rm _site
        git submodule add -b gh-pages git@github.com:your-username/your-repo.git public

2. Or, change the Hugo configuration to use `_site` instead of `public`.

        {
            ..
            "publishdir": "_site",
            ..
        }

## Convert Jekyll templates to Hugo templates
That's the bulk of the work right here. The documentation is your friend. You should refer to [Jekyll's template documentation](http://jekyllrb.com/docs/templates/) if you need to refresh your memory on how you built your blog and [Hugo's template](/layout/templates/) to learn Hugo's way.

As a single reference data point, converting my templates for [heyitsalex.net](http://heyitsalex.net/) took me no more than a few hours.

## Convert Jekyll plugins to Hugo shortcodes
Jekyll has [plugins](http://jekyllrb.com/docs/plugins/); Hugo has [shortcodes](/doc/shortcodes/). It's fairly trivial to do a port.

### Implementation
As an example, I was using a custom [`image_tag`](https://github.com/alexandre-normand/alexandre-normand/blob/74bb12036a71334fdb7dba84e073382fc06908ec/_plugins/image_tag.rb) plugin to generate figures with caption when running Jekyll. As I read about shortcodes, I found Hugo had a nice built-in shortcode that does exactly the same thing.

Jekyll's plugin:

    module Jekyll
      class ImageTag < Liquid::Tag
        @url = nil
        @caption = nil
        @class = nil
        @link = nil
        // Patterns
        IMAGE_URL_WITH_CLASS_AND_CAPTION =
        IMAGE_URL_WITH_CLASS_AND_CAPTION_AND_LINK = /(\w+)(\s+)((https?:\/\/|\/)(\S+))(\s+)"(.*?)"(\s+)->((https?:\/\/|\/)(\S+))(\s*)/i
        IMAGE_URL_WITH_CAPTION = /((https?:\/\/|\/)(\S+))(\s+)"(.*?)"/i
        IMAGE_URL_WITH_CLASS = /(\w+)(\s+)((https?:\/\/|\/)(\S+))/i
        IMAGE_URL = /((https?:\/\/|\/)(\S+))/i
        def initialize(tag_name, markup, tokens)
          super
          if markup =~ IMAGE_URL_WITH_CLASS_AND_CAPTION_AND_LINK
            @class   = $1
            @url     = $3
            @caption = $7
            @link = $9
          elsif markup =~ IMAGE_URL_WITH_CLASS_AND_CAPTION
            @class   = $1
            @url     = $3
            @caption = $7
          elsif markup =~ IMAGE_URL_WITH_CAPTION
            @url     = $1
            @caption = $5
          elsif markup =~ IMAGE_URL_WITH_CLASS
            @class = $1
            @url   = $3
          elsif markup =~ IMAGE_URL
            @url = $1
          end
        end
        def render(context)
          if @class
            source = "<figure class='#{@class}'>"
          else
            source = "<figure>"
          end
          if @link
            source += "<a href=\"#{@link}\">"
          end
          source += "<img src=\"#{@url}\">"
          if @link
            source += "</a>"
          end
          source += "<figcaption>#{@caption}</figcaption>" if @caption
          source += "</figure>"
          source
        end
      end
    end
    Liquid::Template.register_tag('image', Jekyll::ImageTag)

is written as this Hugo shortcode:

    <!-- image -->
    <figure {{ with .Get "class" }}class="{{.}}"{{ end }}>
        {{ with .Get "link"}}<a href="{{.}}">{{ end }}
            <img src="{{ .Get "src" }}" {{ if or (.Get "alt") (.Get "caption") }}alt="{{ with .Get "alt"}}{{.}}{{else}}{{ .Get "caption" }}{{ end }}"{{ end }} />
        {{ if .Get "link"}}</a>{{ end }}
        {{ if or (or (.Get "title") (.Get "caption")) (.Get "attr")}}
        <figcaption>{{ if isset .Params "title" }}
            {{ .Get "title" }}{{ end }}
            {{ if or (.Get "caption") (.Get "attr")}}<p>
            {{ .Get "caption" }}
            {{ with .Get "attrlink"}}<a href="{{.}}"> {{ end }}
                {{ .Get "attr" }}
            {{ if .Get "attrlink"}}</a> {{ end }}
            </p> {{ end }}
        </figcaption>
        {{ end }}
    </figure>
    <!-- image -->

### Usage
I simply changed:

    {% image full http://farm5.staticflickr.com/4136/4829260124_57712e570a_o_d.jpg "One of my favorite touristy-type photos. I secretly waited for the good light while we were "having fun" and took this. Only regret: a stupid pole in the top-left corner of the frame I had to clumsily get rid of at post-processing." ->http://www.flickr.com/photos/alexnormand/4829260124/in/set-72157624547713078/ %}

to this (this example uses a slightly extended version named `fig`, different than the built-in `figure`):

    {{%/* fig class="full" src="http://farm5.staticflickr.com/4136/4829260124_57712e570a_o_d.jpg" title="One of my favorite touristy-type photos. I secretly waited for the good light while we were having fun and took this. Only regret: a stupid pole in the top-left corner of the frame I had to clumsily get rid of at post-processing." link="http://www.flickr.com/photos/alexnormand/4829260124/in/set-72157624547713078/" */%}}

As a bonus, the shortcode named parameters are, arguably, more readable.

## Finishing touches
### Fix content
Depending on the amount of customization that was done with each post with Jekyll, this step will require more or less effort. There are no hard and fast rules here except that `hugo server --watch` is your friend. Test your changes and fix errors as needed.

### Clean up
You'll want to remove the Jekyll configuration at this point. If you have anything else that isn't used, delete it.

## A practical example in a diff
[Hey, it's Alex](http://heyitsalex.net/) was migrated in less than a _father-with-kids day_ from Jekyll to Hugo. You can see all the changes (and screw-ups) by looking at this [diff](https://github.com/alexandre-normand/alexandre-normand/compare/869d69435bd2665c3fbf5b5c78d4c22759d7613a...b7f6605b1265e83b4b81495423294208cc74d610).

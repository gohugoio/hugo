---
title: "Shortcodes"
date: "2013-07-01"
aliases: ["/doc/shortcodes/"]
groups: ["extras"]
groups_weight: 10
---

Because Hugo uses markdown for its simple content format, however there's a lot of things that 
markdown doesn't support well.

We are unwilling to accept being constrained by our simple format. Also unacceptable is writing raw
html in our markdown every time we want to include unsupported content such as a video. To do 
so is in complete opposition to the intent of using a bare bones format for our content and 
utilizing templates to apply styling for display.

To avoid both of these limitations Hugo created shortcodes.

## What is a shortcode?
A shortcode is a simple snippet inside a markdown file that Hugo will render using a predefined template.

An example of a shortcode would be `{{% video http://urlToVideo %}}`

Shortcodes are created by placing a template file in `layouts/shortcodes/`. The
name of the file becomes the name of the shortcode (without the extension).

In your content files a shortcode can be called by using '{{&#37; name parameters
%}}' respectively. Shortcodes are space delimited (parameters with spaces
can be quoted). 

The first word is always the name of the shortcode.  Following
the name are the parameters.

The author of the shortcode can choose if the short code will use positional
parameters or named parameters (but not both). A good rule of thumb is that if
a short code has a single required value in the case of the youtube example
below then positional works very well. For more complex layouts with optional
parameters named parameters work best.

The format for named parameters models that of html with the format name="value"

Lastly like HTML, shortcodes can be singular or paired. An example of a paired
shortcode would be:

    {{% code_highlight %}} A bunch of code here {{% /code_highlight %}} 

Shortcodes are paired with an opening shortcode identical to a single shortcode
and a closing shortcode.

## Creating a shortcode

All that you need to do to create a shortcode is place a template in the layouts/shortcodes directory.

The template name will be the name of the shortcode.

**Inside the template**

To access a parameter by either position or name the index method can be used.

    {{ index .Params 0 }}
    or
    {{ index .Params "class" }}

To check if a parameter has been provided use the isset method provided by Hugo.

    {{ if isset .Params "class"}} class="{{ index .Params "class"}}" {{ end }}

For paired shortcodes the variable .Inner is available which contains all of
the content between the opening and closing shortcodes. **Simply using this
variable is the only difference between single and paired shortcodes.**

## Single Positional Example: youtube

    {{% youtube 09jf3ow9jfw %}}

Would load the template /layouts/shortcodes/youtube.html

    <div class="embed video-player">
    <iframe class="youtube-player" type="text/html" width="640" height="385" src="http://www.youtube.com/embed/{{ index .Params 0 }}" allowfullscreen frameborder="0">
    </iframe>
    </div>

This would be rendered as 

    <div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385" 
        src="http://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
    </div>

## Single Named Example: image with caption
*Example has an extra space so Hugo doesn't actually render it*

    {{ % img src="/media/spf13.jpg" title="Steve Francia" %}}

Would load the template /layouts/shortcodes/img.html

    <!-- image -->
    <figure {{ if isset .Params "class" }}class="{{ index .Params "class" }}"{{ end }}>
        {{ if isset .Params "link"}}<a href="{{ index .Params "link"}}">{{ end }}
            <img src="{{ index .Params "src" }}" {{ if or (isset .Params "alt") (isset .Params "caption") }}alt="{{ if isset .Params "alt"}}{{ index .Params "alt"}}{{else}}{{ index .Params "caption" }}{{ end }}"{{ end }} />
        {{ if isset .Params "link"}}</a>{{ end }}
        {{ if or (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")}}
        <figcaption>{{ if isset .Params "title" }}
            <h4>{{ index .Params "title" }}</h4>{{ end }}
            {{ if or (isset .Params "caption") (isset .Params "attr")}}<p>
            {{ index .Params "caption" }}
            {{ if isset .Params "attrlink"}}<a href="{{ index .Params "attrlink"}}"> {{ end }}
                {{ index .Params "attr" }}
            {{ if isset .Params "attrlink"}}</a> {{ end }}
            </p> {{ end }} 
        </figcaption>
        {{ end }}
    </figure>
    <!-- image -->

Would be rendered as:

    <figure >
        <img src="/media/spf13.jpg"  />
        <figcaption>
            <h4>Steve Francia</h4>
        </figcaption>
    </figure>

## Paired Example: Highlight
*Hugo already ships with the highlight shortcode*

*Example has an extra space so Hugo doesn't actually render it*.

    {{% highlight html %}}
    <html>
        <body> This HTML </body>
    </html>
    {{% /highlight %}}

The template for this utilizes the following code (already include in hugo)

    {{ $lang := index .Params 0 }}{{ highlight .Inner $lang }}

And will be rendered as:

    <div class="highlight" style="background: #272822"><pre style="line-height: 125%"><span style="color: #f92672">&lt;html&gt;</span>
        <span style="color: #f92672">&lt;body&gt;</span> This HTML <span style="color: #f92672">&lt;/body&gt;</span>
    <span style="color: #f92672">&lt;/html&gt;</span>
    </pre></div>

Please notice that this template makes use of a hugo specific template function
called highlight which uses pygments to add the highlighting code.

More shortcode examples can be found at [spf13.com](https://github.com/spf13/spf13.com/tree/master/layouts/shortcodes)

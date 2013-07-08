---
title: "Shortcodes"
Pubdate: "2013-07-01"
---

Because Hugo uses markdown for it's content format, it was clear that there's a lot of things that 
markdown doesn't support well. This is good, the simple nature of markdown is exactly why we chose it.

However we cannot accept being constrained by our simple format. Also unacceptable is writing raw
html in our markdown every time we want to include unsupported content such as a video. To do 
so is in complete opposition to the intent of using a bare bones format for our content and 
utilizing templates to apply styling for display.

To avoid both of these limitations Hugo has full support for shortcodes.

### What is a shortcode?
A shortcode is a simple snippet inside a markdown file that Hugo will render using a template.

Short codes are designated by the opening and closing characters of '{{&#37;' and '%}}' respectively.
Short codes are space delimited. The first word is always the name of the shortcode.  Following the 
name are the parameters. The author of the shortcode can choose if the short code
will use positional parameters or named parameters (but not both). A good rule of thumb is that if a
short code has a single required value in the case of the youtube example below then positional
works very well. For more complex layouts with optional parameters named parameters work best.

The format for named parameters models that of html with the format name="value"

### Example: youtube
*Example has an extra space so Hugo doesn't actually render it*

    {{ % youtube 09jf3ow9jfw %}}

This would be rendered as 

    <div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385" 
        src="http://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
    </div>

### Example: image with caption
*Example has an extra space so Hugo doesn't actually render it*

    {{ % img src="/media/spf13.jpg" title="Steve Francia" %}}

Would be rendered as:

    <figure >
        <img src="/media/spf13.jpg"  />
        <figcaption>
            <h4>Steve Francia</h4>
        </figcaption>
    </figure>


### Creating a shortcode

All that you need to do to create a shortcode is place a template in the layouts/shortcodes directory.

The template name will be the name of the shortcode.

**Inside the template**

To access a parameter by either position or name the index method can be used.

    {{ index .Params 0 }}
    or
    {{ index .Params "class" }}

To check if a parameter has been provided use the isset method provided by Hugo.

    {{ if isset .Params "class"}} class="{{ index .Params "class"}}" {{ end }}



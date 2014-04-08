---
title: "Highlighting"
date: "2013-07-01"
groups: ["extras"]
groups_weight: 15
---

Hugo provides the ability for you to highlight source code from within your
content. Highlighting is performed by an external python based program called
[pygments](http://pygments.org) and is triggered via an embedded shortcode. If pygments is absent from
the path, it will silently simply pass the content along unhighlighted.


## Disclaimers

 * **Warning** pygments is relatively slow and our integration isn't
speed optimized. Expect much longer build times when using highlighting
 * Languages available depends on your pygments installation.
 * While pygments supports a few different output formats and options we currently
only support output=html, style=monokai, noclasses=true, and encoding=utf-8.
 * Styles are inline in order to be supported in syndicated content when references
to style sheets are not carried over.
 * We have sought to have the simplest interface possible, which consequently
limits configuration. An ambitious user is encouraged to extend the current
functionality to offer more customization.

## Usage
Highlight takes exactly one required parameter of language and requires a
closing shortcode.

## Example
{{% highlight html %}}
    {{&#37; highlight html %}}
    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
        {{ range .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>
    {{&#37; /highlight %}}
{{% /highlight %}}


## Example Output

{{% highlight html %}}
<span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
  <span style="color: #f92672">&lt;div&gt;</span>
   <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
    {{ range .Data.Pages }}
        {{ .Render &quot;summary&quot;}}
    {{ end }}
  <span style="color: #f92672">&lt;/div&gt;</span>
<span style="color: #f92672">&lt;/section&gt;</span>
{{% /highlight %}}


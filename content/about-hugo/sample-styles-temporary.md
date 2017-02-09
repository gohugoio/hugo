---
title: Sample Styles
linktitle:
description: This should be deleted before go live.
date: 2016-02-01
publishdate: 2016-02-01
lastmod: 2016-02-01
weight: 99
draft: false
slug:
aliases:
notes:
---

## Code Block Shortcodes

### Code Input Shortcode

{{% input "my-first-template.html" %}}
```go
<ul>
  {{range $.Site.Pages}}
    <li>{{.Title}}</li>
  {{end}}
</ul>
```
{{% /input %}}

{{% input "first-script.js" %}}
```js
function myFunction(arg){
  console.log(arg);
}
```
{{% /input %}}

{{% input "config.yml" %}}
```yaml
---
archetypeDir:               "archetypes"
# hostname (and path) to the root, e.g. http://spf13.com/
baseURL:                    ""
# include content marked as draft
buildDrafts:                false
# include content with publishdate in the future
buildFuture:                false
# include content already expired
buildExpired:               false
# enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.
relativeURLs:               false
canonifyURLs:               false
...
```
{{% /input %}}

### Code Output Shortcode

{{% output "/posts/index.html" %}}
```html
<ul>
    <li>My First Post</li>
    <li>My Second Post</li>
    <li>My Third Post</li>
</ul>
```
{{% /output %}}

{{% output "/first-data.json" %}}
```json
{
    "glossary": {
        "title": "example glossary",
        "GlossDiv": {
            "title": "S"
        }
    }
}
```
{{% /output %}}

{{% output "/script.js" %}}
```javascript
$('my-div').on('click',function(evt){
    evt.preventDefault();
    console.log('Hello world!');
});
```
{{% /output %}}

### Markdown Code Block without Shortcode


```javascript
function myFunction(val){
  console.log(val);
}
```


Here is some body copy. Here is a bit of `inline code`. Here is a [link to Google](https://www.google.com). Here is some *emphasized and italics text*. Here is some **emphasized and bold text**.

## Lists

* Unordered list item lorem ipsum dolor sit amet, consectetur adipisicing elit. Ad dolorum alias vero.
* Unordered list item
    * Nested unordered list item
    * Nested unordered list item
        * Nested again
        * Nested again
* Unordered list item

1. Ordered list item lorem ipsum dolor sit amet, consectetur adipisicing elit. Sapiente, quae nulla dolore eligendi.
2. Ordered list item
    1. Nested ordered list item
    2. Nested ordered list item
        1. Nested again, ordered
        2. Nested again, ordered
3. Ordered list item

## Callouts/CTAs/Admonitions

{{% note "Note Shortcode" %}}
Here is *something in italic* in a shortcode. Here is **something in bold**.
{{% /note %}}

{{% caution "Caution Shortcode" %}}
Here is *something in italic* in a shortcode. Here is **something in bold**.
{{% /caution %}}

{{% warning "Warning Shortcode" %}}
Here is *something in italic* in a shortcode. Here is **something in bold**.
{{% /warning %}}
---
title: transform.HighlightCodeBlock
description: Highlights code received in context within a code block render hook.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: highlight.HighlightResult 
  signatures: ['transform.HighlightCodeBlock CONTEXT [OPTIONS]']
relatedFunctions:
  - transform.CanHighlight
  - transform.Highlight
  - transform.HighlightCodeBlock
---

This function is only useful within a code block render hook.

Given the context passed into a code block render hook, `transform.HighlightCodeBlock` returns a `HighlightResult` object with two methods.

.Wrapped
: (`template.HTML`) Returns highlighted code wrapped in `<div>`, `<pre>`, and `<code>` elements. This is identical to the value returned by the transform.Highlight function.

.Inner
: (`template.HTML`) Returns highlighted code without any wrapping elements, allowing you to create your own wrapper.


```go-html-template
{{ $result := transform.HighlightCodeBlock . }}
{{ $result.Wrapped }}
```

To override the default [highlighting options]:

```go-html-template
{{ $options := merge .Options (dict "linenos" true) }}
{{ $result := transform.HighlightCodeBlock . $options }}
{{ $result.Wrapped }}
```

[highlighting options]: /functions/transform/highlight/#options

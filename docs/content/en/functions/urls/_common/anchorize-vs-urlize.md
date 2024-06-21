---
# Do not remove front matter.
---

The [`anchorize`] and [`urlize`] functions are similar: 

[`anchorize`]: /functions/urls/anchorize/
[`urlize`]: /functions/urls/urlize/

- Use the `anchorize` function to generate an HTML `id` attribute value
- Use the `urlize` function to sanitize a string for usage in a URL

For example:

```go-html-template
{{ $s := "A B C" }}
{{ $s | anchorize }} → a-b-c
{{ $s | urlize }} → a-b-c

{{ $s := "a b   c" }}
{{ $s | anchorize }} → a-b---c
{{ $s | urlize }} → a-b-c

{{ $s := "< a, b, & c >" }}
{{ $s | anchorize }} → -a-b--c-
{{ $s | urlize }} → a-b-c

{{ $s := "main.go" }}
{{ $s | anchorize }} → maingo
{{ $s | urlize }} → main.go

{{ $s := "Hugö" }}
{{ $s | anchorize }} → hugö
{{ $s | urlize }} → hug%C3%B6
```

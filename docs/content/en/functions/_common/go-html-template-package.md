---
# Do not remove front matter.
---

Hugo uses Go's [text/template] and [html/template] packages.

The text/template package implements data-driven templates for generating textual output, while the html/template package implements data-driven templates for generating HTML output safe against code injection.

By default, Hugo uses the html/template package when rendering HTML files.

To generate HTML output that is safe against code injection, the html/template package escapes strings in certain contexts.

[text/template]: https://pkg.go.dev/text/template
[html/template]: https://pkg.go.dev/html/template

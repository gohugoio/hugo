---
_comment: Do not remove front matter.
---

In this project configuration we enable rendering of [emoji shortcodes][], and add emoji shortcodes before (pre) and after (post) each menu entry:

{{< code-toggle file=hugo >}}
enableEmoji = true

[[menus.main]]
name = 'About'
pageRef = '/about'
post = ':point_left:'
pre = ':point_right:'
weight = 10

[[menus.main]]
name = 'Contact'
pageRef = '/contact'
post = ':arrow_left:'
pre = ':arrow_right:'
weight = 20
{{< /code-toggle >}}

To render the menu:

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li>
      {{ .Pre | markdownify }}
      <a href="{{ .URL }}">{{ .Name }}</a>
      {{ .Post | markdownify }}
    </li>
  {{ end }}
</ul>
```

> [!NOTE]
> The `pre` and `post` values of menu entries defined in front matter are HTML escaped unless `text/html` is allowed in the [`allowContent`][] security policy. Menu entries defined in your project configuration are not escaped.

[`allowContent`]: /configuration/security/
[emoji shortcodes]: /quick-reference/emojis/

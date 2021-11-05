---
title: "transform.Transliterate"
description: "Converts a string from Unicode to ASCII."
date: 2021-11-05T01:48:00-07:00
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: ["transliterate"]
signature: ["STRING | transform.Transliterate","STRING | transliterate"]
aliases: []
---
Converts a string from Unicode to ASCII using rules predefined for your site's `defaultContentLanguage`, or using default rules if language-specific rules do not exist.

Hugo provides language-specific transliteration rules for Bosnian (bs), Bulgarian (bg), Catalan (ca), Croatian (hr), Danish (da), Esperanto (eo), German (de), Hungarian (hu), Macedonian (mk), Norwegian Bokmål (nb), Russian (ru), Serbian (sr), Slovenian (sl), Swedish (sv), and Ukrainian (uk).

For a site with English (en) as the default content language:

```go-html-template
{{ "Hugo" | transform.Transliterate }} --> Hugo
{{ "çđħłƚŧ" | transform.Transliterate }} --> cdhllt
{{ "áéíñóú" | transform.Transliterate }} --> aeinou
{{ "ÄÖÜäöüß" | transform.Transliterate }} --> AOUaouss
```

For a site with German (de) as the default content language:

```go-html-template
{{ "ÄÖÜäöüß" | transform.Transliterate }} --> AeOeUeaeoeuess
```

If you have enabled [`transliteratePath`](/getting-started/configuration/#transliteratepath) in your site configuration, you can use `transform.Transliterate` with [`.GetPage`](/functions/getpage/) to retrieve term pages:

```go-html-template
{{ with site.GetPage (path.Join "tags" ("çđħłƚŧ äöü" | transliterate | anchorize)) }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```

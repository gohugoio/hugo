---
title: transform.Transliterate
description: Returns the given string, converting Unicode to ASCII.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: string
  signatures: [transform.Transliterate INPUT]
---

The `transform.Transliterate` function converts a string from Unicode to ASCII using rules predefined for your site's `defaultContentLanguage`, or using default rules if language-specific rules do not exist.

Hugo provides language-specific transliteration rules for Bosnian (bs), Bulgarian (bg), Catalan (ca), Croatian (hr), Danish (da), Esperanto (eo), German (de), Hungarian (hu), Macedonian (mk), Norwegian Bokmål (nb), Russian (ru), Serbian (sr), Slovenian (sl), Swedish (sv), and Ukrainian (uk).

For a site with English (en) as the default content language:

```go-html-template
{{ transform.Transliterate "Hugo" }} → Hugo
{{ transform.Transliterate "çđħłƚŧ" }} → cdhllt
{{ transform.Transliterate "áéíñóú" }} → aeinou
{{ transform.Transliterate "ÄÖÜäöüß" }} → AOUaouss
```

For a site with German (de) as the default content language:

```go-html-template
{{ transform.Transliterate "ÄÖÜäöüß" }} → AeOeUeaeoeuess
```

If you have enabled [`transliteratePath`] in your site configuration, you can use `transform.Transliterate` with the [`GetPage`] method to retrieve term pages:

[`transliteratePath`]: /getting-started/configuration/#transliteratepath
[`GetPage`]: /methods/site/getpage/


```go-html-template
{{ with .Site.GetPage (path.Join "tags" (transform.Transliterate "çđħłƚŧ äöü" | anchorize)) }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```

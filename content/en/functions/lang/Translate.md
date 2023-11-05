---
title: lang.Translate
description: Translates a string using the translation tables in the i18n directory.
categories: []
keywords: []
action:
  aliases: [T, i18n]
  related: []
  returnType: string
  signatures: ['lang.Translate KEY [CONTEXT]']
aliases: [/functions/i18n]
---

Let's say your multilingual site supports two languages, English and Polish. Create a translation table for each language in the `i18n` directory.

```text
i18n/
├── en.toml
└── pl.toml
```

The translation tables can contain both:

- Simple translations
- Translations with pluralization

The Unicode [CLDR Plural Rules chart] describes the pluralization categories for each language.

[CLDR Plural Rules chart]: https://www.unicode.org/cldr/charts/43/supplemental/language_plural_rules.html

The English translation table:

{{< code-toggle file=i18n/en >}}
privacy = 'privacy'
security = 'security'

[day]
one = 'day'
other = 'days'

[day_with_count]
one = '{{ . }} day'
other = '{{ . }} days'
{{< /code-toggle >}}

The Polish translation table:

{{< code-toggle file=i18n/pl >}}
privacy = 'prywatność'
security = 'bezpieczeństwo'

[day]
one = 'miesiąc'
few = 'miesiące'
many = 'miesięcy'
other = 'miesiąca'

[day_with_count]
one = '{{ . }} miesiąc'
few = '{{ . }} miesiące'
many = '{{ . }} miesięcy'
other = '{{ . }} miesiąca'
{{< /code-toggle >}}

{{% note %}}
The examples below use the `T` alias for brevity.
{{% /note %}}

When viewing the English language site:

```go-html-template
{{ T "privacy" }} → privacy
{{ T "security" }} → security

{{ T "day" 0 }} → days
{{ T "day" 1 }} → day
{{ T "day" 2 }} → days
{{ T "day" 5 }} → days

{{ T "day_with_count" 0 }} → 0 days
{{ T "day_with_count" 1 }} → 1 day
{{ T "day_with_count" 2 }} → 2 days
{{ T "day_with_count" 5 }} → 5 days
````

When viewing the Polish language site:

```go-html-template
{{ T "privacy" }} → prywatność
{{ T "security" }} → bezpieczeństwo

{{ T "day" 0 }} → miesięcy
{{ T "day" 1 }} → miesiąc
{{ T "day" 2 }} → miesiące
{{ T "day" 5 }} → miesięcy

{{ T "day_with_count" 0 }} → 0 miesięcy
{{ T "day_with_count" 1 }} → 1 miesiąc
{{ T "day_with_count" 2 }} → 2 miesiące
{{ T "day_with_count" 5 }} → 5 miesięcy
```

In the pluralization examples above, we passed an integer in context (the second argument). You can also pass a map in context, providing a `count` key to control pluralization.

Translation table:

{{< code-toggle file=i18n/en >}}
[age]
one = '{{ .name }} is {{ .count }} year old.'
other = '{{ .name }} is {{ .count }} years old.'
{{< /code-toggle >}}

Template code:

```go-html-template
{{ T "age" (dict "name" "Will" "count" 1) }} → Will is 1 year old.
{{ T "age" (dict "name" "John" "count" 3) }} → John is 3 years old.
```

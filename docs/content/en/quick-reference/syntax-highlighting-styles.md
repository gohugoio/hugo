---
title: Syntax highlighting styles
description: Highlight code examples using one of these styles.
categories: []
keywords: [highlight]
---

## Overview

Hugo provides several methods to add syntax highlighting to code examples:

- Use the [`transform.Highlight`][] function within your templates
- Use the [`highlight`][] shortcode with any [content format](g)
- Use [fenced code blocks][] with the Markdown content format

Regardless of method, use any of the syntax highlighting styles below.

Set the default syntax highlighting style in your project configuration:

{{< code-toggle file=hugo >}}
[markup.highlight]
style = 'monokai'
{{< /code-toggle >}}

See [configure Markup][].

By default, Hugo applies syntax highlighting using inline CSS styles. To use an external CSS file instead, set [`noClasses`][] to `false` in your configuration, then generate the stylesheet with the [`hugo gen chromastyles`][] command.

## Modes

{{< new-in 0.164.0 />}}

Syntax highlighting styles are designed for either light mode or dark mode displays. Some light mode styles have a corresponding dark mode counterpart. When using external CSS stylesheets, you can switch between them based on user preference or environment.

The table below lists the light mode styles that have a dark mode counterpart. Each style name links to its entry in the [Gallery](#gallery) section of this page.

{{% syntax-highlighting-styles-light-with-counterpart %}}

## Gallery

This gallery demonstrates the application of each syntax highlighting style with code examples written in different programming languages. Each style lists its mode (light or dark). For styles with a counterpart in the opposite mode, the counterpart style is also listed.

{{% syntax-highlighting-styles %}}

[`highlight`]: /shortcodes/highlight/
[`hugo gen chromastyles`]: /commands/hugo_gen_chromastyles/
[`noClasses`]: /configuration/markup/#noclasses
[`transform.Highlight`]: /functions/transform/highlight/
[configure Markup]: /configuration/markup/#highlight
[fenced code blocks]: /content-management/syntax-highlighting/#fenced-code-blocks

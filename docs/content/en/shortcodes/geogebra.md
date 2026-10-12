---
title: GeoGebra shortcode
linkTitle: GeoGebra
description: Embed a GeoGebra material in your content using the geogebra shortcode.
categories: []
keywords: []
---

> [!NOTE]
> To override Hugo's embedded `geogebra` shortcode, copy the [source code][] to a file with the same name in the `layouts/_shortcodes` directory.

## Example

To embed a material with this URL:

```text
https://www.geogebra.org/m/zrtmrj3s
```

Include this in your Markdown:

```md
{{</* geogebra zrtmrj3s */>}}
```

The shortcode renders an iframe. By default, its container scales to the available width, up to 850 pixels, with an aspect ratio of 850:600.

## Arguments

For `width`, `height`, and `showToolBar`, shortcode arguments override the [global defaults](#global-defaults). If neither is configured, the built-in defaults apply.

`id`
: (`string`) The material ID. Required unless supplied as the first and only positional argument. Contains only letters and digits.

`width`
: (`int`) The applet width in pixels. Must be positive. Defaults to `services.geogebra.width`, which is `850`.

`height`
: (`int`) The applet height in pixels. Must be positive. Defaults to `services.geogebra.height`, which is `600`.

`class`
: (`string`) The class attribute of the wrapping div. When specified, removes the inline styles from the iframe and its wrapping div. Provide your own CSS to control sizing and borders.

`loading`
: (`string`) The iframe loading attribute, either `eager` or `lazy`. Default is `eager`.

`title`
: (`string`) The iframe title attribute. Default is `GeoGebra material`.

`allowFullScreen`
: (`bool`) Whether to allow the iframe to enter full screen mode. Default is `true`.

`showToolBar`
: (`bool`) Whether to show the toolbar. Defaults to `services.geogebra.showToolBar`, which is `false`. An explicit `false` overrides a global `true`.

The following boolean arguments default to `false`:

| Argument | Description |
| --- | --- |
| `enableRightClick` | Enable right click, zooming, and keyboard editing. |
| `enableShiftDragZoom` | Enable shift-drag and zoom. |
| `enableLabelDrags` | Allow label dragging. |
| `showAlgebraInput` | Show the input bar. |
| `showMenuBar` | Show the menu bar. |
| `showToolBarHelp` | Show toolbar help. |
| `showResetIcon` | Show the reset icon. |
| `showAnimationButton` | Show the animation play button. |
| `showFullscreenButton` | Show the full screen button. Also requires `allowFullScreen=true`. |
| `showZoomButtons` | Show zoom buttons. |

For example:

```md
{{</* geogebra id=zrtmrj3s width=800 height=600 showToolBar=true showZoomButtons=true loading=lazy */>}}
```

See GeoGebra's [material embedding documentation](https://geogebra.github.io/docs/reference/en/Material_Embedding_%28Iframe%29/) for details. This shortcode embeds existing materials; it does not load local `.ggb` files or expose the GeoGebra JavaScript API.

## Global defaults

To set the default size and toolbar visibility for all GeoGebra shortcodes:

{{< code-toggle >}}
[services.geogebra]
width = 800
height = 600
showToolBar = true
{{< /code-toggle >}}

A shortcode can override any of these values independently. For example, this inherits the configured dimensions but hides the toolbar:

```md
{{</* geogebra id=zrtmrj3s showToolBar=false */>}}
```

The built-in defaults are `width=850`, `height=600`, and `showToolBar=false`. Dimensions must be positive integers.

## Mirror configuration

To change the GeoGebra service for all uses of this shortcode, set the base URL in your site configuration:

{{< code-toggle >}}
[services.geogebra]
baseURL = 'https://ggb.example.org'
{{< /code-toggle >}}

The default is `https://www.geogebra.org`. This setting is separate from your site's `baseURL`.

The shortcode appends `/material/iframe/id/<id>/width/<width>/height/<height>/...` to the configured URL. A path prefix and a trailing slash are supported. Use an HTTP or HTTPS URL without credentials, a query string, or a fragment.

The mirror must support the material iframe endpoint and the requested material IDs. Changing this setting changes the iframe entry URL; it does not rewrite requests made inside the iframe. A mirror that only hosts `deployggb.js` and app resources is insufficient. Verify the mirror's network requests if you need to avoid connections to GeoGebra's official services.

## Privacy

To disable the shortcode across your site:

{{< code-toggle >}}
[privacy.geogebra]
disable = true
{{< /code-toggle >}}

Default is `false`.

[source code]: <{{% eturl geogebra %}}>

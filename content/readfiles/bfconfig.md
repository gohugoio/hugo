`taskLists`
: default: **`true`**<br>
    Blackfriday flag: <br>
    Purpose: `false` turns off GitHub-style automatic task/TODO list generation

`smartypants`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_USE_SMARTYPANTS`** <br>
    Purpose: `false` disables smart punctuation substitutions, including smart quotes, smart dashes, smart fractions, etc. If `true`, it may be fine-tuned with the `angledQuotes`, `fractions`, `smartDashes`, and `latexDashes` flags (see below).

`angledQuotes`
: default: **`false`**<br>
    Blackfriday flag: **`HTML_SMARTYPANTS_ANGLED_QUOTES`**<br>
    Purpose: `true` enables smart, angled double quotes. Example: "Hugo" renders to renders to «Hugo» instead of “Hugo”.

`fractions`
: default: **`true`**<br>
    Blackfriday flag: **`HTML_SMARTYPANTS_FRACTIONS`** <br>
    Purpose: <code>false</code> disables smart fractions.<br>
    Example: `5/12` renders to <sup>5</sup>&frasl;<sub>12</sub>(<code>&lt;sup&gt;5&lt;/sup&gt;&amp;frasl;&lt;sub&gt;12&lt;/sub&gt;</code>).<br> <strong>Caveat:</strong> Even with <code>fractions = false</code>, Blackfriday still converts `1/2`, `1/4`, and `3/4` respectively to ½ (<code>&amp;frac12;</code>), ¼ (<code>&amp;frac14;</code>) and ¾ (<code>&amp;frac34;</code>), but only these three.</small>

`smartDashes`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_SMARTY_DASHES`** <br>
    Purpose: `false` disables smart dashes; i.e., the conversion of multiple hyphens into an en dash or em dash. If `true`, its behavior can be modified with the `latexDashes` flag below.

`latexDashes`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_SMARTYPANTS_LATEX_DASHES`** <br>
    Purpose: `false` disables LaTeX-style smart dashes and selects conventional smart dashes. Assuming `smartDashes`: <br>
    If `true`, `--` is translated into &ndash; (`&ndash;`), whereas `---` is translated into &mdash; (`&mdash;`). <br>
    However, *spaced* single hyphen between two words is translated into an en&nbsp;dash&mdash; e.g., "`12 June - 3 July`" becomes `12 June ndash; 3 July` upon rendering.

`hrefTargetBlank`
: default: **`false`** <br>
    Blackfriday flag: **`HTML_HREF_TARGET_BLANK`** <br>
    Purpose: `true` opens external links in a new window or tab.

`plainIDAnchors`
: default **`true`** <br>
    Blackfriday flag: **`FootnoteAnchorPrefix` and `HeaderIDSuffix`** <br>
    Purpose: `true` renders any heading and footnote IDs without the document ID. <br>
    Example: renders `#my-heading` instead of `#my-heading:bec3ed8ba720b970`

`extensions`
: default: **`[]`** <br>
    Blackfriday flag: **`EXTENSION_*`** <br>
    Purpose: Enable one or more Blackfriday's Markdown extensions (if they aren't Hugo defaults). <br>
    Example: Include `hardLineBreak` in the list to enable Blackfriday's `EXTENSION_HARD_LINK_BREAK`

`extensionsmask`
: default: **`[]`** <br>
    Blackfriday flag: **`EXTENSION_*`** <br>
    Purpose: Enable one or more of Blackfriday's Markdown extensions (if they aren't Hugo defaults). <br>
    Example: Include `autoHeaderIds` in the list of disable Blackfriday's `EXTENSION_AUTO_HEADER_IDS`.
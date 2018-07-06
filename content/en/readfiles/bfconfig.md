## Blackfriday Options

`taskLists`
: default: **`true`**<br>
    Blackfriday flag: <br>
    Purpose: `false` turns off GitHub-style automatic task/TODO list generation.

`smartypants`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_USE_SMARTYPANTS`** <br>
    Purpose: `false` disables smart punctuation substitutions, including smart quotes, smart dashes, smart fractions, etc. If `true`, it may be fine-tuned with the `angledQuotes`, `fractions`, `smartDashes`, and `latexDashes` flags (see below).

`smartypantsQuotesNBSP`
: default: **`false`** <br>
    Blackfriday flag: **`HTML_SMARTYPANTS_QUOTES_NBSP`** <br>
    Purpose: `true` enables French style Guillemets with non-breaking space inside the quotes.

`angledQuotes`
: default: **`false`**<br>
    Blackfriday flag: **`HTML_SMARTYPANTS_ANGLED_QUOTES`**<br>
    Purpose: `true` enables smart, angled double quotes. Example: "Hugo" renders to «Hugo» instead of “Hugo”.

`fractions`
: default: **`true`**<br>
    Blackfriday flag: **`HTML_SMARTYPANTS_FRACTIONS`** <br>
    Purpose: <code>false</code> disables smart fractions.<br>
    Example: `5/12` renders to <sup>5</sup>&frasl;<sub>12</sub>(<code>&lt;sup&gt;5&lt;/sup&gt;&amp;frasl;&lt;sub&gt;12&lt;/sub&gt;</code>).<br> <strong>Caveat:</strong> Even with <code>fractions = false</code>, Blackfriday still converts `1/2`, `1/4`, and `3/4` respectively to ½ (<code>&amp;frac12;</code>), ¼ (<code>&amp;frac14;</code>) and ¾ (<code>&amp;frac34;</code>), but only these three.</small>

`smartDashes`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_SMARTY_DASHES`** <br>
    Purpose: `false` disables smart dashes; i.e., the conversion of multiple hyphens into an en-dash or em-dash. If `true`, its behavior can be modified with the `latexDashes` flag below.

`latexDashes`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_SMARTYPANTS_LATEX_DASHES`** <br>
    Purpose: `false` disables LaTeX-style smart dashes and selects conventional smart dashes. Assuming `smartDashes`: <br>
    If `true`, `--` is translated into &ndash; (`&ndash;`), whereas `---` is translated into &mdash; (`&mdash;`). <br>
    However, *spaced* single hyphen between two words is translated into an en&nbsp;dash&mdash; e.g., "`12 June - 3 July`" becomes `12 June &ndash; 3 July` upon rendering.

`hrefTargetBlank`
: default: **`false`** <br>
    Blackfriday flag: **`HTML_HREF_TARGET_BLANK`** <br>
    Purpose: `true` opens <s>external links</s> **absolute** links in a new window or tab. While the `target="_blank"` attribute is typically used for external links, Blackfriday does that for _all_ absolute links ([ref](https://discourse.gohugo.io/t/internal-links-in-same-tab-external-links-in-new-tab/11048/8)). One needs to make note of this if they use absolute links throughout, for internal links too (for example, by setting `canonifyURLs` to `true` or via `absURL`).

`nofollowLinks`
: default: **`false`** <br>
    Blackfriday flag: **`HTML_NOFOLLOW_LINKS`** <br>
    Purpose: `true` creates <s>external links</s> **absolute** links with `nofollow` being added to their `rel` attribute. Thereby crawlers are advised to not follow the link. While the `rel="nofollow"` attribute is typically used for external links, Blackfriday does that for _all_ absolute links. One needs to make note of this if they use absolute links throughout, for internal links too (for example, by setting `canonifyURLs` to `true` or via `absURL`).

`noreferrerLinks`
: default: **`false`** <br>
    Blackfriday flag: **`HTML_NOREFERRER_LINKS`** <br>
    Purpose: `true` creates <s>external links</s> **absolute** links with `noreferrer` being added to their `rel` attribute. Thus when following the link no referrer information will be leaked. While the `rel="noreferrer"` attribute is typically used for external links, Blackfriday does that for _all_ absolute links. One needs to make note of this if they use absolute links throughout, for internal links too (for example, by setting `canonifyURLs` to `true` or via `absURL`).

`plainIDAnchors`
: default **`true`** <br>
    Blackfriday flag: **`FootnoteAnchorPrefix` and `HeaderIDSuffix`** <br>
    Purpose: `true` renders any heading and footnote IDs without the document ID. <br>
    Example: renders `#my-heading` instead of `#my-heading:bec3ed8ba720b970`

`extensions`
: default: **`[]`** <br>
    Purpose: Enable one or more Blackfriday's Markdown extensions (**`EXTENSION_*`**). <br>
    Example: Include `hardLineBreak` in the list to enable Blackfriday's `EXTENSION_HARD_LINE_BREAK`. <br>
    *See [Blackfriday extensions](#blackfriday-extensions) section for information on all extensions.*

`extensionsmask`
: default: **`[]`** <br>
    Purpose: Disable one or more of Blackfriday's Markdown extensions (**`EXTENSION_*`**). <br>
    Example: Include `autoHeaderIds` as `false` in the list to disable Blackfriday's `EXTENSION_AUTO_HEADER_IDS`. <br>
    *See [Blackfriday extensions](#blackfriday-extensions) section for information on all extensions.*

## Blackfriday extensions

`noIntraEmphasis`
: default: *enabled* <br>
    Purpose: The "\_" character is commonly used inside words when discussing
    code, so having Markdown interpret it as an emphasis command is usually the
    wrong thing.  When enabled, Blackfriday lets you treat all emphasis markers
    as normal characters when they occur inside a word.

`tables`
: default: *enabled* <br>
    Purpose: When enabled, tables can be created by drawing them in the input
    using the below syntax:
    Example:

           Name | Age
        --------|------
            Bob | 27
          Alice | 23

`fencedCode`
: default: *enabled* <br>
    Purpose: When enabled, in addition to the normal 4-space indentation to mark
    code blocks, you can explicitly mark them and supply a language (to make
    syntax highlighting simple).

    You can use 3 or more backticks to mark the beginning of the block, and the
    same number to mark the end of the block.

    Example:

         ```md
        # Heading Level 1
        Some test
        ## Heading Level 2
        Some more test
        ```

`autolink`
: default: *enabled* <br>
    Purpose: When enabled, URLs that have not been explicitly marked as links
    will be converted into links.

`strikethrough`
: default: *enabled* <br>
    Purpose: When enabled, text wrapped with two tildes will be crossed out. <br>
    Example: `~~crossed-out~~`

`laxHtmlBlocks`
: default: *disabled* <br>
    Purpose: When enabled, loosen up HTML block parsing rules.

`spaceHeaders`
: default: *enabled* <br>
    Purpose: When enabled, be strict about prefix header rules.

`hardLineBreak`
: default: *disabled* <br>
    Purpose: When enabled, newlines in the input translate into line breaks in
    the output.


`tabSizeEight`
: default: *disabled* <br>
    Purpose: When enabled, expand tabs to eight spaces instead of four.

`footnotes`
: default: *enabled* <br>
    Purpose: When enabled, Pandoc-style footnotes will be supported.  The
    footnote marker in the text that will become a superscript text; the
    footnote definition will be placed in a list of footnotes at the end of the
    document. <br>
    Example:

        This is a footnote.[^1]

        [^1]: the footnote text.

`noEmptyLineBeforeBlock`
: default: *disabled* <br>
    Purpose: When enabled, no need to insert an empty line to start a (code,
    quote, ordered list, unordered list) block.


`headerIds`
: default: *enabled* <br>
    Purpose: When enabled, allow specifying header IDs with `{#id}`.

`titleblock`
: default: *disabled* <br>
    Purpose: When enabled, support [Pandoc-style title blocks][1].

`autoHeaderIds`
: default: *enabled* <br>
    Purpose: When enabled, auto-create the header ID's from the headline text.

`backslashLineBreak`
: default: *enabled* <br>
    Purpose: When enabled, translate trailing backslashes into line breaks.

`definitionLists`
: default: *enabled* <br>
    Purpose: When enabled, a simple definition list is made of a single-line
    term followed by a colon and the definition for that term. <br>
    Example:

        Cat
        : Fluffy animal everyone likes

        Internet
        : Vector of transmission for pictures of cats

    Terms must be separated from the previous definition by a blank line.

`joinLines`
: default: *enabled* <br>
    Purpose: When enabled, delete newlines and join the lines.

[1]: http://pandoc.org/MANUAL.html#extension-pandoc_title_block

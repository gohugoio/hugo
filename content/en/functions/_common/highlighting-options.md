---
_comment: Do not remove front matter.
---

anchorLineNos
: (`bool`) Whether to render each line number as an HTML anchor element, setting the `id` attribute of the surrounding `span` element to the line number. Irrelevant if `lineNos` is `false`. Default is `false`.

codeFences
: (`bool`) Whether to highlight fenced code blocks. Default is `true`.

guessSyntax
: (`bool`) Whether to automatically detect the language if the `LANG` argument is blank or set to a language for which there is no corresponding [lexer]. Falls back to a plain text lexer if unable to automatically detect the language. Default is `false`.

[lexer]: /getting-started/glossary/#lexer

{{% note %}}
The Chroma syntax highlighter includes lexers for approximately 250 languages, but only 5 of these have implemented automatic language detection.
{{% /note %}}

hl_Lines
: (`string`) A space-delimited list of lines to emphasize within the highlighted code. To emphasize lines 2, 3, 4, and 7, set this value to `2-4 7`. This option is independent of the `lineNoStart` option.

hl_inline
: (`bool`) Whether to render the highlighted code without a wrapping container. Default is `false`.

lineAnchors
: (`string`) When rendering a line number as an HTML anchor element, prepend this value to the `id` attribute of the surrounding `span` element. This provides unique `id` attributes when a page contains two or more code blocks. Irrelevant if `lineNos` or `anchorLineNos` is `false`.

lineNoStart
: (`int`) The number to display at the beginning of the first line. Irrelevant if `lineNos` is `false`. Default is `1`.

lineNos
: (`bool`) Whether to display a number at the beginning of each line. Default is `false`.

lineNumbersInTable
: (`bool`) Whether to render the highlighted code in an HTML table with two cells. The left table cell contains the line numbers, while the right table cell contains the code. Irrelevant if `lineNos` is `false`. Default is `true`.

noClasses
: (`bool`) Whether to use inline CSS styles instead of an external CSS file. To use an external CSS file, set this value to `false` and generate the CSS file using the `hugo gen chromastyles` command. Default is `true`.

style
: (`string`) The CSS styles to apply to the highlighted code. See the [style gallery] for examples. Case-sensitive. Default is `monokai`.

tabWidth
: (`int`) Substitute this number of spaces for each tab character in your highlighted code. Irrelevant if `noClasses` is `false`. Default is `4`.

wrapperClass
{{< new-in 0.140.2 >}}
: (`string`) The class or classes to use for the outermost element of the highlighted code. Default is `highlight`.

{{% note %}}
Instead of specifying both `lineNos` and `lineNumbersInTable`, you can use the following shorthand notation:

lineNos=inline
: equivalent to `lineNos=true` and `lineNumbersInTable=false`

lineNos=table
: equivalent to `lineNos=true` and `lineNumbersInTable=true`
{{% /note %}}

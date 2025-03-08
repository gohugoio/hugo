---
_comment: Do not remove front matter.
---

anchorLineNos
: (`bool`) Whether to render each line number as an HTML anchor element, setting the `id` attribute of the surrounding `span` element to the line number. Irrelevant if `lineNos` is `false`. Default is `false`.

codeFences
: (`bool`) Whether to highlight fenced code blocks. Default is `true`.

guessSyntax
: (`bool`) Whether to automatically detect the language if the `LANG` argument is blank or set to a language for which there is no corresponding [lexer](g). Falls back to a plain text lexer if unable to automatically detect the language. Default is `false`.

  > [!note]
  > The Chroma syntax highlighter includes lexers for approximately 250 languages, but only 5 of these have implemented automatic language detection.

hl_Lines
: (`string`) A space-delimited list of lines to emphasize within the highlighted code. To emphasize lines 2, 3, 4, and 7, set this value to `2-4 7`. This option is independent of the `lineNoStart` option.

hl_inline
: (`bool`) Whether to render the highlighted code without a wrapping container. Default is `false`.

lineAnchors
: (`string`) When rendering a line number as an HTML anchor element, prepend this value to the `id` attribute of the surrounding `span` element. This provides unique `id` attributes when a page contains two or more code blocks. Irrelevant if `lineNos` or `anchorLineNos` is `false`.

lineNoStart
: (`int`) The number to display at the beginning of the first line. Irrelevant if `lineNos` is `false`. Default is `1`.

lineNos
: (`any`) Controls line number display. Default is `false`.
  - `true`: Enable line numbers, controlled by `lineNumbersInTable`.
  - `false`: Disable line numbers.
  - `inline`: Enable inline line numbers (sets `lineNumbersInTable` to `false`).
  - `table`: Enable table-based line numbers (sets `lineNumbersInTable` to `true`).

lineNumbersInTable
: (`bool`) Whether to render the highlighted code in an HTML table with two cells. The left table cell contains the line numbers, while the right table cell contains the code. Irrelevant if `lineNos` is `false`. Default is `true`.

noClasses
: (`bool`) Whether to use inline CSS styles instead of an external CSS file. Default is `true`. To use an external CSS file, set this value to `false` and generate the CSS file from the command line:

  ```text
  hugo gen chromastyles --style=monokai > syntax.css
  ```

style
: (`string`) The CSS styles to apply to the highlighted code. Case-sensitive. Default is `monokai`. See [syntax highlighting styles].

tabWidth
: (`int`) Substitute this number of spaces for each tab character in your highlighted code. Irrelevant if `noClasses` is `false`. Default is `4`.

wrapperClass
: {{< new-in 0.140.2 />}}
: (`string`) The class or classes to use for the outermost element of the highlighted code. Default is `highlight`.

[syntax highlighting styles]: /quick-reference/syntax-highlighting-styles/

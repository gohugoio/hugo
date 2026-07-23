# Goldmark v2 (beta) — notes from a Hugo test upgrade

Re: https://github.com/yuin/goldmark/discussions/559 — upgraded Hugo to
`goldmark/v2 v2.0.0-beta.5`. Builds clean, core output matches v1 (links,
headings, auto-IDs incl. CJK de-dup, Chroma fences, footnotes, def lists, task
lists). Friction, worst first:

## Would like reconsidered

- **`Parse(source)` takes no context.** Can't pass a `parser.Context` in or read
  one back. Broke, at once: per-doc ID injection (→ construction-time
  `WithIDGenerator`), TOC-via-context (→ post-parse AST walk), and per-doc flags.
  An option to supply/return the context would fix all three.
- **`IDGenerator` is stateless-only.** Lost v1's full `parser.IDs` (`Put` +
  observing final suffixed IDs). No hook to collect IDs for a TOC; generator is
  shared across docs so can't hold state. De-dup itself is fine.
- **Attribute values are `text.MultilineValue` only.** No typed parse anymore, so
  `{lineNos=true hl_lines=[1,2]}` comes back as strings — regresses our fence
  highlight options. Also `parseAttributeWord` rejects non-letter-initial values,
  so `lineNoStart=5` now needs quoting.
- **Nodes can't hold arbitrary attr values.** We used `SetAttributeString(name,
  any)` to pass internal signals (image isBlock/ordinal) parser→renderer; now
  string-encoded.
- **Render `Hook` is whole-render only.** The new `renderer.Hook` /
  `html.WithHooks` (`PreRender`/`PostRender`) fires once per top-level
  `Render(w, src, n)`, bracketing the entire walk — not per node. That rules out
  the thing I'd most want a render hook for: recording output byte offsets around
  each node so the rendered buffer can be sliced after the fact (e.g. pull one
  table's HTML out as `b[lo:hi]`). Per-node enter/leave callbacks — or just
  exposing start/stop output offsets on nodes — would enable this.

## Extension model

- No unified `Markdown`/`Extender`; every extension split into `parser.Extension`
  and/or `html.Extension`. Fine, just a lot of surface.
- Override ordering trap: `html.New` prepends CommonMark and applies extension
  renderers *after* direct options, so a custom renderer for a built-in kind must
  be an `html.Extension`, not a `WithNodeRenderer` option, or it's clobbered. v1
  priorities made this explicit.
- Node renderers take `io.Writer` and register per-kind. Writer gets wrapped
  unless it implements `util.ErrorBufWriter` — non-obvious; had to add `Error()`
  to our writer to pass it through for the `w.(*render.Context)` assertion.

## Mechanical (pervasive but easy)

- `ast.String` gone (now owned `Text`); `**bold**` is now `ast.Strong`/
  `KindStrong` and `Emphasis` lost its `Level` (code that only handled
  `KindEmphasis`, or switched on level, silently drops bold — bit us in the TOC);
  `CodeSpan` text moved to `.Value` (no child nodes); `FencedCodeBlock`→
  `CodeBlock`+`Kind`; `Node.Type()`/`NodeType`/`TypeBlock`→`ast.BlockNode`;
  `Lines()` gone from block nodes (→`.Value`); `Remove/Replace/InsertChild`
  dropped the receiver arg; `Dump(src,level)`+`DumpHelper`→`Dump(src) *NodeDump`.
- Field types → `text.Value`/`MultilineValue`: `Link/Image.Destination/.Title`,
  `Text.Value` (was `.Segment`), `RawHTML.Value` (was `.Segments`),
  `CodeBlock.Value/.Info`.
- `AutoLink`: `AutoLinkType`/`Protocol`/`URL()`/`Label()` gone → `Destination`
  (mailto-prefixed)/`Label`/`Text`. No equivalent for our `linkifyProtocol`.
- `Writer.Write/RawWrite/SecureWrite`→`WriteText/RawWriteText/WriteHTML`; writer
  no longer on `html.Config`.
- Built-ins split (`NewXxxParser`+`NewXxxHTMLRenderer`); `NewCJK` gone
  (`WithEscapedSpace`+`WithEastAsianLineBreaks`); footnote opts moved to renderer
  + renamed; `NewTypographer`→`NewTypographerParser`; `Document.Meta()`→`Metadata()`.
- Typographer now emits Unicode chars, not entities (`&ldquo;`→`“`).

## Ecosystem

`goldmark-emoji` and `hugo-goldmark-extensions` (extras, passthrough) now have v2
releases and are wired back in, so Hugo's own extension deps are covered. Still
the usual chicken-and-egg for the wider third-party ecosystem.

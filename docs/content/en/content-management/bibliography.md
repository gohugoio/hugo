---
title: Bibliographies in Markdown
linkTitle: Bibliography
description: Include citations and a bibliography in Markdown using LaTeX markup.
categories: [content management]
keywords: [latex,pandoc,citation,reference,bibliography]
menu:
  docs:
    parent: content-management
    weight: 320
weight: 320
toc: true
---

{{< new-in 0.144.0 />}}

## Citations and Bibliographies

[Pandoc](https://pandoc.org) is a universal document converter and can be used to convert markdown files.

With **Pandoc >= 2.11**, you can use [citations](https://pandoc.org/MANUAL.html#extension-citations).
One way is to employ [BibTeX files](https://en.wikibooks.org/wiki/LaTeX/Bibliography_Management#BibTeX) to cite:

```
---
title: Citation document
---
---
bibliography: assets/bibliography.bib
...
This is a citation: @Doe2022
```

Note that Hugo will **not** pass its metadata YAML block to Pandoc; however, it will pass the **second** meta data block, denoted with `---` and `...` to Pandoc.
Thus, all Pandoc-specific settings should go there.

You can also add all elements from a bibliography file (without citing them explicitly) using:

```
---
title: My Publications
---
---
bibliography: assets/bibliography.bib
nocite: |
  @*
...
```

It is also possible to provide a custom [CSL style](https://citationstyles.org/authors/) by passing `csl: path-to-style.csl` as a Pandoc option.

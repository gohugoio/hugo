---
aliases:
- /doc/supported-formats/
lastmod: 2016-07-22
date: 2016-07-22
menu:
  main:
    parent: content
prev: /content/summaries
next: /content/multilingual
title: Markdown Extras
weight: 66
toc: false
---

Hugo provides some convenient markdown extensions.

## Task lists

Hugo supports GitHub styled task lists (TODO lists) for the Blackfriday renderer (md-files). See [Blackfriday config](/overview/configuration/#configure-blackfriday-rendering) for how to turn it off.

Example:

```markdown
- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed
```

Renders as:

- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed


And produces this HTML:

```html

<ul class="task-list">
<li><input type="checkbox" disabled="" class="task-list-item"> a task list item</li>
<li><input type="checkbox" disabled="" class="task-list-item"> list syntax required</li>
<li><input type="checkbox" disabled="" class="task-list-item"> incomplete</li>
<li><input type="checkbox" checked="" disabled="" class="task-list-item"> completed</li>
</ul>
```

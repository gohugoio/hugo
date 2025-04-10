---
_comment: Do not remove front matter.
---

<!--
This description list intentionally excludes the `pageRef` and `url` properties. Add those properties manually after using the include shortcode to include this list.
-->

identifier
: (`string`) Required when two or more menu entries have the same `name`, or when localizing the `name` using translation tables. Must start with a letter, followed by letters, digits, or underscores.

name
: (`string`) The text to display when rendering the menu entry.

params
: (`map`) User-defined properties for the menu entry.

parent
: (`string`) The `identifier` of the parent menu entry. If `identifier` is not defined, use `name`. Required for child entries in a nested menu.

post
: (`string`) The HTML to append when rendering the menu entry.

pre
: (`string`) The HTML to prepend when rendering the menu entry.

title
: (`string`) The HTML `title` attribute of the rendered menu entry.

weight
: (`int`) A non-zero integer indicating the entry's position relative the root of the menu, or to its parent for a child entry. Lighter entries float to the top, while heavier entries sink to the bottom.

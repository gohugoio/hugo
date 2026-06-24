---
_comment: Do not remove front matter.
---

When the `images` front matter parameter is set, Hugo processes each value. For internal paths, it searches page resources then global resources, using the resource permalink if found or converting the path to an absolute URL if not. External URLs are used as-is.

When `images` is not set, Hugo searches page resources for a name matching `*feature*`, falling back to `*cover*` or `*thumbnail*` if none is found. If still no image is found, Hugo uses the first entry in the site configuration's `params.images` array, if present, and processes it as described above.

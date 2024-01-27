---
# Do not remove front matter.
---

Path|Pattern|Match
:--|:--|:--
`images/foo/a.jpg`|`images/foo/*.jpg`|`true`
`images/foo/a.jpg`|`images/foo/*.*`|`true`
`images/foo/a.jpg`|`images/foo/*`|`true`
`images/foo/a.jpg`|`images/*/*.jpg`|`true`
`images/foo/a.jpg`|`images/*/*.*`|`true`
`images/foo/a.jpg`|`images/*/*`|`true`
`images/foo/a.jpg`|`*/*/*.jpg`|`true`
`images/foo/a.jpg`|`*/*/*.*`|`true`
`images/foo/a.jpg`|`*/*/*`|`true`
`images/foo/a.jpg`|`**/*.jpg`|`true`
`images/foo/a.jpg`|`**/*.*`|`true`
`images/foo/a.jpg`|`**/*`|`true`
`images/foo/a.jpg`|`**`|`true`
`images/foo/a.jpg`|`*/*.jpg`|`false`
`images/foo/a.jpg`|`*.jpg`|`false`
`images/foo/a.jpg`|`*.*`|`false`
`images/foo/a.jpg`|`*`|`false`

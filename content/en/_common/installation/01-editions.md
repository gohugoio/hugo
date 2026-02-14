---
_comment: Do not remove front matter.
---

## Editions

Hugo offers a standard edition with core features, plus extended and extended/deploy editions with more. Use the standard edition unless you need the features below.

<!-- TODO Remove the transpiler row somewhere around v0.166.0 -->

Feature|extended edition|extended/deploy edition
:--|:-:|:-:
[Transpile Sass to CSS] via embedded LibSass. Note that embedded LibSass was deprecated in v0.153.0 and will be removed in a future release. Use the [Dart Sass] transpiler instead, which is compatible with any edition.|:heavy_check_mark:|:heavy_check_mark:
Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See&nbsp;[details].|:x:|:heavy_check_mark:

[dart sass]: /functions/css/sass/#dart-sass
[transpile sass to css]: /functions/css/sass/
[details]: /host-and-deploy/deploy-with-hugo-deploy/

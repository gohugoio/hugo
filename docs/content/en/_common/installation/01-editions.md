---
_comment: Do not remove front matter.
---

## Editions

Hugo is available in several editions. Use the standard edition unless you need additional features.

Feature|standard|deploy (1)|extended|extended/deploy
:--|:-:|:-:|:-:|:-:
Core features|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Direct cloud deployment (2)|:x:|:heavy_check_mark:|:x:|:heavy_check_mark:
LibSass support (3)|:x:|:x:|:heavy_check_mark:|:heavy_check_mark:

(1) {{< new-in v0.159.2 />}}

(2) Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See&nbsp;[details].

(3) [Transpile Sass to CSS] via embedded LibSass. Note that embedded LibSass was deprecated in v0.153.0 and will be removed in a future release. Use the [Dart Sass] transpiler instead, which is compatible with any edition.

[dart sass]: /functions/css/sass/#dart-sass
[transpile sass to css]: /functions/css/sass/
[details]: /host-and-deploy/deploy-with-hugo-deploy/

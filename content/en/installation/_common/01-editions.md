---
# Do not remove front matter.
---

Hugo is available in three editions: standard, extended, and extended/deploy. While the standard edition provides core functionality, the extended and extended/deploy editions offer advanced features.

Feature|extended edition|extended/deploy edition
:--|:-:|:-:
Encode to the WebP format when [processing images]. You can decode WebP images with any edition.|:heavy_check_mark:|:heavy_check_mark:
[Transpile Sass to CSS] using the embedded LibSass transpiler. You can use the [Dart Sass] transpiler with any edition.|:heavy_check_mark:|:heavy_check_mark:
Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See [details].|:x:|:heavy_check_mark:

[dart sass]: /hugo-pipes/transpile-sass-to-css/#dart-sass
[processing images]: /content-management/image-processing/
[transpile sass to css]: /hugo-pipes/transpile-sass-to-css/
[details]: /hosting-and-deployment/hugo-deploy/

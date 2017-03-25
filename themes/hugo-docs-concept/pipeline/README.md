## Gulp-based Asset Pipeline for Hugo Docs

### Tools

You can see details for the following development dependencies in package.json, as well as details about the build process in gulpfile.js.

* Node 7.3.0
* Gulp
* Babel
* JS Uglify
* Autoprefixer
* Imagemin
* Image Resize
    * **Note:** Configured to automatically resize and optimize images in the `source-images` directory but was *not* actively used during site development.
* Concat
* Rename


### Gulpfile.js

```js
//gulpfile.js
const gulp = require('gulp');
const babel = require('gulp-babel');
const concat = require('gulp-concat');
const imagefull = 1200;
const imagehalf = 450;
const imagemin = require("gulp-imagemin");
const imageresize = require('gulp-image-resize');
const imagethumb = 80;
const os = require("os");
const parallel = require("concurrent-transform");
const plumber = require('gulp-plumber');
const pngquant = require('imagemin-pngquant');
const prefix = require('gulp-autoprefixer');
const pump = require('pump');
const rename = require("gulp-rename");
const sass = require('gulp-sass');
const sassfiles = ["./scss/**/*.scss"];
const sourcemaps = require('gulp-sourcemaps');
const uglify = require('gulp-uglify');
```
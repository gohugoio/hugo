## Welcome to the SRC folder for the Gohugo Theme.

The contents of this folder are used to generate CSS and javascript. You may never have to touch anything here,  unless you want to more deeply customize your styles.

## Tools

### NPM

We use [NPM](https://www.npmjs.com/) for package management The theme's `.gitignore` file should be kept intact to make sure that all files in the `node_modules` folder are not pushed to the repository.

### Webpack

We use Webpack to manage our asset pipeline. Arguably, Webpack is overkill for this use-case, but we're using it here because once it's set up (which we've done for you), it's really easy to use. If you want to use an external script, just add it via Yarn, and reference it in main.js. You'll find instructions in the js/main.js file.

### PostCSS
PostCSS is just CSS. You'll find `postcss.config.js` in the css folder. There you'll find that we're using [`postcss-import`](https://github.com/postcss/postcss-import) which allows us import css files directly from the node_modules folder, [`postcss-cssnext`](http://cssnext.io/features/) which gives us the power to use upcoming CSS features today. If you miss Sass you can find PostCss modules for those capabilities, too.


### Tachyons

This theme uses the [Tachyons CSS Library](http://tachyons.io/). It's about 15kb gzipped, highly modular, and each class is atomic so you never have to worry about overwriting your styles. It's a great library for themes because you can make most all the style changes you need right in your layouts.

## How to Use

You'll find the commands to run in `src/package.json`.

For development, you'll need Node with NPM installed:

```
$ cd themes/gohugo-theme/src/

$ npm install

$ npm start

```
This will process both the postcss and scripts.

For production, instead of `npm start`, run `npm run build:production,` which will output minified versions of your files.

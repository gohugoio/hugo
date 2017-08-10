module.exports = {
  plugins: {
   'postcss-import': {},
   'postcss-cssnext': {
	     browsers: ['last 2 versions', '> 5%'],
	     },
    'cssnano': {
      discardComments: {
        removeAll: true
      },
      minifyFontValues: false,
      autoprefixer: false
    }
	}
};

const webpack = require("webpack");
const path = require("path");
const CleanWebpackPlugin = require("clean-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const PurgecssPlugin = require("purgecss-webpack-plugin");
const glob = require("glob-all");

class TailwindExtractor {
	static extract(content) {
		return content.match(/[A-z0-9-:\/]+/g);
	}
}

module.exports = env => {
	console.log("NODE_ENV: ", process.env.NODE_ENV); // 'local'

	var CONFIG = {
		mode: "production",
		entry: {
			app: "./assets/index.js"			
		},
		output: {
			path: path.resolve(__dirname, "assets/output"),
			filename: "js/[name].js" // string
		},
		module: {
			rules: [				
				{
					test: /\.css$/,
					use: [
						MiniCssExtractPlugin.loader,
						{
							loader: "css-loader",
							options: {
								importLoaders: 1,
								minimize: true || {
									discardComments: {
										removeAll: true
									},
									minifyFontValues: false,
									autoprefixer: false
								}
							}
						},
						{
							loader: "postcss-loader",
							options: {
								ident: "postcss",
								plugins: () => [
									require("postcss-import"),									
									require("postcss-cssnext")({
										browsers: ['last 2 versions', '> 5%']  
									})
								]
							}
						}
					]
				},
				{
					test: /\.(woff|woff2|eot|ttf|otf|svg)$/,
					use: [
						{
							loader: "file-loader",
							options: {
								publicPath: "/fonts",
								name: "[name].[ext]",
								outputPath: "./../../static/fonts/"
								// useRelativePath: true
							}
						}
					]
				}
			]
		},
		plugins: [
			new webpack.DefinePlugin({
				ENVIRONMENT: JSON.stringify(process.env.NODE_ENV)
			}),
			new CleanWebpackPlugin(["./assets/output"], {
				root: __dirname,
				verbose: true,
				dry: false,
				allowExternal: true
			}),
			new MiniCssExtractPlugin({
				filename: "css/[name].css"
			})
		]
	};

	if (process.env.NODE_ENV == "production") {
		console.log("now with purge");
		CONFIG.plugins.push(
			new CleanWebpackPlugin(["./assets/output"], {
				root: __dirname,
				verbose: true,
				dry: false,
				allowExternal: true
			}),
			new PurgecssPlugin({
				paths: glob.sync([
					path.join(__dirname, "layouts/**/*.html")
				]),
				extractors: [
					{
						extractor: TailwindExtractor,
						extensions: ["html"]
					}
				],
				fontFace: false,
				whitelist: [
					"pagination",
					"#TableOfContents ul li",
					"chroma",
					"expand",
					"hljs"
				]
			})
		);
	}

	return CONFIG;
};

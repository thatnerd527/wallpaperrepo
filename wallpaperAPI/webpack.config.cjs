const path = require('path');
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
  entry: "./src/index.ts",
  mode: "development", // "production" or "development"
  devtool: "eval-source-map", // "source-map" or "eval-source-map"
  module: {
    rules: [
      {
        test: /\.ts?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
  plugins: [new HtmlWebpackPlugin()],
  output: {
    filename: "bundle.js",
    library: "WallpaperAPI",
    path: path.resolve(__dirname, "dist"),
  },
  devServer: {
    static: path.join(__dirname),
    compress: true,
    port: 4000,
  },
};
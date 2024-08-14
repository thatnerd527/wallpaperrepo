// Generated using webpack-cli https://github.com/webpack/webpack-cli

const path = require('path');

const isProduction = process.env.NODE_ENV == 'production';

const webpack = require("webpack");
const config = {
  entry: "./index.js",
  output: {
    path: path.resolve(__dirname, "electronload"),
    chunkFormat: "module",
  },

  plugins: [
    // Add your plugins here
    // Learn more about plugins from https://webpack.js.org/configuration/plugins/
    new webpack.IgnorePlugin({
      resourceRegExp: /\/src\/main\.ts$/,
    }),
  ],
  target: "es2020",
  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/i,
        loader: "ts-loader",
        options: {
          configFile: "tsconfig.webpack.json",
        },
        exclude: ["/node_modules/", "/src/"],
      },
      // Add your rules for custom modules here
      // Learn more about loaders from https://webpack.js.org/loaders/
    ],
  },
  resolve: {
      extensions: [".tsx", ".ts", ".jsx", ".js"],
      fallback: {
          fs: false,
          stream: false,
          url: false,
          zlib: false,
          buffer: false,
          crypto: false,
          http: false,
          net: false,
          tls: false,
          https: false,
            path: require.resolve("path-browserify"),
      }
  },
};

module.exports = () => {
    if (isProduction) {
        config.mode = 'production';
        

    } else {
        config.mode = 'development';

    }
    return config;
};

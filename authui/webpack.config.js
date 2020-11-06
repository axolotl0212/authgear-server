const path = require("path");
const webpack = require("webpack");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
  mode: "production",
  target: "browserslist",
  entry: {
    authgear: path.resolve(__dirname, "src/authgear.js"),
    "password-policy": path.resolve(__dirname, "src/password-policy.js"),
  },
  output: {
    path: path.resolve(__dirname, "../resources/authgear/static"),
    filename: "[name].js"
  },
  plugins: [
    new webpack.ProgressPlugin(),
    new MiniCssExtractPlugin({
      filename: "authgear.css"
    })
  ],
  module: {
    rules: [
      {
        test: /\.js$/,
        loader: "babel-loader"
      },
      {
        test: /.css$/,
        use: [MiniCssExtractPlugin.loader, "css-loader", "postcss-loader"]
      }
    ]
  }
};

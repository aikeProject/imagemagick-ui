const fs = require("fs");
const path = require("path");

let cssConfig = {};

if (process.env.NODE_ENV === "production") {
  cssConfig = {
    extract: {
      filename: "[name].css",
      chunkFilename: "[name].css"
    }
  };
}
const appDirectory = fs.realpathSync(process.cwd());
const resolveApp = relativePath => path.resolve(appDirectory, relativePath);

module.exports = {
  chainWebpack: config => {
    const limit = 9999999999999999;
    config.module
      .rule("images")
      .test(/\.(png|gif|jpg)(\?.*)?$/i)
      .use("url-loader")
      .loader("url-loader")
      .tap(options => Object.assign(options, { limit: limit }));
    config.module
      .rule("fonts")
      .test(/\.(woff2?|eot|ttf|otf|svg)(\?.*)?$/i)
      .use("url-loader")
      .loader("url-loader")
      .options({
        limit: limit
      });
    // 导入方式变化
    // import installElementPlus from "./plugins/element"; 这样的导入方式变为
    // import installElementPlus from "plugins/element";
    config.resolve.modules.add(resolveApp("src"));
  },
  css: cssConfig,
  configureWebpack: {
    output: {
      filename: "[name].js"
    },
    optimization: {
      splitChunks: false
    }
  },
  devServer: {
    disableHostCheck: true
  }
};

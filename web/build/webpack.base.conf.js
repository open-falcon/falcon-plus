const path = require('path')
const config = require('../config')
const projectRoot = path.join(__dirname, '../')
const projectSrc = path.join(projectRoot, 'src')
const env = process.env.NODE_ENV

module.exports = {
  entry: {
    app: './src/index.js'
  },
  output: {
    path: config.build.assetsRoot,
    publicPath: env === 'production' ? config.build.assetsPublicPath : config.dev.assetsPublicPath,
    filename: '[name].js'
  },
  resolve: {
    extensions: ['', '.js'],
    alias: {
      '~src': projectSrc,
      '~utils': path.join(projectSrc, 'utils'),
      '~coms': path.join(projectSrc, 'components'),
      '~sass': path.join(projectSrc, 'sass'), // for js
      'sass': path.join(projectSrc, 'sass'), // for scss
    }
  },

  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel',
        include: projectRoot,
        exclude: /node_modules/
      },
      {
        test: /\.scss$/,
        loaders: [
          'style',
          'css?modules&localIdentName=[local]--[hash:base64:5]&sourceMap',
          'sass?sourceMap'
        ]
      }
    ]
  }
}

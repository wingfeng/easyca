const path = require("path")

function resolve(dir) {
  return path.join(__dirname, dir)
}

module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:9000',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api' //需要rewrite重写的,
        }
      }
    }
  },
  chainWebpack: config => {
    config.resolve.alias
      .set('@/assets', resolve('src/assets'))
      .set('@/components', resolve('src/components'))
      .set('@/layout', resolve('src/layout'))
      .set('@/api', resolve('src/api'))
      .set('@/views', resolve('src/views'))
      .set('@/mixins', resolve('src/mixins'))
      .set('@/utils', resolve('src/utils'))

    config.module.rule('md')
      .test(/\.md$/)
      .use('html-loader')
      .loader('html-loader')
      .end()
      .use('markdown-loader')
      .loader('markdown-loader')
      .end()
  }
}

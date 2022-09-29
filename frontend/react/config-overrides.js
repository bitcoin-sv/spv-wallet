const webpack = require('webpack');

module.exports = {
  // The Webpack config to use when compiling your react app for development or production.
  webpack: function (config, env) {
    const okv = {
      fallback: {
        assert: require.resolve('assert'),
        buffer: require.resolve("buffer/"),
        crypto: require.resolve('crypto-browserify'),
        stream: require.resolve("stream-browserify"),
      },
    };

    const plugins = [
      new webpack.ProvidePlugin({
        process: {env: {}},
        Buffer: ['buffer', 'Buffer'],
      })
    ];

    if (!config.resolve) {
      config.resolve = okv;
    } else {
      config.resolve.fallback = { ...config.resolve.fallback, ...okv.fallback };
    }

    if (!config.plugins) {
      config.plugins = plugins;
    } else {
      config.plugins = [ ...config.plugins, ...plugins ];
    }

    return config;
  },
};

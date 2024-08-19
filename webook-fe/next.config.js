const webpack = require('webpack');

module.exports = {
  webpack: (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      'lodash-es': 'lodash',
    };
    return config;
  },
};
/* eslint-disable @typescript-eslint/no-var-requires */
/** @type {import('next').NextConfig} */
const path = require('path');
const withLess = require('next-with-less');
const withTM = require('next-transpile-modules')([
  '@arco-design/web-react',
  '@arco-themes/react-arco-pro',
]);
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

const setting = require('./src/settings.json');

module.exports = withLess(
  withTM({
    lessLoaderOptions: {
      lessOptions: {
        modifyVars: {
          'arcoblue-6': setting.themeColor,
        },
      },
    },
    webpack: (config, { isServer, buildId, dev, webpack }) => {
      const apiBaseUrl = dev
        ? '/api' // Development API URL
        : 'http://101.133.129.34:8051/v1'; // Production API URL

      const env = { 'process.env.API_URL': JSON.stringify(apiBaseUrl) };
      config.plugins.push(new webpack.DefinePlugin(env));

      if (process.env.ANALYZE) {
        config.plugins.push(
          new BundleAnalyzerPlugin({
            analyzerMode: 'static',
            reportFilename: isServer
              ? '../analyze/server.html'
              : './analyze/client.html',
            openAnalyzer: true,
          })
        );
      }
      config.module.rules.push({
        test: /\.svg$/,
        use: ['@svgr/webpack'],
      });

      config.resolve.alias['@/assets'] = path.resolve(
        __dirname,
        './src/public/assets'
      );
      config.resolve.alias['@'] = path.resolve(__dirname, './src');

      return config;
    },
    async redirects() {
      return [
        {
          source: '/',
          destination: '/video',
          permanent: true,
        },
        {
          source: '/user',
          destination: '/user/self',
          permanent: true,
        },
      ];
    },
    pageExtensions: ['tsx'],
    async rewrites() {
      return [
        {
          source: '/api/:path*',
          destination: `http://101.133.129.34:8051/v1/:path*`,
        },
        // { source: '/video/:path*', destination: `http://s348vstvo.bkt.clouddn.com/:path*` },
      ];
    },
  })
);

/* eslint-disable @typescript-eslint/no-var-requires */
/** @type {import('next').NextConfig} */
const path = require('path');
const withLess = require('next-with-less');
const withTM = require('next-transpile-modules')([
  '@arco-design/web-react',
  '@arco-themes/react-arco-pro',
]);

const setting = require("./src/settings.json");

module.exports = withLess(
  withTM({
    lessLoaderOptions: {
      lessOptions: {
        modifyVars: {
          'arcoblue-6': setting.themeColor,
        },
      },
    },
    webpack: (config) => {
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
        
      ];
    },
    pageExtensions: ['tsx'],
    async rewrites() { 
      return [ 
       //接口请求 前缀带上/api-text/
        { source: '/v1-api/:path*', destination: `http://101.133.129.34:8051/:path*` }, 
        // { source: '/video/:path*', destination: `http://s348vstvo.bkt.clouddn.com/:path*` }, 
      ]
    },
  })
);

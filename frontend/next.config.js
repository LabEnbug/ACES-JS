/** @type {import('next').NextConfig} */
const nextConfig = {
    async rewrites() { 
        return [ 
         //接口请求 前缀带上/api-text/
          { source: '/v1-api/:path*', destination: `http://101.133.129.34:8051/:path*` }, 
          // { source: '/video/:path*', destination: `http://s348vstvo.bkt.clouddn.com/:path*` }, 
        ]
      },
}

module.exports = nextConfig

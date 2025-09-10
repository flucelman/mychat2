import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import svgLoader from 'vite-svg-loader';
import path from 'path';

export default defineConfig({
  base: './',
  plugins: [
    vue(),
    svgLoader({
      svgoConfig: {
        plugins: [
          {
            name: 'preset-default',
            params: {
              overrides: {
                removeViewBox: false,
              },
            },
          },
          // 确保SVG使用currentColor
          {
            name: 'addAttributesToSVGElement',
            params: {
              attributes: [{ fill: 'currentColor' }]
            }
          }
        ]
      }
    })
  ],
  server: {
    host: '0.0.0.0', // 自定义主机名
    port: 9000, // 自定义端口
    strictPort: true,
    cors: true,
    // 临时测试：完全禁用HMR
    // hmr: false,
    hmr: {
      host: '192.168.41.125',
      port: 9001,
      clientPort: 9001,
      protocol: 'ws'
    },
  //  // 检测到/api则代理到后端
  //   proxy: {
  //    '/api': {
  //       target: 'http://47.108.218.167:8080',
  //       changeOrigin: true,
  //      }
  //    }
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    }
  },

});
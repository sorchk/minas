import { defineConfig, loadEnv } from 'vite'
import path from 'path'
import createVitePlugins from './vite/plugins';

export default ({ mode, command }) => {
  const env = loadEnv(mode, process.cwd());
  const { VITE_APP_BASE } = env;

  return defineConfig({
    plugins: createVitePlugins(env, command === 'build'),
    resolve: {
      alias: {
        vue: "vue/dist/vue.esm-bundler.js",
        'vue-i18n': "vue-i18n/dist/vue-i18n.cjs.js",
        '~': path.resolve(__dirname, './'),
        '@': path.resolve(__dirname, './src'),
      }
    },
    build: {
      cssCodeSplit: false,
      // rollupOptions: {
      //   output: {
      //     manualChunks(id) {
      //       if (id.includes('node_modules')) {
      //         return id.toString().split('node_modules/')[1].split('/')[0].toString();
      //       }
      //     }
      //   }
      // },
      outDir: '../server/www/dist',
      chunkSizeWarningLimit: 4096,
      rollupOptions: {
        output: {
          chunkFileNames: 'static/js/x-[name]-[hash].js',
          entryFileNames: 'static/js/x-[name]-[hash].js',
          assetFileNames: 'static/[ext]/x-[name]-[hash].[ext]',
          manualChunks(id) {
            if (id.includes('node_modules')) {
              return id.toString().match(/\/node_modules\/(?!.pnpm)(?<moduleName>[^\/]*)\//)?.groups!.moduleName ?? 'vender';
            }
          },
        },
      },
    },
    base: VITE_APP_BASE,
    server: {
      port: 3002,
      proxy: {
        '/minas/api/v1': {// '/v1'是代理标识，用于告诉node，url前面是/v1的就是使用代理的
          target: 'http://127.0.0.1:8002', //目标地址，一般是指后台服务器地址
          rewrite: (path) => path.replace(/^\/minas\/api\/v1/, '/minas/v1'),
          changeOrigin: true,//是否跨域
          ws: true,
        },
        '/minas/dav': {
          target: 'http://127.0.0.1:8002', //目标地址，一般是指后台服务器地址
          changeOrigin: true,//是否跨域
          ws: true,
        },
      },
    },
  })
};
import { resolve } from 'path'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  build: {
    target: 'esnext',
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        path: resolve(__dirname, 'path.html'),
      },
    },
    emptyOutDir: true,
    outDir: '../static',
  },
  plugins: [
    tailwindcss(),
    vue()
  ],
  base:"/ui/",
  preview:{
    host: "0.0.0.0",
    port:8000,
  },
  server:{
    proxy: {
      '/open': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
      '/fs': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://127.0.0.1:8080',   //代理接口
        changeOrigin: true,
        //rewrite: (path) => path.replace(/^\/api/, '')
      }
    },
    host: "0.0.0.0",
    port:8000
  },
})

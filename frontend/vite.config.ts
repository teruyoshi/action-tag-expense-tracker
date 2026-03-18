import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    allowedHosts: ['frontend'],
    watch: {
      usePolling: true,
    },
    proxy: {
      '/tags': 'http://backend:8080',
      '/events': 'http://backend:8080',
      '/expenses': 'http://backend:8080',
      '/summary': 'http://backend:8080',
      '/balance': 'http://backend:8080',
    },
  },
})

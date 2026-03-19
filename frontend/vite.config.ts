import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const enableWatch = process.env.VITE_WATCH === 'true'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    allowedHosts: ['frontend'],
    watch: enableWatch ? { usePolling: true } : null,
    hmr: enableWatch,
    proxy: {
      '/tags': 'http://backend:8080',
      '/events': 'http://backend:8080',
      '/expenses': 'http://backend:8080',
      '/summary': 'http://backend:8080',
      '/balance': 'http://backend:8080',
    },
  },
})

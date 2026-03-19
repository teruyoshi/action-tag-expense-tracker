import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const enableWatch = process.env.VITE_WATCH === 'true'

// https://vite.dev/config/
const isCI = process.env.VITE_CI === 'true' || process.env.CI === 'true'

export default defineConfig({
  plugins: [react()],
  define: {
    __IS_CI__: JSON.stringify(isCI),
  },
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

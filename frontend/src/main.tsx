import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import './index.css'
import App from './App.tsx'

declare const __IS_CI__: boolean
const isCI = __IS_CI__

const app = (
  <BrowserRouter>
    <App />
  </BrowserRouter>
)

createRoot(document.getElementById('root')!).render(isCI ? app : <StrictMode>{app}</StrictMode>)

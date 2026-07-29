export const APP_STORE_URL = 'https://apps.apple.com/app/id6758725997'
export const PLAY_STORE_URL = 'https://play.google.com/store/apps/details?id=com.talhacan.kuranai'
export const API_BASE = import.meta.env.VITE_API_BASE || 'https://kuranai-api.talhacan.com'
export const TOKEN_KEY = 'kuranai_web_tokens'
export const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID || '157832066999-3tgb2daavmrd1d1vd10ptrato1iu382t.apps.googleusercontent.com'
export const APPLE_CLIENT_ID = import.meta.env.VITE_APPLE_CLIENT_ID || 'com.talhacan.kuranai.web'
export const APPLE_REDIRECT_URI = import.meta.env.VITE_APPLE_REDIRECT_URI || 'https://kuranai.talhacan.com/auth/apple/callback'
export const WEB_DEVICE_HEADERS = {
  'X-Device-Type': 'web',
  'X-Device-Info': 'Kuran AI Web',
}

export const API_BASE = import.meta.env.VITE_API_BASE || 'https://kuranai-api.talhacan.com'
export const TOKEN_KEY = 'kuranai_manager_session'

export async function request(path, options = {}, token = '') {
  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      Accept: 'application/json',
      'X-Device-Type': 'web',
      'X-Device-Info': 'Kuran AI Manager Dashboard',
      ...(options.body ? { 'Content-Type': 'application/json' } : {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })
  const payload = await response.json().catch(() => null)
  if (!response.ok || payload?.isSuccess === false) {
    const error = new Error(payload?.message || 'İstek tamamlanamadı.')
    error.status = response.status
    throw error
  }
  return payload?.data ?? payload
}

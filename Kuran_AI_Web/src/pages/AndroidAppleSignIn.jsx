import { useEffect, useState } from 'react'
import { API_BASE, APPLE_CLIENT_ID, APPLE_REDIRECT_URI } from '../config/app'
import { loadScript } from '../auth/providerScripts'
import { useLanguage } from '../i18n/LanguageContext'

export default function AndroidAppleSignIn() {
  const { language, setLanguage, text } = useLanguage()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [completed, setCompleted] = useState(false)

  useEffect(() => {
    const requested = new URLSearchParams(window.location.search).get('lang')
    if (['tr', 'en', 'ar'].includes(requested) && requested !== language) {
      setLanguage(requested)
    }
  }, [language, setLanguage])

  const signIn = async () => {
    setLoading(true)
    setError('')
    try {
      await loadScript(
        'apple-signin-sdk',
        'https://appleid.cdn-apple.com/appleauth/static/jsapi/appleid/1/en_US/appleid.auth.js',
      )
      window.AppleID.auth.init({
        clientId: APPLE_CLIENT_ID,
        scope: 'name email',
        redirectURI: APPLE_REDIRECT_URI,
        state: crypto.randomUUID().replaceAll('-', ''),
        nonce: crypto.randomUUID().replaceAll('-', ''),
        usePopup: true,
      })
      const result = await window.AppleID.auth.signIn()
      const identityToken = result?.authorization?.id_token
      if (!identityToken) {
        throw new Error(text(
          'Apple kimlik belirteci alınamadı.',
          'Apple identity token was not returned.',
          'تعذّر الحصول على رمز هوية Apple.',
        ))
      }
      const response = await fetch(`${API_BASE}/api/user/auth/login/apple/mobile/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Device-Type': 'android',
          'X-Device-Info': 'Kuran AI Android Apple Sign In',
        },
        body: JSON.stringify({ identityToken }),
      })
      const payload = await response.json()
      if (!response.ok || !payload?.isSuccess || !payload?.data?.code) {
        throw new Error(payload?.message || text(
          'Apple ile giriş tamamlanamadı.',
          'Sign in with Apple could not be completed.',
          'تعذّر إكمال تسجيل الدخول باستخدام Apple.',
        ))
      }
      setCompleted(true)
      window.location.assign(`kuranai://auth/apple?code=${encodeURIComponent(payload.data.code)}`)
    } catch (signInError) {
      if (signInError?.error !== 'popup_closed_by_user') {
        setError(signInError?.message || text(
          'Apple ile giriş tamamlanamadı.',
          'Sign in with Apple could not be completed.',
          'تعذّر إكمال تسجيل الدخول باستخدام Apple.',
        ))
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <main style={{
      minHeight: '100dvh',
      display: 'grid',
      placeItems: 'center',
      padding: 24,
      background: '#07120f',
      color: '#fff8e7',
      textAlign: 'center',
    }}>
      <section style={{
        width: 'min(100%, 430px)',
        padding: '36px 28px',
        border: '1px solid rgba(228, 182, 45, .35)',
        borderRadius: 24,
        background: '#0a1e1a',
      }}>
        <div aria-hidden="true" style={{ color: '#e4b62d', fontSize: 44 }}>ق</div>
        <h1 style={{ marginBottom: 8 }}>Kuran AI</h1>
        <p style={{ color: '#c9c4b8', lineHeight: 1.6 }}>
          {completed
            ? text(
              'Giriş tamamlandı. Uygulamaya dönülüyor…',
              'Sign-in completed. Returning to the app…',
              'اكتمل تسجيل الدخول. جارٍ العودة إلى التطبيق…',
            )
            : text(
              'Android uygulamasında güvenli şekilde devam etmek için Apple hesabınızla giriş yapın.',
              'Sign in with your Apple Account to continue securely in the Android app.',
              'سجّل الدخول باستخدام حساب Apple للمتابعة بأمان في تطبيق Android.',
            )}
        </p>
        {!completed && (
          <button
            type="button"
            onClick={signIn}
            disabled={loading}
            style={{
              width: '100%',
              minHeight: 50,
              marginTop: 18,
              border: 0,
              borderRadius: 12,
              background: '#fff',
              color: '#050505',
              font: 'inherit',
              fontWeight: 700,
              cursor: loading ? 'wait' : 'pointer',
            }}
          >
            {loading
              ? text('Bağlanıyor…', 'Connecting…', 'جارٍ الاتصال…')
              : text('Apple ile Devam Et', 'Continue with Apple', 'المتابعة باستخدام Apple')}
          </button>
        )}
        {error && <p role="alert" style={{ color: '#ff9a8d', marginTop: 18 }}>{error}</p>}
      </section>
    </main>
  )
}

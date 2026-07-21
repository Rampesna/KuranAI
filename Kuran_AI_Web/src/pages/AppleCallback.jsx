import { useEffect } from 'react'
import { useLanguage } from '../i18n/LanguageContext'

export default function AppleCallback() {
  const { text } = useLanguage()

  useEffect(() => {
    const timer = window.setTimeout(() => window.close(), 800)
    return () => window.clearTimeout(timer)
  }, [])

  return (
    <section style={{ minHeight: '60vh', display: 'grid', placeItems: 'center', textAlign: 'center', padding: 24 }}>
      <div>
        <h1>{text('Apple ile giriş tamamlanıyor…', 'Completing Sign in with Apple…')}</h1>
        <p>{text('Bu pencereyi kapatabilirsiniz.', 'You can close this window.')}</p>
      </div>
    </section>
  )
}

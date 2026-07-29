import { useState } from 'react'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './page.module.css'

const API_BASE = import.meta.env.VITE_API_BASE || 'https://kuranai-api.talhacan.com'

export default function Contact() {
  const { text } = useLanguage()
  const [form, setForm] = useState({ name: '', email: '', subject: 'general', message: '' })
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')

  const handleChange = (e) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess(false)

    try {
      const res = await fetch(`${API_BASE}/api/public/support`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      const data = await res.json()
      if (data.isSuccess) {
        setSuccess(true)
        setForm({ name: '', email: '', subject: 'general', message: '' })
      } else {
        setError(data.message || text('Bir sorun oluştu. Lütfen yeniden deneyin.', 'Something went wrong. Please try again.'))
      }
    } catch {
      setError(text('Ağ hatası. Bağlantınızı kontrol edip yeniden deneyin.', 'Network error. Please check your connection and try again.'))
    } finally {
      setLoading(false)
    }
  }

  const isValid = form.name.trim() && form.email.trim() && form.message.trim().length >= 10

  return (
    <div className={styles.page}>
      <h1>{text('Bize Ulaşın', 'Contact Us')}</h1>
      <p className={styles.updated}>{text('Sorunuz, öneriniz veya destek talebiniz mi var? Mesajınızı gönderin; size geri dönelim.', "Have a question, suggestion, or need help? Send us a message and we'll get back to you.")}</p>

      {success ? (
        <div className={styles.success}>
          {text('✓ Mesajınız başarıyla gönderildi. En kısa sürede size dönüş yapacağız. Teşekkürler!', "✓ Your message has been sent successfully. We'll get back to you as soon as possible. Thank you!")}
        </div>
      ) : (
        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <label htmlFor="name">{text('Ad *', 'Name *')}</label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder={text('Adınız', 'Your name')}
              value={form.name}
              onChange={handleChange}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="email">{text('E-posta *', 'Email *')}</label>
            <input
              id="email"
              name="email"
              type="email"
              placeholder="your@email.com"
              value={form.email}
              onChange={handleChange}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="subject">{text('Konu', 'Subject')}</label>
            <select id="subject" name="subject" value={form.subject} onChange={handleChange}>
              <option value="general">{text('Genel Soru', 'General Question')}</option>
              <option value="bug">{text('Hata Bildirimi', 'Bug Report')}</option>
              <option value="subscription">{text('Abonelik Sorunu', 'Subscription Issue')}</option>
              <option value="account">{text('Hesap Sorunu', 'Account Issue')}</option>
              <option value="feedback">{text('Görüş / Öneri', 'Feedback / Suggestion')}</option>
              <option value="privacy">{text('Gizlilik Konusu', 'Privacy Concern')}</option>
            </select>
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="message">{text('Mesaj * (en az 10 karakter)', 'Message * (min. 10 characters)')}</label>
            <textarea
              id="message"
              name="message"
              placeholder={text('Sorunuzu veya talebinizi açıklayın…', 'Describe your issue or question...')}
              value={form.message}
              onChange={handleChange}
              required
              minLength={10}
            />
          </div>

          {error && <div className={styles.error}>{error}</div>}

          <button type="submit" className={styles.submitBtn} disabled={loading || !isValid}>
            {loading ? text('Gönderiliyor…', 'Sending...') : text('Mesajı Gönder', 'Send Message')}
          </button>
        </form>
      )}

      <h2>{text('Diğer İletişim Yolları', 'Other Ways to Reach Us')}</h2>
      <p>{text('Kuran AI App Store ve Google Play sayfaları veya uygulama içi destek seçenekleri üzerinden de bize ulaşabilirsiniz.', 'You can also reach us through the App Store and Google Play pages of Kuran AI or via the in-app support options.', 'يمكنك أيضًا التواصل معنا عبر صفحتي Kuran AI على App Store وGoogle Play أو خيارات الدعم داخل التطبيق.')}</p>

      <h2>{text('İlgili Sayfalar', 'Related Pages')}</h2>
      <p>
        <a href="/privacy-policy">{text('Gizlilik Politikası', 'Privacy Policy')}</a> · <a href="/terms">{text('Kullanım Koşulları', 'Terms of Use')}</a> · <a href="/support">{text('Destek', 'Support')}</a>
      </p>
    </div>
  )
}

import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { API_BASE, APP_STORE_URL, TOKEN_KEY } from '../config/app'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Packages.module.css'

const fallbackPackages = [
  { id: 'bronze', slug: 'bronze', name: 'Bronze', nameEn: 'Bronze', dailyQuestionLimit: 100, dailyTTSLimit: 5, monthlyPriceTL: 29.99, yearlyPriceTL: 299.90 },
  { id: 'silver', slug: 'silver', name: 'Gümüş', nameEn: 'Silver', dailyQuestionLimit: 250, dailyTTSLimit: 50, monthlyPriceTL: 49.99, yearlyPriceTL: 499.90 },
  { id: 'gold', slug: 'gold', name: 'Altın', nameEn: 'Gold', dailyQuestionLimit: null, dailyTTSLimit: 100, monthlyPriceTL: 79.99, yearlyPriceTL: 799.90 },
  { id: 'diamond', slug: 'diamond', name: 'Elmas', nameEn: 'Diamond', dailyQuestionLimit: null, dailyTTSLimit: null, monthlyPriceTL: 149.99, yearlyPriceTL: 1499.90 },
]

function readAccessToken() {
  try { return JSON.parse(localStorage.getItem(TOKEN_KEY))?.accessToken || '' }
  catch { return '' }
}

export default function Packages() {
  const { isTurkish, text } = useLanguage()
  const [packages, setPackages] = useState(fallbackPackages)
  const [overview, setOverview] = useState(null)

  useEffect(() => {
    const accessToken = readAccessToken()
    if (!accessToken) return

    const headers = { Authorization: `Bearer ${accessToken}` }
    Promise.allSettled([
      fetch(`${API_BASE}/api/v1/subscription/packages`, { headers }).then(response => response.json()),
      fetch(`${API_BASE}/api/v1/subscription/me`, { headers }).then(response => response.json()),
    ]).then(([packageResult, overviewResult]) => {
      const livePackages = packageResult.status === 'fulfilled' && packageResult.value?.isSuccess
        ? packageResult.value.data?.packages
        : null
      if (Array.isArray(livePackages) && livePackages.length) setPackages(livePackages)
      if (overviewResult.status === 'fulfilled' && overviewResult.value?.isSuccess) {
        setOverview(overviewResult.value.data)
      }
    })
  }, [])

  const currentSlug = overview?.subscription?.package?.slug
  const price = value => new Intl.NumberFormat(isTurkish ? 'tr-TR' : 'en-US', {
    style: 'currency',
    currency: 'TRY',
    maximumFractionDigits: 2,
  }).format(value)
  const limit = value => value == null || value >= 999999 ? text('Sınırsız', 'Unlimited') : value

  return (
    <div className={styles.packagesPage}>
      <section className={styles.hero}>
        <p className={styles.kicker}>{text('Kuran AI paketleri', 'Kuran AI plans')}</p>
        <h1>{text('Webde kullanın, paketinizi iOS’tan seçin.', 'Use it on the web, choose your plan on iOS.')}</h1>
        <p>{text(
          'Sohbeti web tarayıcınızdan kullanmaya devam edebilirsiniz. Abonelik satın alma ve geri yükleme işlemleri güvenli biçimde Kuran AI iOS uygulaması ve Apple App Store üzerinden yapılır.',
          'Keep using chat in your web browser. Subscription purchases and restores are handled securely through the Kuran AI iOS app and Apple App Store.',
        )}</p>
        <div className={styles.heroActions}>
          <a href={APP_STORE_URL} className={styles.primaryButton} target="_blank" rel="noopener noreferrer">
            <span aria-hidden="true"></span> {text('iOS uygulamasında aç', 'Open in the iOS app')}
          </a>
          <Link to="/app" className={styles.secondaryButton}>{text('Web sohbetine dön', 'Return to web chat')}</Link>
        </div>
      </section>

      {overview && (
        <section className={styles.usagePanel} aria-label={text('Mevcut kullanımınız', 'Your current usage')}>
          <div>
            <span>{text('Mevcut paket', 'Current plan')}</span>
            <strong>{overview.subscription?.hasPackage ? (isTurkish ? overview.subscription.package?.name : overview.subscription.package?.nameEn) : text('Ücretsiz', 'Free')}</strong>
          </div>
          <div>
            <span>{text('Bugünkü sorular', 'Questions today')}</span>
            <strong>{overview.usage?.messageUsed ?? 0} / {limit(overview.usage?.messageLimit)}</strong>
          </div>
          <div>
            <span>{text('Bugünkü sesli okuma', 'Text-to-speech today')}</span>
            <strong>{overview.usage?.ttsUsed ?? 0} / {limit(overview.usage?.ttsLimit)}</strong>
          </div>
        </section>
      )}

      <section className={styles.planGrid} aria-label={text('Paket seçenekleri', 'Plan options')}>
        {packages.map((pkg, index) => {
          const isCurrent = currentSlug === pkg.slug
          return (
            <article key={pkg.id} className={`${styles.planCard} ${pkg.slug === 'diamond' ? styles.featured : ''}`}>
              <div className={styles.planHeading}>
                <div>
                  <span className={styles.planIndex}>{String(index + 1).padStart(2, '0')}</span>
                  <h2>{isTurkish ? pkg.name : pkg.nameEn}</h2>
                </div>
                {isCurrent && <span className={styles.currentBadge}>{text('Mevcut', 'Current')}</span>}
              </div>
              <ul>
                <li><strong>{limit(pkg.dailyQuestionLimit)}</strong> {text('günlük yapay zekâ sorusu', 'AI questions per day')}</li>
                <li><strong>{limit(pkg.dailyTTSLimit)}</strong> {text('günlük sesli okuma', 'text-to-speech uses per day')}</li>
                <li>{text('Web ve iOS’ta ortak hesap', 'One account across web and iOS')}</li>
              </ul>
              <div className={styles.prices}>
                <p><span>{text('Aylık', 'Monthly')}</span><strong>{price(pkg.monthlyPriceTL)}</strong></p>
                <p><span>{text('Yıllık', 'Yearly')}</span><strong>{price(pkg.yearlyPriceTL)}</strong></p>
              </div>
              <a href={APP_STORE_URL} target="_blank" rel="noopener noreferrer" className={styles.planButton}>
                {isCurrent ? text('iOS’ta paketi yönet', 'Manage plan on iOS') : text('iOS’ta bu paketi seç', 'Choose this plan on iOS')}
              </a>
            </article>
          )
        })}
      </section>

      <section className={styles.purchaseNotice}>
        <span aria-hidden="true"></span>
        <div>
          <h2>{text('Satın alma yalnızca iOS uygulamasında', 'Purchases are available only in the iOS app')}</h2>
          <p>{text(
            'Butonlar Kuran AI’nın App Store sayfasını açar. Uygulamayı yükledikten veya açtıktan sonra Paketler ekranından aboneliğinizi seçebilirsiniz. Aynı hesapla giriş yaptığınızda paketiniz web kullanımınıza da otomatik yansır.',
            'The buttons open Kuran AI on the App Store. After installing or opening the app, choose your subscription from the Plans screen. Sign in with the same account and your plan automatically applies to web usage too.',
          )}</p>
        </div>
      </section>
    </div>
  )
}

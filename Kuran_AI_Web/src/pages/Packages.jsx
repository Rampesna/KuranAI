import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { API_BASE, APP_STORE_URL, PLAY_STORE_URL, TOKEN_KEY } from '../config/app'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Packages.module.css'

const fallbackPackages = [
  { id: 'bronze', slug: 'bronze', name: 'Bronze', nameEn: 'Bronze', dailyQuestionLimit: 100, dailyTTSLimit: 5, monthlyPriceTL: 29.99, yearlyPriceTL: 299.99 },
  { id: 'silver', slug: 'silver', name: 'Gümüş', nameEn: 'Silver', dailyQuestionLimit: 250, dailyTTSLimit: 50, monthlyPriceTL: 49.99, yearlyPriceTL: 499.99 },
  { id: 'gold', slug: 'gold', name: 'Altın', nameEn: 'Gold', dailyQuestionLimit: null, dailyTTSLimit: 100, monthlyPriceTL: 79.99, yearlyPriceTL: 799.99 },
  { id: 'diamond', slug: 'diamond', name: 'Elmas', nameEn: 'Diamond', dailyQuestionLimit: null, dailyTTSLimit: null, monthlyPriceTL: 149.99, yearlyPriceTL: 1499.99 },
]

function readAccessToken() {
  try { return JSON.parse(localStorage.getItem(TOKEN_KEY))?.accessToken || '' }
  catch { return '' }
}

export default function Packages() {
  const { language, isTurkish, isArabic, text } = useLanguage()
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
  const price = value => new Intl.NumberFormat(language === 'tr' ? 'tr-TR' : language === 'ar' ? 'ar-SA' : 'en-US', {
    style: 'currency',
    currency: 'TRY',
    maximumFractionDigits: 2,
  }).format(value)
  const limit = value => value == null || value >= 999999
    ? text('Sınırsız', 'Unlimited', 'غير محدود')
    : value
  const packageName = pkg => {
    if (!isArabic) return isTurkish ? pkg.name : pkg.nameEn
    return ({ bronze: 'برونزي', silver: 'فضي', gold: 'ذهبي', diamond: 'ماسي' })[pkg.slug] || pkg.nameEn
  }

  return (
    <div className={styles.packagesPage}>
      <section className={styles.hero}>
        <p className={styles.kicker}>{text('Kuran AI paketleri', 'Kuran AI plans', 'باقات Kuran AI')}</p>
        <h1>{text(
          'Webde kullanın, paketinizi mobil uygulamadan seçin.',
          'Use it on the web, choose your plan in the mobile app.',
          'استخدمه على الويب واختر باقتك من تطبيق الهاتف.',
        )}</h1>
        <p>{text(
          'Sohbeti web tarayıcınızdan kullanmaya devam edebilirsiniz. Abonelik satın alma ve geri yükleme işlemleri iOS’ta App Store, Android’de Google Play üzerinden güvenle yapılır.',
          'Keep using chat in your web browser. Purchases and restores are handled securely by the App Store on iOS and Google Play on Android.',
          'يمكنك متابعة المحادثة عبر الويب. تتم عمليات الشراء والاستعادة بأمان عبر App Store على iOS وGoogle Play على Android.',
        )}</p>
        <div className={styles.heroActions}>
          <a href={APP_STORE_URL} className={styles.primaryButton} target="_blank" rel="noopener noreferrer">
            {text('App Store’da aç', 'Open on the App Store', 'فتح في App Store')}
          </a>
          <a href={PLAY_STORE_URL} className={styles.primaryButton} target="_blank" rel="noopener noreferrer">
            <span aria-hidden="true">▶</span> {text('Google Play’de aç', 'Open on Google Play', 'فتح في Google Play')}
          </a>
          <Link to="/app" className={styles.secondaryButton}>{text('Web sohbetine dön', 'Return to web chat', 'العودة إلى محادثة الويب')}</Link>
        </div>
      </section>

      {overview && (
        <section className={styles.usagePanel} aria-label={text('Mevcut kullanımınız', 'Your current usage')}>
          <div>
            <span>{text('Mevcut paket', 'Current plan')}</span>
            <strong>{overview.subscription?.hasPackage ? packageName(overview.subscription.package) : text('Ücretsiz', 'Free')}</strong>
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
                  <h2>{packageName(pkg)}</h2>
                </div>
                {isCurrent && <span className={styles.currentBadge}>{text('Mevcut', 'Current')}</span>}
              </div>
              <ul>
                <li><strong>{limit(pkg.dailyQuestionLimit)}</strong> {text('günlük yapay zekâ sorusu', 'AI questions per day')}</li>
                <li><strong>{limit(pkg.dailyTTSLimit)}</strong> {text('günlük sesli okuma', 'text-to-speech uses per day')}</li>
                <li>{text(
                  'Web, iOS ve Android’de ortak hesap',
                  'One account across web, iOS, and Android',
                  'حساب واحد على الويب وiOS وAndroid',
                )}</li>
              </ul>
              <div className={styles.prices}>
                <p><span>{text('Aylık', 'Monthly')}</span><strong>{price(pkg.monthlyPriceTL)}</strong></p>
                <p><span>{text('Yıllık', 'Yearly')}</span><strong>{price(pkg.yearlyPriceTL)}</strong></p>
              </div>
              <div className={styles.planActions}>
                <a href={APP_STORE_URL} target="_blank" rel="noopener noreferrer" className={styles.planButton}>
                  {isCurrent
                    ? text('iOS’ta yönet', 'Manage on iOS', 'الإدارة على iOS')
                    : text('iOS’ta seç', 'Choose on iOS', 'الاختيار على iOS')}
                </a>
                <a href={PLAY_STORE_URL} target="_blank" rel="noopener noreferrer" className={styles.planButton}>
                  ▶ {isCurrent
                    ? text('Android’de yönet', 'Manage on Android', 'الإدارة على Android')
                    : text('Android’de seç', 'Choose on Android', 'الاختيار على Android')}
                </a>
              </div>
            </article>
          )
        })}
      </section>

      <section className={styles.purchaseNotice}>
        <span aria-hidden="true">◉</span>
        <div>
          <h2>{text(
            'Satın alma mobil uygulamalarda',
            'Purchases are available in the mobile apps',
            'الشراء متاح في تطبيقات الهاتف',
          )}</h2>
          <p>{text(
            'App Store veya Google Play’den Kuran AI’yı açın, Paketler ekranından aboneliğinizi seçin. Aynı hesapla giriş yaptığınızda paketiniz web kullanımınıza da otomatik yansır.',
            'Open Kuran AI from the App Store or Google Play and choose your subscription on the Plans screen. Sign in with the same account and your plan automatically applies to web usage too.',
            'افتح Kuran AI من App Store أو Google Play واختر اشتراكك من شاشة الباقات. عند تسجيل الدخول بالحساب نفسه تُطبّق باقتك تلقائيًا على استخدام الويب.',
          )}</p>
        </div>
      </section>
    </div>
  )
}

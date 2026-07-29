import { Link } from 'react-router-dom'
import { APP_STORE_URL, PLAY_STORE_URL } from '../config/app'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Landing.module.css'

export default function Landing() {
  const { text } = useLanguage()
  const features = [
    { icon: '✦', title: text('Yapay Zekâ Sohbeti', 'AI Chat'), desc: text('Kur’an meali, tefsir ve hadis hakkında sorularınızı kaynak odaklı bir asistana sorun.', 'Ask a source-focused assistant about Quran translations, Tafsir, and Hadith.') },
    { icon: '۞', title: text('Günün Ayeti', 'Verse of the Day'), desc: text('Her güne Kur’an-ı Kerim’den anlamlı bir ayetle başlayın.', 'Start every day with a meaningful verse from the Holy Quran.') },
    { icon: '◷', title: text('Namaz Vakitleri', 'Prayer Times'), desc: text('Konumunuza göre hesaplanan günlük namaz vakitlerini takip edin.', 'Follow daily prayer times calculated for your location.') },
    { icon: '◉', title: text('Sesli Deneyim', 'Voice Experience'), desc: text('iOS ve Android uygulamalarında sorunuzu sesle iletin ve cevapları dinleyin.', 'Ask by voice and listen to answers in the iOS and Android apps.', 'اطرح سؤالك بصوتك واستمع إلى الإجابات في تطبيقي iOS وAndroid.') },
  ]

  return (
    <div className={styles.landing}>
      <section className={styles.hero}>
        <div className={styles.heroGlow} />
        <div className={styles.heroContent}>
          <p className={styles.eyebrow}>{text('Kur’an • Tefsir • Hadis', 'Quran • Tafsir • Hadith')}</p>
          <p className={styles.badge}>{text('Yapay zekâ destekli İslami bilgi asistanı', 'AI-powered Islamic knowledge assistant')}</p>
          <h1 className={styles.title}>
            {text('Vahiy ile ', 'Speak with ')}<span className={styles.gold}>{text('konuş.', 'revelation.')}</span>
          </h1>
          <p className={styles.subtitle}>
            {text(
              'Kur’an ayetlerini, tefsirleri ve hadis kaynaklarını tek bir sade deneyimde keşfedin. Webden hemen kullanın veya mobil uygulamayı indirin.',
              'Explore Quran verses, Tafsir, and Hadith sources in one calm experience. Start on the web or download the mobile app.',
              'استكشف آيات القرآن والتفاسير ومصادر الحديث في تجربة هادئة واحدة. ابدأ عبر الويب أو نزّل تطبيق الهاتف.',
            )}
          </p>
          <div className={styles.cta}>
            <Link to="/app" className={styles.primaryBtn}>
              {text('Webde kullan', 'Use on the web')} <span>→</span>
            </Link>
            <a href={APP_STORE_URL} className={styles.downloadBtn} target="_blank" rel="noopener noreferrer">
              <svg width="18" height="21" viewBox="0 0 20 24" fill="currentColor" aria-hidden="true"><path d="M15.77 12.86c-.03-2.88 2.35-4.26 2.46-4.33-1.34-1.96-3.43-2.23-4.17-2.26-1.77-.18-3.47 1.05-4.37 1.05-.9 0-2.29-1.03-3.77-1-1.94.03-3.73 1.13-4.73 2.87-2.02 3.5-.52 8.7 1.45 11.55.96 1.39 2.11 2.96 3.62 2.9 1.45-.06 2-.94 3.75-.94s2.25.94 3.78.91c1.56-.03 2.55-1.42 3.5-2.81 1.1-1.61 1.56-3.17 1.58-3.25-.03-.01-3.04-1.17-3.07-4.63zM12.93 4.47c.8-.97 1.34-2.31 1.19-3.65-1.15.05-2.55.77-3.37 1.73-.74.85-1.38 2.22-1.21 3.53 1.28.1 2.59-.65 3.39-1.61z"/></svg>
              App Store
            </a>
            <a href={PLAY_STORE_URL} className={styles.downloadBtn} target="_blank" rel="noopener noreferrer">
              <span aria-hidden="true">▶</span>
              Google Play
            </a>
          </div>
          <p className={styles.heroNote}>{text('Ücretsiz hesapla günde 3 soru.', '3 questions per day with a free account.')}</p>
        </div>
      </section>

      <section className={styles.features}>
        <div className={styles.featuresInner}>
          <p className={styles.sectionKicker}>{text('Tek bir asistan, derin bir keşif', 'One assistant, a deeper exploration')}</p>
          <h2 className={styles.sectionTitle}>{text('Kuran AI ile neler yapabilirsiniz?', 'What can you do with Kuran AI?')}</h2>
          <div className={styles.grid}>
            {features.map(feature => (
              <article key={feature.title} className={styles.card}>
                <span className={styles.cardIcon}>{feature.icon}</span>
                <h3>{feature.title}</h3>
                <p>{feature.desc}</p>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className={styles.webSection}>
        <div className={styles.webPanel}>
          <div>
            <p className={styles.sectionKicker}>{text('Tarayıcıdan erişim', 'Access from your browser')}</p>
            <h2>{text('İndirmeden hemen sorun.', 'Ask now, without installing.')}</h2>
            <p>{text('Hesabınız iOS ve Android uygulamalarıyla ortaktır. Sohbetlerinize webden başlayın, telefonunuzdan devam edin.', 'Your account is shared with the iOS and Android apps. Start on the web and continue on your phone.', 'حسابك مشترك مع تطبيقي iOS وAndroid. ابدأ محادثتك على الويب وتابعها من هاتفك.')}</p>
          </div>
          <Link to="/app" className={styles.primaryBtn}>{text('Sohbeti aç', 'Open chat')} <span>→</span></Link>
        </div>
        <div className={styles.iosPurchaseNote}>
          <span aria-hidden="true">⇄</span>
          <p>{text('Paket satın alma işlemleri iOS’ta App Store, Android’de Google Play üzerinden yapılır; paketiniz web hesabınıza otomatik yansır.', 'Plan purchases are completed through the App Store on iOS or Google Play on Android; your plan automatically applies to your web account.', 'تتم مشتريات الباقات عبر App Store على iOS أو Google Play على Android، وتُطبّق الباقة تلقائيًا على حساب الويب.')}</p>
          <Link to="/packages">{text('Paketleri incele', 'View plans')} →</Link>
        </div>
      </section>

      <section className={styles.linksSection}>
        <div className={styles.linksInner}>
          <h2 className={styles.sectionTitle}>{text('Güven, şeffaflık ve destek', 'Trust, transparency, and support')}</h2>
          <div className={styles.linksGrid}>
            <Link to="/privacy-policy" className={styles.linkCard}><span>01</span><div><h3>{text('Gizlilik Politikası', 'Privacy Policy')}</h3><p>{text('Verilerinizi nasıl işlediğimiz', 'How we handle your data')}</p></div></Link>
            <Link to="/terms" className={styles.linkCard}><span>02</span><div><h3>{text('Kullanım Koşulları', 'Terms of Use')}</h3><p>{text('Hizmet ve abonelik koşulları', 'Service and subscription terms')}</p></div></Link>
            <Link to="/support" className={styles.linkCard}><span>03</span><div><h3>{text('Destek', 'Support')}</h3><p>{text('Sık sorulan sorular ve yardım', 'FAQs and product help')}</p></div></Link>
            <Link to="/contact" className={styles.linkCard}><span>04</span><div><h3>{text('İletişim', 'Contact')}</h3><p>{text('Bize mesaj gönderin', 'Send us a message')}</p></div></Link>
          </div>
        </div>
      </section>
    </div>
  )
}

import { Link } from 'react-router-dom'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './page.module.css'

export default function Support() {
  const { language } = useLanguage()
  const tr = language === 'tr'
  return (
    <div className={styles.page}>
      <h1>{tr ? 'Destek' : 'Support'}</h1>
      <p className={styles.updated}>{tr ? 'Kuran AI uygulaması için yardım ve destek.' : 'Help and support for the Kuran AI application.'}</p>
      <h2>{tr ? 'Nasıl Yardım Alabilirim?' : 'How to Get Help'}</h2>
      <p>{tr ? 'Uygulama, abonelik veya hesabınızla ilgili yardıma ihtiyacınız varsa ' : 'If you need help with the app, subscription, or your account, use our '}<Link to="/contact">{tr ? 'İletişim sayfasındaki' : 'Contact page'}</Link>{tr ? ' destek formunu kullanın.' : ' support form.'}</p>
      <h2>{tr ? 'Hesap ve Giriş' : 'Account & Login'}</h2>
      <ul>
        <li><strong>{tr ? 'Giriş yapamıyorum:' : 'I cannot sign in:'}</strong> {tr ? 'İlk kayıtta kullandığınız e-posta veya kullanıcı adını kullandığınızdan emin olun. Sorun sürerse hesap e-postanızla bize ulaşın.' : 'Make sure you use the email or username from your original registration. Contact us with your account email if the issue continues.'}</li>
        <li><strong>{tr ? 'Hesabımı nasıl silerim?' : 'How do I delete my account?'}</strong> {tr ? 'Profil → Güvenlik Ayarları → Hesabı Sil yolunu izleyin. Bu işlem tüm verilerinizi kalıcı olarak kaldırır.' : 'Go to Profile → Security Settings → Delete Account. This permanently removes your data.'}</li>
      </ul>
      <h2>{tr ? 'Abonelik ve Satın Almalar' : 'Subscription & Purchases'}</h2>
      <ul>
        <li><strong>{tr ? 'Aboneliğim görünmüyor:' : 'My subscription is missing:'}</strong> {tr ? 'Abonelik ekranından “Satın Alımları Geri Yükle”yi deneyin ve satın almada kullanılan Apple kimliğiyle giriş yaptığınızı kontrol edin.' : 'Try “Restore Purchases” and verify that you are signed in with the Apple ID used for the purchase.'}</li>
        <li><strong>{tr ? 'Nasıl iptal ederim?' : 'How do I cancel?'}</strong> {tr ? 'Ayarlar → Apple Kimliği → Abonelikler → Kuran AI yolunu izleyin.' : 'Go to Settings → Apple ID → Subscriptions → Kuran AI.'}</li>
        <li><strong>{tr ? 'İade istiyorum:' : 'I want a refund:'}</strong> {tr ? 'İadeleri Apple yönetir. ' : 'Apple manages refunds. Visit '}<a href="https://reportaproblem.apple.com" target="_blank" rel="noopener noreferrer">reportaproblem.apple.com</a>{tr ? ' adresinden talep oluşturun.' : ' to submit a request.'}</li>
      </ul>
      <h2>{tr ? 'Uygulama Sorunları' : 'App Issues'}</h2>
      <ul>
        <li><strong>{tr ? 'Uygulama kapanıyor:' : 'The app crashes:'}</strong> {tr ? 'En güncel sürümün yüklü olduğundan emin olun. Sorun sürerse cihaz modeliniz ve iOS sürümünüzle bize ulaşın.' : 'Install the latest version. If it continues, contact us with your device model and iOS version.'}</li>
        <li><strong>{tr ? 'Yanıtlar yavaş:' : 'AI responses are slow:'}</strong> {tr ? 'Sunucu yoğunluğundan kaynaklanabilir; kısa bir süre sonra yeniden deneyin.' : 'Server load may be the cause; try again shortly.'}</li>
        <li><strong>{tr ? 'Namaz vakitleri yanlış:' : 'Prayer times look wrong:'}</strong> {tr ? 'Cihaz ayarlarından Kuran AI için konum iznini kontrol edin.' : 'Check location permission for Kuran AI in device settings.'}</li>
      </ul>
      <h2>{tr ? 'Gizlilik ve Veriler' : 'Privacy & Data'}</h2>
      <p>{tr ? 'Toplanan veriler ve OpenAI ile veri paylaşımının ayrıntıları için ' : 'For details about collected data and processing by OpenAI, read our '}<Link to="/privacy-policy">{tr ? 'Gizlilik Politikamızı' : 'Privacy Policy'}</Link>{tr ? ' inceleyin.' : '.'}</p>
      <h2>{tr ? 'Hâlâ Yardıma mı İhtiyacınız Var?' : 'Still Need Help?'}</h2>
      <p><Link to="/contact">{tr ? 'Bize ulaşın' : 'Contact us'}</Link>{tr ? '; en kısa sürede dönüş yapalım.' : ' and we will respond as soon as possible.'}</p>
    </div>
  )
}

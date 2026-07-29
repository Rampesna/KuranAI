import { Link } from 'react-router-dom'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './page.module.css'

const copy = {
  tr: {
    title: 'Destek', updated: 'Kuran AI uygulaması için yardım ve destek.',
    how: 'Nasıl Yardım Alabilirim?', howText: 'Uygulama, abonelik veya hesabınızla ilgili yardıma ihtiyacınız varsa', contact: 'İletişim sayfasındaki', form: 'destek formunu kullanın.',
    account: 'Hesap ve Giriş',
    loginQ: 'Giriş yapamıyorum:', loginA: 'İlk kayıtta kullandığınız e-posta veya kullanıcı adını kullandığınızdan emin olun. Sorun sürerse hesap e-postanızla bize ulaşın.',
    deleteQ: 'Hesabımı nasıl silerim?', deleteA: 'Profil → Güvenlik Ayarları → Hesabı Sil yolunu izleyin. Bu işlem tüm verilerinizi kalıcı olarak kaldırır.',
    subscriptions: 'Abonelik ve Satın Almalar',
    missingQ: 'Aboneliğim görünmüyor:', missingA: 'Abonelik ekranından “Satın Alımları Geri Yükle”yi deneyin ve satın almada kullanılan Apple veya Google hesabıyla giriş yaptığınızı kontrol edin.',
    cancelQ: 'Nasıl iptal ederim?', cancelA: 'Satın aldığınız platformun App Store veya Google Play abonelik ayarlarından Kuran AI aboneliğini açın.',
    refundQ: 'İade istiyorum:', refundA: 'İadeleri satın alma platformu yönetir. App Store için reportaproblem.apple.com, Android için Google Play sipariş geçmişini kullanın.',
    issues: 'Uygulama Sorunları',
    crashQ: 'Uygulama kapanıyor:', crashA: 'En güncel sürümün yüklü olduğundan emin olun. Sorun sürerse cihaz modeliniz ve işletim sistemi sürümünüzle bize ulaşın.',
    slowQ: 'Yanıtlar yavaş:', slowA: 'Sunucu yoğunluğundan kaynaklanabilir; kısa bir süre sonra yeniden deneyin.',
    prayerQ: 'Namaz vakitleri yanlış:', prayerA: 'Cihaz ayarlarından Kuran AI için konum iznini kontrol edin.',
    privacy: 'Gizlilik ve Veriler', privacyPrefix: 'Toplanan veriler ve OpenAI ile veri paylaşımının ayrıntıları için', privacyLink: 'Gizlilik Politikamızı', privacySuffix: 'inceleyin.',
    still: 'Hâlâ Yardıma mı İhtiyacınız Var?', reach: 'Bize ulaşın', reachSuffix: '; en kısa sürede dönüş yapalım.',
  },
  en: {
    title: 'Support', updated: 'Help and support for the Kuran AI application.',
    how: 'How to Get Help', howText: 'If you need help with the app, subscription, or your account, use our', contact: 'Contact page', form: 'support form.',
    account: 'Account & Login',
    loginQ: 'I cannot sign in:', loginA: 'Make sure you use the email or username from your original registration. Contact us with your account email if the issue continues.',
    deleteQ: 'How do I delete my account?', deleteA: 'Go to Profile → Security Settings → Delete Account. This permanently removes your data.',
    subscriptions: 'Subscriptions & Purchases',
    missingQ: 'My subscription is missing:', missingA: 'Try “Restore Purchases” and verify that you are signed in with the Apple or Google account used for the purchase.',
    cancelQ: 'How do I cancel?', cancelA: 'Open Kuran AI under the App Store or Google Play subscription settings for the platform where you purchased it.',
    refundQ: 'I want a refund:', refundA: 'The purchase platform manages refunds. Use reportaproblem.apple.com for App Store or Google Play order history for Android.',
    issues: 'App Issues',
    crashQ: 'The app crashes:', crashA: 'Install the latest version. If it continues, contact us with your device model and operating-system version.',
    slowQ: 'AI responses are slow:', slowA: 'Server load may be the cause; try again shortly.',
    prayerQ: 'Prayer times look wrong:', prayerA: 'Check location permission for Kuran AI in device settings.',
    privacy: 'Privacy & Data', privacyPrefix: 'For details about collected data and processing by OpenAI, read our', privacyLink: 'Privacy Policy', privacySuffix: '.',
    still: 'Still Need Help?', reach: 'Contact us', reachSuffix: ' and we will respond as soon as possible.',
  },
  ar: {
    title: 'الدعم', updated: 'المساعدة والدعم لتطبيق Kuran AI.',
    how: 'كيف أحصل على المساعدة؟', howText: 'إذا احتجت إلى مساعدة بشأن التطبيق أو الاشتراك أو الحساب، فاستخدم', contact: 'صفحة التواصل', form: 'ونموذج الدعم فيها.',
    account: 'الحساب وتسجيل الدخول',
    loginQ: 'لا أستطيع تسجيل الدخول:', loginA: 'تأكد من استخدام البريد الإلكتروني أو اسم المستخدم الذي سجلت به أول مرة. إذا استمرت المشكلة فتواصل معنا باستخدام بريد حسابك.',
    deleteQ: 'كيف أحذف حسابي؟', deleteA: 'انتقل إلى الملف الشخصي ← إعدادات الأمان ← حذف الحساب. يؤدي ذلك إلى إزالة جميع بياناتك نهائيًا.',
    subscriptions: 'الاشتراكات والمشتريات',
    missingQ: 'اشتراكي غير ظاهر:', missingA: 'جرّب «استعادة المشتريات» وتأكد من تسجيل الدخول بحساب Apple أو Google الذي تم الشراء به.',
    cancelQ: 'كيف ألغي الاشتراك؟', cancelA: 'افتح اشتراك Kuran AI من إعدادات اشتراكات App Store أو Google Play بحسب منصة الشراء.',
    refundQ: 'أريد استرداد المبلغ:', refundA: 'تدير منصة الشراء عمليات الاسترداد. استخدم reportaproblem.apple.com لـ App Store أو سجل طلبات Google Play على Android.',
    issues: 'مشكلات التطبيق',
    crashQ: 'يتوقف التطبيق:', crashA: 'تأكد من تثبيت أحدث إصدار. إذا استمرت المشكلة فتواصل معنا مع ذكر طراز جهازك وإصدار نظام التشغيل.',
    slowQ: 'إجابات الذكاء الاصطناعي بطيئة:', slowA: 'قد يكون السبب ضغط الخادم؛ حاول مرة أخرى بعد قليل.',
    prayerQ: 'مواقيت الصلاة غير صحيحة:', prayerA: 'تحقق من إذن الموقع لتطبيق Kuran AI في إعدادات الجهاز.',
    privacy: 'الخصوصية والبيانات', privacyPrefix: 'للتعرّف إلى البيانات التي نجمعها وكيفية معالجتها بواسطة OpenAI، اقرأ', privacyLink: 'سياسة الخصوصية', privacySuffix: '.',
    still: 'ما زلت بحاجة إلى مساعدة؟', reach: 'تواصل معنا', reachSuffix: ' وسنرد عليك في أقرب وقت ممكن.',
  },
}

export default function Support() {
  const { language } = useLanguage()
  const c = copy[language]
  return (
    <div className={styles.page}>
      <h1>{c.title}</h1>
      <p className={styles.updated}>{c.updated}</p>
      <h2>{c.how}</h2>
      <p>{c.howText} <Link to="/contact">{c.contact}</Link> {c.form}</p>
      <h2>{c.account}</h2>
      <ul>
        <li><strong>{c.loginQ}</strong> {c.loginA}</li>
        <li><strong>{c.deleteQ}</strong> {c.deleteA}</li>
      </ul>
      <h2>{c.subscriptions}</h2>
      <ul>
        <li><strong>{c.missingQ}</strong> {c.missingA}</li>
        <li><strong>{c.cancelQ}</strong> {c.cancelA}</li>
        <li><strong>{c.refundQ}</strong> {c.refundA}</li>
      </ul>
      <h2>{c.issues}</h2>
      <ul>
        <li><strong>{c.crashQ}</strong> {c.crashA}</li>
        <li><strong>{c.slowQ}</strong> {c.slowA}</li>
        <li><strong>{c.prayerQ}</strong> {c.prayerA}</li>
      </ul>
      <h2>{c.privacy}</h2>
      <p>{c.privacyPrefix} <Link to="/privacy-policy">{c.privacyLink}</Link>{c.privacySuffix}</p>
      <h2>{c.still}</h2>
      <p><Link to="/contact">{c.reach}</Link>{c.reachSuffix}</p>
    </div>
  )
}

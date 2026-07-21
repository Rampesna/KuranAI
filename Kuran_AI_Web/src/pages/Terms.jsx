import { useLanguage } from '../i18n/LanguageContext'
import styles from './page.module.css'

const copy = {
  tr: {
    title: 'Kullanım Koşulları (EULA)', updated: 'Son güncelleme: 21 Temmuz 2026. Kuran AI uygulaması ve bağlantılı hizmetler için geçerlidir.',
    sections: [
      ['1. Koşulların Kabulü', 'Kuran AI uygulamasını indirerek, kurarak veya kullanarak bu koşulları kabul etmiş olursunuz. Kabul etmiyorsanız uygulamayı kullanmayın.'],
      ['2. Hizmetin Tanımı', 'Kuran AI; Kur’an meali, tefsir ve hadis hakkında yapay zekâ destekli sohbet, sesli soru-cevap, günlük ayet ve namaz vakti özellikleri sunar. Hizmet kesintisiz veya hatasız çalışma garantisi olmaksızın “olduğu gibi” sağlanır.'],
      ['3. Hesap ve Güvenlik', 'Hesap bilgilerinizin güvenliğini korumaktan siz sorumlusunuz. Yetkisiz erişim şüphenizde bizimle hemen iletişime geçin.'],
      ['4. Abonelikler ve Ödemeler', 'Kuran AI, daha yüksek kullanım limitleri sağlayan otomatik yenilenen Bronze, Silver, Gold ve Diamond abonelikleri sunabilir. Ücret Apple kimliğinizden tahsil edilir. Yenilemeyi dönem bitiminden en az 24 saat önce cihazınızın Ayarlar → Apple Kimliği → Abonelikler bölümünden kapatabilirsiniz. İade talepleri Apple tarafından değerlendirilir.'],
      ['5. Yapay Zekâ İçeriği', 'Yanıtlar yalnızca bilgilendirme ve eğitim amaçlıdır; fetva veya yetkili dinî hüküm değildir. Dinî rehberlik için yetkin uzmanlara danışın. Yapay zekâ yanıtlarının doğruluğu veya eksiksizliği garanti edilmez.'],
      ['6. Kabul Edilebilir Kullanım', 'Uygulamayı yalnızca hukuka uygun amaçlarla kullanabilirsiniz. Tersine mühendislik, toplu otomatik istek, hizmeti bozma veya yasa dışı kullanım yasaktır ve hesabın sınırlandırılmasına yol açabilir.'],
      ['7. Fikrî Mülkiyet', 'Uygulama, tasarım, logo, kod ve ilgili içerikler geliştiricinin mülkiyetindedir. İzinsiz kopyalama, dağıtım veya ticari kullanım yasaktır.'],
      ['8. Üçüncü Taraf Hizmetleri', 'Yapay zekâ sohbeti ve sesli transkripsiyon için OpenAI API kullanılır. Veri işleme ayrıntıları için Gizlilik Politikamızı inceleyin.'],
      ['9. Sorumluluğun Sınırlandırılması', 'Yasaların izin verdiği en geniş ölçüde geliştirici, uygulamanın kullanımı veya kullanılamamasından doğan dolaylı ya da sonuçsal zararlardan sorumlu değildir.'],
      ['10. Geçerli Hukuk ve İletişim', 'Bu koşullar Türkiye Cumhuriyeti hukukuna tabidir. Sorularınız için İletişim sayfamızı veya uygulama içi destek kanalını kullanabilirsiniz.'],
    ],
  },
  en: {
    title: 'Terms of Use (EULA)', updated: 'Last updated: July 21, 2026. Applies to the Kuran AI application and related services.',
    sections: [
      ['1. Acceptance', 'By downloading, installing, or using Kuran AI, you agree to these terms. If you do not agree, do not use the application.'],
      ['2. Service', 'Kuran AI provides AI-assisted chat about Quran translations, Tafsir, and Hadith, voice questions, a daily verse, and prayer-time calculations. The service is provided “as is” without a guarantee of uninterrupted or error-free operation.'],
      ['3. Account and Security', 'You are responsible for protecting your account credentials. Contact us immediately if you suspect unauthorized access.'],
      ['4. Subscriptions and Payments', 'Kuran AI may offer auto-renewable Bronze, Silver, Gold, and Diamond subscriptions with higher limits. Apple charges your Apple ID. You can disable renewal at least 24 hours before the period ends under Settings → Apple ID → Subscriptions. Apple handles refund requests.'],
      ['5. AI-Generated Content', 'Responses are for informational and educational purposes only and are not authoritative religious rulings. Consult qualified scholars for religious guidance. Accuracy or completeness is not guaranteed.'],
      ['6. Acceptable Use', 'You may use the application only for lawful purposes. Reverse engineering, bulk automated requests, service disruption, and illegal use are prohibited and may result in account restrictions.'],
      ['7. Intellectual Property', 'The application, design, logos, code, and related content belong to the developer. Unauthorized copying, distribution, or commercial use is prohibited.'],
      ['8. Third-Party Services', 'OpenAI’s API powers AI chat and voice transcription. See our Privacy Policy for data-processing details.'],
      ['9. Limitation of Liability', 'To the maximum extent allowed by law, the developer is not liable for indirect or consequential damages arising from use of or inability to use the application.'],
      ['10. Governing Law and Contact', 'These terms are governed by the laws of the Republic of Turkey. Use our Contact page or in-app support channel for questions.'],
    ],
  },
}

export default function Terms() {
  const { language } = useLanguage()
  const content = copy[language]
  return (
    <div className={styles.page}>
      <h1>{content.title}</h1>
      <p className={styles.updated}>{content.updated}</p>
      {content.sections.map(([title, body]) => <section key={title}><h2>{title}</h2><p>{body}</p></section>)}
      <p>{language === 'tr' ? 'Apple üzerinden yapılan ödemeler ' : 'Payments made through Apple are subject to Apple’s '}<a href="https://www.apple.com/legal/internet-services/itunes/dev/stdeula/" target="_blank" rel="noopener noreferrer">{language === 'tr' ? 'Apple Standart Son Kullanıcı Lisans Sözleşmesine' : 'Standard End User License Agreement'}</a>{language === 'tr' ? ' tabidir.' : '.'}</p>
    </div>
  )
}

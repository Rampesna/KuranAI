import { useLanguage } from '../i18n/LanguageContext'
import styles from './page.module.css'

export default function Privacy() {
  const { language } = useLanguage()

  if (language === 'ar') {
    return (
      <div className={styles.page}>
        <h1>سياسة الخصوصية</h1>
        <p className={styles.updated}>آخر تحديث: 29 يوليو 2026. تنطبق على تطبيق Kuran AI والخدمات المرتبطة به.</p>
        <h2>1. البيانات التي نجمعها وأسباب جمعها</h2>
        <p>يجمع Kuran AI الحد الأدنى من البيانات اللازمة لتقديم خدماته:</p>
        <ul>
          <li><strong>معلومات الحساب:</strong> يُستخدم البريد الإلكتروني والاسم وبيانات الملف الشخصي للمصادقة وإدارة الحساب.</li>
          <li><strong>صورة الملف الشخصي:</strong> اختيارية، وتُحفظ فقط لعرضها في ملفك الشخصي.</li>
          <li><strong>رسائل المحادثة:</strong> تُحفظ الأسئلة والإجابات المنشأة بالذكاء الاصطناعي لتوفير سجل المحادثات. يُرسل محتوى الرسالة إلى واجهة OpenAI API لإنشاء الإجابة.</li>
          <li><strong>الإدخال الصوتي:</strong> يُرسل التسجيل الصوتي عبر خادمنا إلى OpenAI لتحويله إلى نص، ولا يحتفظ به Kuran AI بصورة دائمة.</li>
          <li><strong>الموقع:</strong> يُستخدم فقط أثناء تشغيل التطبيق لحساب مواقيت الصلاة المحلية، ولا يُحفظ على خوادمنا.</li>
          <li><strong>بيانات الاستخدام:</strong> يُتابَع استخدام الميزات يوميًا لإدارة حدود الباقات.</li>
          <li><strong>بيانات الشراء:</strong> تتحقق Apple أو Google من الاشتراكات. لا نستلم بيانات بطاقتك أو وسيلة الدفع.</li>
        </ul>
        <h2>2. خدمة الذكاء الاصطناعي التابعة لطرف ثالث (OpenAI)</h2>
        <p>بعد موافقتك الصريحة، تعالج OpenAI, Inc. نص الرسالة أو التسجيل الصوتي لإنشاء إجابة أو نص مكتوب. لا نرفق اسمك أو بريدك الإلكتروني بهذه الطلبات. تعالج OpenAI هذه البيانات وفق <a href="https://openai.com/policies/privacy-policy" target="_blank" rel="noopener noreferrer">سياسة الخصوصية</a> و<a href="https://openai.com/policies/terms-of-use" target="_blank" rel="noopener noreferrer">شروط الاستخدام</a> الخاصة بها. يمكنك الرفض ومواصلة استخدام الميزات التي لا تتطلب الذكاء الاصطناعي.</p>
        <h2>3. مشاركة البيانات</h2>
        <p><strong>لا نبيع بياناتك الشخصية ولا نؤجرها ولا نشاركها مع المعلنين أو وسطاء البيانات.</strong> لا تُرسل البيانات إلا إلى مزوّدي الخدمة اللازمين لتشغيل المنتج، بمن فيهم OpenAI كما هو موضح أعلاه، وApple أو Google للمصادقة والتحقق من الاشتراكات.</p>
        <h2>4. التخزين والأمان</h2>
        <p>تُنقل البيانات عبر HTTPS/TLS. الوصول إلى خوادمنا وقواعد بياناتنا مقيّد ومحمي بضوابط متوافقة مع معايير الصناعة.</p>
        <h2>5. الاحتفاظ بالبيانات وحذف الحساب</h2>
        <p>نحتفظ ببيانات الحساب ما دام حسابك نشطًا. يمكنك حذف حسابك والبيانات المرتبطة به نهائيًا من <strong>الملف الشخصي ← إعدادات الأمان ← حذف الحساب</strong>.</p>
        <h2>6. خصوصية الأطفال</h2>
        <p>لا يستهدف Kuran AI الأطفال دون 13 عامًا، ولا يجمع بياناتهم الشخصية عن علم.</p>
        <h2>7. التغييرات والتواصل</h2>
        <p>سنبلغك بالتغييرات الجوهرية عبر التطبيق أو البريد الإلكتروني. لأسئلة الخصوصية استخدم <a href="/contact">صفحة التواصل</a>.</p>
      </div>
    )
  }

  if (language === 'en') {
    return (
      <div className={styles.page}>
        <h1>Privacy Policy</h1>
        <p className={styles.updated}>Last updated: July 29, 2026. Applies to the Kuran AI application and related services.</p>
        <h2>1. Data We Collect and Why</h2>
        <p>Kuran AI collects only the minimum data necessary to provide its services:</p>
        <ul>
          <li><strong>Account information:</strong> Email address, name, and profile details are used for authentication and account management.</li>
          <li><strong>Profile photo:</strong> Optional. Stored only for display in your profile.</li>
          <li><strong>Chat messages:</strong> Questions and AI-generated responses are stored to provide chat history. Message content is sent to OpenAI's API for response generation.</li>
          <li><strong>Voice input:</strong> Audio is sent through our server to OpenAI for transcription and is not permanently stored by Kuran AI.</li>
          <li><strong>Location:</strong> Used only while the app is in use to calculate local prayer times and is not stored on our servers.</li>
          <li><strong>Usage data:</strong> Daily feature usage is tracked to manage plan limits.</li>
          <li><strong>Purchase data:</strong> Apple or Google verifies subscriptions. We do not receive card or payment-method details.</li>
        </ul>
        <h2>2. Third-Party AI Service (OpenAI)</h2>
        <p>With your explicit permission, message text and voice recordings are processed by OpenAI, Inc. to generate chat responses or transcripts. We do not include your name or email in these requests. OpenAI processes this data under its <a href="https://openai.com/policies/privacy-policy" target="_blank" rel="noopener noreferrer">Privacy Policy</a> and <a href="https://openai.com/policies/terms-of-use" target="_blank" rel="noopener noreferrer">Terms of Use</a>. You may decline and continue using features that do not require AI.</p>
        <h2>3. Data Sharing</h2>
        <p><strong>We do not sell, rent, or share your personal data with advertisers or data brokers.</strong> Data is shared only with service providers required to operate the product, including OpenAI as described above and Apple or Google for authentication and subscription verification.</p>
        <h2>4. Storage and Security</h2>
        <p>Data is transmitted over HTTPS/TLS. Access to our servers and databases is restricted and protected using industry-standard controls.</p>
        <h2>5. Retention and Account Deletion</h2>
        <p>We retain account data while your account is active. You can permanently delete your account and associated data from <strong>Profile → Security Settings → Delete Account</strong>.</p>
        <h2>6. Children's Privacy</h2>
        <p>Kuran AI is not directed at children under 13 and does not knowingly collect their personal information.</p>
        <h2>7. Changes and Contact</h2>
        <p>Material changes will be communicated through the app or email. For privacy questions, use our <a href="/contact">Contact page</a>.</p>
      </div>
    )
  }

  return (
    <div className={styles.page}>
      <h1>Gizlilik Politikası</h1>
      <p className={styles.updated}>Son güncelleme: 29 Temmuz 2026. Kuran AI uygulaması ve bağlantılı hizmetler için geçerlidir.</p>
      <h2>1. Topladığımız Veriler ve Amaçları</h2>
      <p>Kuran AI yalnızca hizmeti sunmak için gerekli olan asgari veriyi toplar:</p>
      <ul>
        <li><strong>Hesap bilgileri:</strong> E-posta, ad ve profil bilgileri kimlik doğrulama ve hesap yönetimi için kullanılır.</li>
        <li><strong>Profil fotoğrafı:</strong> İsteğe bağlıdır ve yalnızca profilinizde gösterilmek üzere saklanır.</li>
        <li><strong>Sohbet mesajları:</strong> Sorular ve yapay zekâ yanıtları sohbet geçmişini sunmak için saklanır. Mesaj içeriği yanıt üretimi amacıyla OpenAI API'sine gönderilir.</li>
        <li><strong>Sesli giriş:</strong> Ses kaydı, yazıya çevrilmek üzere sunucumuz üzerinden OpenAI'a gönderilir ve Kuran AI tarafından kalıcı olarak saklanmaz.</li>
        <li><strong>Konum:</strong> Yalnızca uygulama kullanılırken yerel namaz vakitlerini hesaplamak için kullanılır; sunucularımızda saklanmaz.</li>
        <li><strong>Kullanım verileri:</strong> Plan limitlerini yönetmek amacıyla günlük özellik kullanım sayıları tutulur.</li>
        <li><strong>Satın alma verileri:</strong> Abonelik Apple veya Google tarafından doğrulanır. Kart ya da ödeme yöntemi bilgilerinizi almayız.</li>
      </ul>
      <h2>2. Üçüncü Taraf Yapay Zekâ Hizmeti (OpenAI)</h2>
      <p>Açık izninizle mesaj metni ve ses kayıtları, yanıt ya da transkript üretmek üzere OpenAI, Inc. tarafından işlenir. Bu isteklere adınız veya e-posta adresiniz eklenmez. OpenAI verileri kendi <a href="https://openai.com/policies/privacy-policy" target="_blank" rel="noopener noreferrer">Gizlilik Politikası</a> ve <a href="https://openai.com/policies/terms-of-use" target="_blank" rel="noopener noreferrer">Kullanım Koşulları</a> kapsamında işler. İzin vermeyip yapay zekâ gerektirmeyen özellikleri kullanmaya devam edebilirsiniz.</p>
      <h2>3. Veri Paylaşımı</h2>
      <p><strong>Kişisel verilerinizi reklamverenlere veya veri aracılarına satmayız, kiralamayız ve paylaşmayız.</strong> Veriler yalnızca hizmeti çalıştırmak için gerekli sağlayıcılarla; yukarıda açıklandığı şekliyle OpenAI ve kimlik/abonelik doğrulaması için Apple veya Google dâhil olmak üzere paylaşılır.</p>
      <h2>4. Saklama ve Güvenlik</h2>
      <p>Veriler HTTPS/TLS üzerinden iletilir. Sunucu ve veri tabanı erişimleri sınırlandırılmış, sektör standardı kontrollerle korunmuştur.</p>
      <h2>5. Saklama Süresi ve Hesap Silme</h2>
      <p>Hesap verileri hesabınız etkin olduğu sürece tutulur. Hesabınızı ve ilişkili verileri <strong>Profil → Güvenlik Ayarları → Hesabı Sil</strong> yolundan kalıcı olarak silebilirsiniz.</p>
      <h2>6. Çocukların Gizliliği</h2>
      <p>Kuran AI, 13 yaş altındaki çocuklara yönelik değildir ve bilerek kişisel verilerini toplamaz.</p>
      <h2>7. Değişiklikler ve İletişim</h2>
      <p>Önemli değişiklikler uygulama veya e-posta üzerinden bildirilir. Gizlilik sorularınız için <a href="/contact">İletişim sayfamızı</a> kullanabilirsiniz.</p>
    </div>
  )
}

import { Link } from 'react-router-dom'
import { APP_STORE_URL, PLAY_STORE_URL } from '../config/app'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Landing.module.css'

export default function Landing() {
  const { text } = useLanguage()

  const capabilities = [
    {
      icon: 'ق',
      title: text('Kur’an odaklı sohbet', 'Quran-focused chat', 'محادثة تتمحور حول القرآن'),
      desc: text(
        'Ayet, kavram ve sureler hakkında sorun; cevabı anlaşılır başlıklar ve ilgili kaynaklarla okuyun.',
        'Ask about verses, concepts, and surahs; read clear, structured answers with relevant sources.',
        'اسأل عن الآيات والمفاهيم والسور، واقرأ إجابات واضحة ومنظمة مع المصادر ذات الصلة.',
      ),
      tone: 'gold',
    },
    {
      icon: 'ت',
      title: text('Meal, tefsir ve hadis', 'Translation, Tafsir, and Hadith', 'الترجمة والتفسير والحديث'),
      desc: text(
        'İhtiyacınıza göre kaynak türünü seçin ve aynı konuya farklı bir bakış açısından yaklaşın.',
        'Choose the source type you need and explore the same topic from a different perspective.',
        'اختر نوع المصدر الذي تحتاجه واستكشف الموضوع نفسه من منظور مختلف.',
      ),
      tone: 'violet',
    },
    {
      icon: '◉',
      title: text('Sesle sor, dinleyerek devam et', 'Ask by voice, continue by listening', 'اسأل بصوتك وتابع بالاستماع'),
      desc: text(
        'Mobil uygulamada doğal biçimde konuşun; otomatik durma algılama kaydı doğru anda tamamlasın.',
        'Speak naturally in the mobile app; automatic silence detection completes recording at the right moment.',
        'تحدث بطبيعتك في تطبيق الهاتف، ويُنهي اكتشاف الصمت التلقائي التسجيل في الوقت المناسب.',
      ),
      tone: 'blue',
    },
    {
      icon: '⌘',
      title: text('Her cihazda tek hesap', 'One account across every device', 'حساب واحد على جميع أجهزتك'),
      desc: text(
        'Web, iOS ve Android’de aynı hesapla oturum açın; sohbetlerinize kaldığınız yerden devam edin.',
        'Sign in with the same account on web, iOS, and Android, then continue where you left off.',
        'سجّل الدخول بالحساب نفسه على الويب وiOS وAndroid وتابع من حيث توقفت.',
      ),
      tone: 'green',
    },
  ]

  const steps = [
    {
      no: '01',
      title: text('Sorunuzu sorun', 'Ask your question', 'اطرح سؤالك'),
      desc: text('Yazarak veya mobil uygulamada sesinizi kullanarak başlayın.', 'Start by typing or using your voice in the mobile app.', 'ابدأ بالكتابة أو باستخدام صوتك في تطبيق الهاتف.'),
    },
    {
      no: '02',
      title: text('Kaynağı belirleyin', 'Choose the source', 'اختر المصدر'),
      desc: text('Kur’an meali, tefsir veya hadis görünümünü seçin.', 'Select Quran translation, Tafsir, or Hadith.', 'اختر ترجمة القرآن أو التفسير أو الحديث.'),
    },
    {
      no: '03',
      title: text('Sade bir cevap okuyun', 'Read a clear answer', 'اقرأ إجابة واضحة'),
      desc: text('Başlıkları, vurguları ve kaynak notlarını kolayca takip edin.', 'Follow headings, key points, and source notes with ease.', 'تابع العناوين والنقاط المهمة وملاحظات المصادر بسهولة.'),
    },
  ]

  return (
    <div className={styles.landing}>
      <section className={styles.hero}>
        <div className={styles.heroGrid} aria-hidden="true" />
        <div className={styles.heroGlow} aria-hidden="true" />
        <div className={styles.heroInner}>
          <div className={styles.heroCopy}>
            <div className={styles.productBadge}>
              <span className={styles.liveDot} />
              {text('Web, iOS ve Android’de hazır', 'Ready on web, iOS, and Android', 'متاح على الويب وiOS وAndroid')}
            </div>
            <p className={styles.eyebrow}>{text('ANLAMAYA ODAKLANAN DİJİTAL REHBER', 'A DIGITAL GUIDE BUILT FOR UNDERSTANDING', 'دليل رقمي صُمم للفهم')}</p>
            <h1>
              {text('Kur’an’ı keşfetmenin ', 'A calmer way to ', 'طريقة أكثر هدوءًا ')}
              <span>{text('daha sade yolu.', 'explore the Quran.', 'لاستكشاف القرآن.')}</span>
            </h1>
            <p className={styles.heroLead}>
              {text(
                'Kur’an meali, tefsir ve hadis kaynaklarıyla çalışan yapay zekâ destekli asistan. Sorunuzu sorun, kaynağınızı seçin, anlaşılır bir cevapla ilerleyin.',
                'An AI-assisted companion grounded in Quran translations, Tafsir, and Hadith. Ask a question, choose your source, and move forward with a clear answer.',
                'مساعد مدعوم بالذكاء الاصطناعي يعتمد على ترجمات القرآن والتفسير والحديث. اطرح سؤالك واختر مصدرك وتابع بإجابة واضحة.',
              )}
            </p>
            <div className={styles.heroActions}>
              <Link to="/app" className={styles.primaryButton}>
                {text('Ücretsiz kullanmaya başla', 'Start using it free', 'ابدأ الاستخدام مجانًا')}
                <span aria-hidden="true">→</span>
              </Link>
              <a href="#how-it-works" className={styles.secondaryButton}>
                {text('Nasıl çalışır?', 'How does it work?', 'كيف يعمل؟')}
              </a>
            </div>
            <div className={styles.heroFootnote}>
              <span>✓ {text('Kredi kartı gerekmez', 'No credit card required', 'لا حاجة إلى بطاقة ائتمان')}</span>
              <span>✓ {text('Günde 3 ücretsiz soru', '3 free questions daily', '3 أسئلة مجانية يوميًا')}</span>
            </div>
          </div>

          <div className={styles.productStage}>
            <div className={styles.stageAura} aria-hidden="true" />
            <div className={styles.appWindow}>
              <div className={styles.windowTop}>
                <div className={styles.windowBrand}><span>ق</span><div><strong>Kuran AI</strong><small>{text('Kaynak odaklı asistan', 'Source-focused assistant', 'مساعد يركز على المصادر')}</small></div></div>
                <div className={styles.windowStatus}><i /> {text('Hazır', 'Ready', 'جاهز')}</div>
              </div>
              <div className={styles.sourceTabs}>
                <span className={styles.sourceActive}>{text('Kur’an Meali', 'Quran', 'ترجمة القرآن')}</span>
                <span>{text('Tefsir', 'Tafsir', 'التفسير')}</span>
                <span>{text('Hadis', 'Hadith', 'الحديث')}</span>
              </div>
              <div className={styles.questionBubble}>
                {text('Sabır hakkında Kur’an ne söyler?', 'What does the Quran say about patience?', 'ماذا يقول القرآن عن الصبر؟')}
              </div>
              <div className={styles.answerCard}>
                <div className={styles.answerIdentity}><span>ق</span><strong>Kuran AI</strong><small>{text('Şimdi', 'Now', 'الآن')}</small></div>
                <h3>{text('Sabır, bilinçli bir dayanıklılıktır', 'Patience is conscious resilience', 'الصبر ثبات واعٍ')}</h3>
                <p>
                  {text(
                    'Kur’an’da sabır; yalnızca beklemek değil, zorluk karşısında doğru tavrı korumak ve güveni sürdürmek olarak ele alınır.',
                    'In the Quran, patience is more than waiting: it is preserving the right response and sustaining trust through difficulty.',
                    'الصبر في القرآن ليس مجرد انتظار، بل هو الحفاظ على الموقف الصحيح واستمرار الثقة عند الشدائد.',
                  )}
                </p>
                <div className={styles.citations}>
                  <span>Bakara 2:153</span>
                  <span>Asr 103:3</span>
                  <span>{text('Kaynakları gör', 'View sources', 'عرض المصادر')} ↗</span>
                </div>
              </div>
              <div className={styles.promptBar}><span>{text('Ayet, konu veya kavram hakkında sorun…', 'Ask about a verse, topic, or concept…', 'اسأل عن آية أو موضوع أو مفهوم…')}</span><b>↑</b></div>
            </div>
            <div className={`${styles.floatingCard} ${styles.verseCard}`}>
              <span>{text('GÜNÜN AYETİ', 'VERSE OF THE DAY', 'آية اليوم')}</span>
              <strong>“Şüphesiz zorlukla beraber bir kolaylık vardır.”</strong>
              <small>İnşirah 94:5</small>
            </div>
            <div className={`${styles.floatingCard} ${styles.syncCard}`}>
              <i>✓</i>
              <div><strong>{text('Sohbet eşitlendi', 'Conversation synced', 'تمت مزامنة المحادثة')}</strong><small>Web · iOS · Android</small></div>
            </div>
          </div>
        </div>
      </section>

      <section className={styles.platformBar} aria-label={text('Desteklenen platformlar', 'Supported platforms', 'المنصات المدعومة')}>
        <div>
          <p>{text('Tek hesap. Her yerde aynı deneyim.', 'One account. The same experience everywhere.', 'حساب واحد. التجربة نفسها في كل مكان.')}</p>
          <div className={styles.platformLogos}>
            <span>● App Store</span><i />
            <span>▶ Google Play</span><i />
            <span>◉ {text('Web uygulaması', 'Web app', 'تطبيق الويب')}</span>
          </div>
        </div>
      </section>

      <section className={styles.capabilities}>
        <div className={styles.sectionIntro}>
          <p className={styles.eyebrow}>{text('BİLGİYE GİDEN SADE YOL', 'A CLEAR PATH TO KNOWLEDGE', 'طريق واضح إلى المعرفة')}</p>
          <h2>{text('Merak ettiğiniz konuya odaklanın.', 'Stay focused on what you want to understand.', 'ركّز على ما تريد فهمه.')}</h2>
          <p>{text('Karmaşık menüler yerine, ihtiyacınız olan kaynağı ve sohbeti merkeze alan bir deneyim.', 'A focused experience that puts your source and conversation first, without complicated menus.', 'تجربة مركزة تضع المصدر والمحادثة أولًا، من دون قوائم معقدة.')}</p>
        </div>
        <div className={styles.capabilityGrid}>
          {capabilities.map(item => (
            <article key={item.title} className={`${styles.capabilityCard} ${styles[item.tone]}`}>
              <span className={styles.capabilityIcon}>{item.icon}</span>
              <h3>{item.title}</h3>
              <p>{item.desc}</p>
              <div className={styles.cardLine} />
            </article>
          ))}
        </div>
      </section>

      <section id="how-it-works" className={styles.workflow}>
        <div className={styles.workflowInner}>
          <div className={styles.workflowIntro}>
            <p className={styles.eyebrow}>{text('NASIL ÇALIŞIR?', 'HOW IT WORKS', 'كيف يعمل؟')}</p>
            <h2>{text('Sorudan kaynağa, üç sade adım.', 'From question to source in three clear steps.', 'من السؤال إلى المصدر في ثلاث خطوات واضحة.')}</h2>
            <p>{text('Kuran AI, araştırma akışınızı bölmeden doğru bağlama ulaşmanızı kolaylaştırır.', 'Kuran AI helps you reach the right context without interrupting your train of thought.', 'يساعدك Kuran AI على الوصول إلى السياق الصحيح دون قطع تسلسل أفكارك.')}</p>
            <Link to="/app" className={styles.textLink}>{text('Web uygulamasını aç', 'Open the web app', 'افتح تطبيق الويب')} <span>→</span></Link>
          </div>
          <div className={styles.steps}>
            {steps.map(step => (
              <article key={step.no}>
                <span>{step.no}</span>
                <div><h3>{step.title}</h3><p>{step.desc}</p></div>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className={styles.crossPlatform}>
        <div className={styles.crossPlatformInner}>
          <div className={styles.phoneFrame}>
            <div className={styles.phoneTop}><i /><span>Kuran AI</span><b>•••</b></div>
            <div className={styles.mobileVerse}><span>۞</span><small>{text('Günün Ayeti', 'Verse of the Day', 'آية اليوم')}</small><p>{text('“Kalpler ancak Allah’ı anmakla huzur bulur.”', '“Hearts find comfort in the remembrance of Allah.”', '«ألا بذكر الله تطمئن القلوب»')}</p><b>Ra’d 13:28</b></div>
            <div className={styles.mobileTopics}><span>{text('Kur’an', 'Quran', 'القرآن')}</span><span>{text('Tefsir', 'Tafsir', 'التفسير')}</span><span>{text('Hadis', 'Hadith', 'الحديث')}</span></div>
            <div className={styles.mobilePrompt}><span>{text('Sorunuzu yazın veya söyleyin…', 'Type or speak your question…', 'اكتب سؤالك أو انطقه…')}</span><b>◉</b></div>
          </div>
          <div className={styles.crossCopy}>
            <p className={styles.eyebrow}>{text('CEBİNİZDE, MASANIZDA, HER YERDE', 'IN YOUR POCKET, AT YOUR DESK, EVERYWHERE', 'في جيبك وعلى مكتبك وفي كل مكان')}</p>
            <h2>{text('Sohbetiniz cihaz değiştirince yarım kalmaz.', 'Your conversation continues when your device changes.', 'تستمر محادثتك عند الانتقال بين الأجهزة.')}</h2>
            <p>{text('Google, Apple veya geçiş anahtarıyla güvenli biçimde oturum açın. Aynı hesabı web, iOS ve Android’de kullanın.', 'Sign in securely with Google, Apple, or a passkey. Use the same account on web, iOS, and Android.', 'سجّل الدخول بأمان باستخدام Google أو Apple أو مفتاح المرور، واستخدم الحساب نفسه على الويب وiOS وAndroid.')}</p>
            <ul>
              <li><span>✓</span>{text('Sesli soru ve otomatik durma algılama', 'Voice questions with automatic silence detection', 'أسئلة صوتية مع اكتشاف الصمت تلقائيًا')}</li>
              <li><span>✓</span>{text('Sohbet geçmişine her cihazdan erişim', 'Conversation history on every device', 'سجل المحادثات على كل جهاز')}</li>
              <li><span>✓</span>{text('Apple ve Google üzerinden güvenli abonelik', 'Secure subscriptions through Apple and Google', 'اشتراكات آمنة عبر Apple وGoogle')}</li>
            </ul>
            <div className={styles.storeButtons}>
              <a href={APP_STORE_URL} target="_blank" rel="noopener noreferrer"><small>{text('Şuradan indirin', 'Download on the', 'نزّل من')}</small><strong>App Store</strong></a>
              <a href={PLAY_STORE_URL} target="_blank" rel="noopener noreferrer"><small>{text('Şuradan edinin', 'Get it on', 'احصل عليه من')}</small><strong>Google Play</strong></a>
            </div>
          </div>
        </div>
      </section>

      <section className={styles.trust}>
        <div className={styles.trustInner}>
          <div className={styles.trustMark}>◇</div>
          <div className={styles.trustCopy}>
            <p className={styles.eyebrow}>{text('GÜVEN VE ŞEFFAFLIK', 'TRUST AND TRANSPARENCY', 'الثقة والشفافية')}</p>
            <h2>{text('Teknoloji yardımcı olur; kaynaklar yön gösterir.', 'Technology assists; sources provide direction.', 'التقنية تساعد، والمصادر ترشد.')}</h2>
            <p>{text('Kuran AI cevapları öğrenme ve araştırma amacıyla sunar. Dinî hüküm ve kişisel fetva gerektiren konularda ehil uzmanlara danışmanızı önerir.', 'Kuran AI provides answers for learning and research. For religious rulings and personal fatwas, we recommend consulting a qualified scholar.', 'يقدم Kuran AI الإجابات للتعلم والبحث. وفي مسائل الأحكام الشرعية والفتاوى الشخصية نوصي بالرجوع إلى عالم مؤهل.')}</p>
          </div>
          <div className={styles.trustLinks}>
            <Link to="/privacy-policy">{text('Gizlilik Politikası', 'Privacy Policy', 'سياسة الخصوصية')} <span>↗</span></Link>
            <Link to="/terms">{text('Kullanım Koşulları', 'Terms of Use', 'شروط الاستخدام')} <span>↗</span></Link>
            <Link to="/support">{text('Destek Merkezi', 'Support Center', 'مركز الدعم')} <span>↗</span></Link>
          </div>
        </div>
      </section>

      <section className={styles.finalCta}>
        <div className={styles.finalGlow} aria-hidden="true" />
        <div>
          <span className={styles.finalMark}>ق</span>
          <p className={styles.eyebrow}>{text('MERAKINIZLA BAŞLAYIN', 'START WITH YOUR CURIOSITY', 'ابدأ بفضولك')}</p>
          <h2>{text('Bugün neyi keşfetmek istersiniz?', 'What would you like to explore today?', 'ماذا تود أن تستكشف اليوم؟')}</h2>
          <p>{text('Hemen webden başlayın. İlk üç sorunuz ücretsiz.', 'Start on the web now. Your first three questions are free.', 'ابدأ الآن عبر الويب. أسئلتك الثلاثة الأولى مجانية.')}</p>
          <Link to="/app" className={styles.primaryButton}>{text('Kuran AI’ı aç', 'Open Kuran AI', 'افتح Kuran AI')} <span>→</span></Link>
        </div>
      </section>
    </div>
  )
}

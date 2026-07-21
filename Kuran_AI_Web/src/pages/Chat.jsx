import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { API_BASE, TOKEN_KEY } from '../config/app'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Chat.module.css'

const CONSENT_KEY = 'kuranai_ai_consent'

function readTokens() {
  try { return JSON.parse(localStorage.getItem(TOKEN_KEY)) }
  catch { return null }
}

async function parseResponse(response) {
  const data = await response.json().catch(() => null)
  if (!response.ok || !data?.isSuccess) {
    throw new Error(data?.message || `HTTP ${response.status}`)
  }
  return data.data
}

export default function Chat() {
  const { text } = useLanguage()
  const [tokens, setTokens] = useState(readTokens)
  const [user, setUser] = useState(null)
  const [authMode, setAuthMode] = useState('login')
  const [authForm, setAuthForm] = useState({ name: '', username: '', email: '', credential: '', password: '' })
  const [authLoading, setAuthLoading] = useState(false)
  const [authError, setAuthError] = useState('')
  const [sourceType, setSourceType] = useState('meal')
  const [messages, setMessages] = useState([])
  const [question, setQuestion] = useState('')
  const [sessionId, setSessionId] = useState(null)
  const [streaming, setStreaming] = useState(false)
  const [chatError, setChatError] = useState('')
  const [subscriptionOverview, setSubscriptionOverview] = useState(null)
  const [limitInfo, setLimitInfo] = useState(null)
  const [showConsent, setShowConsent] = useState(false)
  const [pendingQuestion, setPendingQuestion] = useState('')
  const messagesEndRef = useRef(null)

  const storeTokens = (nextTokens) => {
    setTokens(nextTokens)
    if (nextTokens) localStorage.setItem(TOKEN_KEY, JSON.stringify(nextTokens))
    else localStorage.removeItem(TOKEN_KEY)
  }

  const refreshAccessToken = async (currentTokens) => {
    if (!currentTokens?.refreshToken) return null
    try {
      const response = await fetch(`${API_BASE}/api/user/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken: currentTokens.refreshToken }),
      })
      const nextTokens = await parseResponse(response)
      storeTokens(nextTokens)
      return nextTokens
    } catch {
      storeTokens(null)
      setUser(null)
      return null
    }
  }

  const authorizedFetch = async (path, options = {}) => {
    let activeTokens = tokens
    if (!activeTokens?.accessToken) throw new Error(text('Oturum açmanız gerekiyor.', 'You need to sign in.'))
    const run = token => fetch(`${API_BASE}${path}`, {
      ...options,
      headers: { ...options.headers, Authorization: `Bearer ${token}` },
    })
    let response = await run(activeTokens.accessToken)
    if (response.status === 401) {
      activeTokens = await refreshAccessToken(activeTokens)
      if (!activeTokens) throw new Error(text('Oturumunuz sona erdi. Lütfen tekrar giriş yapın.', 'Your session expired. Please sign in again.'))
      response = await run(activeTokens.accessToken)
    }
    return response
  }

  useEffect(() => {
    if (!tokens?.accessToken) return
    let active = true
    authorizedFetch('/api/user/auth/me')
      .then(parseResponse)
      .then(data => { if (active) setUser(data) })
      .catch(() => { if (active) storeTokens(null) })
    authorizedFetch('/api/v1/subscription/me')
      .then(parseResponse)
      .then(data => { if (active) setSubscriptionOverview(data) })
      .catch(() => { /* Chat remains usable if usage details are temporarily unavailable. */ })
    return () => { active = false }
    // Existing tokens are validated once when their access token changes.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tokens?.accessToken])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, streaming])

  const login = async (credential, password) => {
    const response = await fetch(`${API_BASE}/api/user/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ credential, password }),
    })
    const nextTokens = await parseResponse(response)
    storeTokens(nextTokens)
  }

  const handleAuth = async (event) => {
    event.preventDefault()
    setAuthLoading(true)
    setAuthError('')
    try {
      if (authMode === 'register') {
        const response = await fetch(`${API_BASE}/api/user/auth/register`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            name: authForm.name.trim(),
            username: authForm.username.trim(),
            email: authForm.email.trim(),
            password: authForm.password,
          }),
        })
        await parseResponse(response)
        await login(authForm.email.trim(), authForm.password)
      } else {
        await login(authForm.credential.trim(), authForm.password)
      }
    } catch (error) {
      setAuthError(error.message)
    } finally {
      setAuthLoading(false)
    }
  }

  const logout = async () => {
    try { await authorizedFetch('/api/user/auth/logout', { method: 'POST' }) }
    catch { /* Local session is cleared even if the network request fails. */ }
    storeTokens(null)
    setUser(null)
    setMessages([])
    setSessionId(null)
    setSubscriptionOverview(null)
    setLimitInfo(null)
  }

  const performSend = async (content) => {
    const trimmed = content.trim()
    if (!trimmed || streaming) return
    setQuestion('')
    setChatError('')
    setLimitInfo(null)
    setStreaming(true)
    setMessages(previous => [...previous, { role: 'user', content: trimmed }, { role: 'assistant', content: '' }])

    try {
      const response = await authorizedFetch('/api/v1/chat/message', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sessionId, sourceType, content: trimmed }),
      })

      if (!response.ok) {
        const errorBody = await response.json().catch(() => null)
        if (response.status === 402 && errorBody?.data?.code === 'MESSAGE_LIMIT_EXCEEDED') {
          setLimitInfo(errorBody.data)
        }
        throw new Error(errorBody?.message || text('Mesaj gönderilemedi.', 'Could not send the message.'))
      }

      const reader = response.body?.getReader()
      if (!reader) throw new Error(text('Yanıt akışı başlatılamadı.', 'The response stream could not start.'))

      const decoder = new TextDecoder()
      let buffer = ''
      while (true) {
        const { value, done } = await reader.read()
        buffer += decoder.decode(value || new Uint8Array(), { stream: !done })
        const events = buffer.split('\n\n')
        buffer = events.pop() || ''

        for (const event of events) {
          const dataLine = event.split('\n').find(line => line.startsWith('data:'))
          if (!dataLine) continue
          const payload = JSON.parse(dataLine.replace(/^data:\s*/, ''))
          if (payload.type === 'started' || payload.type === 'done') setSessionId(payload.sessionId)
          if (payload.type === 'done') {
            setSubscriptionOverview(previous => previous ? {
              ...previous,
              usage: { ...previous.usage, messageUsed: (previous.usage?.messageUsed || 0) + 1 },
            } : previous)
          }
          if (payload.type === 'error') throw new Error(payload.message)
          if (payload.type === 'chunk') {
            setMessages(previous => previous.map((message, index) => (
              index === previous.length - 1 ? { ...message, content: message.content + payload.content } : message
            )))
          }
        }
        if (done) break
      }
    } catch (error) {
      setChatError(error.message)
      setMessages(previous => previous.filter((message, index) => index !== previous.length - 1 || message.content))
    } finally {
      setStreaming(false)
    }
  }

  const handleSend = (event) => {
    event?.preventDefault()
    if (!question.trim()) return
    if (localStorage.getItem(CONSENT_KEY) !== 'granted') {
      setPendingQuestion(question)
      setShowConsent(true)
      return
    }
    performSend(question)
  }

  const allowAI = () => {
    localStorage.setItem(CONSENT_KEY, 'granted')
    setShowConsent(false)
    const content = pendingQuestion
    setPendingQuestion('')
    performSend(content)
  }

  const suggestions = sourceType === 'hadis'
    ? [text('Niyet ile ilgili sahih hadisleri açıklar mısın?', 'Can you explain authentic hadiths about intention?'), text('Komşuluk hakkı hakkında hadisler nelerdir?', 'What do hadiths say about neighbors’ rights?')]
    : sourceType === 'tefsir'
      ? [text('Fâtiha suresinin ana mesajı nedir?', 'What is the main message of Surah Al-Fatiha?'), text('Bakara 286. ayetin tefsirini özetler misin?', 'Can you summarize the Tafsir of Al-Baqarah 2:286?')]
      : [text('Sabır hakkında Kur’an ne söyler?', 'What does the Quran say about patience?'), text('Rahman suresinden bir ayeti açıklar mısın?', 'Can you explain a verse from Surah Ar-Rahman?')]

  const usage = subscriptionOverview?.usage
  const currentPackage = subscriptionOverview?.subscription?.package
  const formatLimit = value => value == null || value >= 999999 ? '∞' : value

  if (!tokens) {
    return (
      <section className={styles.authShell}>
        <div className={styles.authIntro}>
          <span className={styles.miniMark}>ق</span>
          <p className={styles.kicker}>{text('Web sohbeti', 'Web chat')}</p>
          <h1>{text('Sorunuzla başlayın.', 'Begin with your question.')}</h1>
          <p>{text('Kur’an meali, tefsir ve hadis kaynakları üzerine yapay zekâ destekli sohbet. iOS hesabınızla giriş yapabilir veya ücretsiz hesap oluşturabilirsiniz.', 'AI-assisted chat about Quran translations, Tafsir, and Hadith. Sign in with your iOS account or create a free account.')}</p>
          <ul>
            <li>{text('Ücretsiz günde 3 soru', '3 free questions per day')}</li>
            <li>{text('Sohbetler cihazlarınız arasında ortak', 'Chats shared across your devices')}</li>
            <li>{text('Kaynak türünü siz seçersiniz', 'You choose the source type')}</li>
          </ul>
        </div>

        <form className={styles.authCard} onSubmit={handleAuth}>
          <div className={styles.authTabs}>
            <button type="button" className={authMode === 'login' ? styles.selectedTab : ''} onClick={() => { setAuthMode('login'); setAuthError('') }}>{text('Giriş yap', 'Sign in')}</button>
            <button type="button" className={authMode === 'register' ? styles.selectedTab : ''} onClick={() => { setAuthMode('register'); setAuthError('') }}>{text('Hesap oluştur', 'Create account')}</button>
          </div>

          {authMode === 'register' && <>
            <label>{text('Adınız', 'Your name')}<input required value={authForm.name} onChange={e => setAuthForm({ ...authForm, name: e.target.value })} autoComplete="name" /></label>
            <label>{text('Kullanıcı adı', 'Username')}<input required minLength={2} value={authForm.username} onChange={e => setAuthForm({ ...authForm, username: e.target.value })} autoComplete="username" /></label>
            <label>{text('E-posta', 'Email')}<input required type="email" value={authForm.email} onChange={e => setAuthForm({ ...authForm, email: e.target.value })} autoComplete="email" /></label>
          </>}
          {authMode === 'login' && <label>{text('E-posta veya kullanıcı adı', 'Email or username')}<input required value={authForm.credential} onChange={e => setAuthForm({ ...authForm, credential: e.target.value })} autoComplete="username" /></label>}
          <label>{text('Şifre', 'Password')}<input required type="password" minLength={6} value={authForm.password} onChange={e => setAuthForm({ ...authForm, password: e.target.value })} autoComplete={authMode === 'login' ? 'current-password' : 'new-password'} /></label>
          {authError && <p className={styles.error}>{authError}</p>}
          <button className={styles.authSubmit} disabled={authLoading}>{authLoading ? text('Bağlanıyor…', 'Connecting…') : authMode === 'login' ? text('Sohbete gir', 'Enter chat') : text('Hesabımı oluştur', 'Create my account')}</button>
          <p className={styles.legal}>{text('Devam ederek ', 'By continuing, you accept the ')}<Link to="/terms">{text('Kullanım Koşullarını', 'Terms of Use')}</Link>{text(' ve ', ' and ') }<Link to="/privacy-policy">{text('Gizlilik Politikasını', 'Privacy Policy')}</Link>{text(' kabul edersiniz.', '.')}</p>
        </form>
      </section>
    )
  }

  return (
    <section className={styles.chatShell}>
      <aside className={styles.sidebar}>
        <div>
          <p className={styles.kicker}>{text('Kaynak', 'Source')}</p>
          <div className={styles.sourceList}>
            {[['meal', text('Kur’an Meali', 'Quran Translation')], ['tefsir', text('Tefsir', 'Tafsir')], ['hadis', text('Hadis', 'Hadith')]].map(([value, label]) => (
              <button key={value} className={sourceType === value ? styles.sourceActive : ''} onClick={() => { setSourceType(value); setSessionId(null); setMessages([]) }}>
                <span>{value === 'meal' ? 'ق' : value === 'tefsir' ? 'ت' : 'ح'}</span>{label}
              </button>
            ))}
          </div>
        </div>
        <Link to="/packages" className={styles.planBlock}>
          <span>{currentPackage ? '✦' : '＋'}</span>
          <div>
            <strong>{currentPackage ? text(currentPackage.name, currentPackage.nameEn) : text('Ücretsiz paket', 'Free plan')}</strong>
            <small>{usage ? `${usage.messageUsed} / ${formatLimit(usage.messageLimit)} ${text('soru bugün', 'questions today')}` : text('Paketleri görüntüle', 'View plans')}</small>
          </div>
          <b>›</b>
        </Link>
        <div className={styles.accountBlock}>
          <span className={styles.avatar}>{(user?.name || user?.username || 'K').charAt(0).toUpperCase()}</span>
          <div><strong>{user?.name || user?.username || text('Kullanıcı', 'User')}</strong><small>{user?.email || ''}</small></div>
          <button onClick={logout} title={text('Çıkış yap', 'Sign out')}>↗</button>
        </div>
      </aside>

      <div className={styles.conversation}>
        <header className={styles.chatHeader}>
          <div><span>{sourceType === 'meal' ? 'ق' : sourceType === 'tefsir' ? 'ت' : 'ح'}</span><div><strong>{sourceType === 'meal' ? text('Kur’an Meali', 'Quran Translation') : sourceType === 'tefsir' ? text('Tefsir', 'Tafsir') : text('Hadis', 'Hadith')}</strong><small>{text('Kaynak destekli yapay zekâ asistanı', 'Source-assisted AI')}</small></div></div>
          <button onClick={() => { setMessages([]); setSessionId(null); setChatError('') }}>＋ {text('Yeni sohbet', 'New chat')}</button>
        </header>

        <div className={styles.messages} aria-live="polite">
          {messages.length === 0 ? (
            <div className={styles.emptyState}>
              <span className={styles.largeMark}>ق</span>
              <p className={styles.kicker}>{text('Vahiy ile konuş.', 'Speak with revelation.')}</p>
              <h1>{text('Bugün neyi keşfetmek istersiniz?', 'What would you like to explore today?')}</h1>
              <div className={styles.suggestions}>{suggestions.map(item => <button key={item} onClick={() => setQuestion(item)}>{item}<span>↗</span></button>)}</div>
            </div>
          ) : messages.map((message, index) => (
            <article key={`${message.role}-${index}`} className={message.role === 'user' ? styles.userMessage : styles.assistantMessage}>
              <span className={styles.messageRole}>{message.role === 'user' ? text('Siz', 'You') : 'Kuran AI'}</span>
              <p>{message.content || (streaming && index === messages.length - 1 ? text('Düşünüyor…', 'Thinking…') : '')}</p>
            </article>
          ))}
          {chatError && <p className={styles.error}>{chatError}</p>}
          {limitInfo && (
            <div className={styles.limitBanner}>
              <div><strong>{text('Günlük soru limitiniz doldu.', 'You have reached today’s question limit.')}</strong><p>{text('Sohbete devam etmek için paketinizi Kuran AI iOS uygulamasından seçebilirsiniz.', 'To continue chatting, choose a plan in the Kuran AI iOS app.')}</p></div>
              <Link to="/packages">{text('iOS paketlerini aç', 'Open iOS plans')} →</Link>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>

        <form className={styles.composer} onSubmit={handleSend}>
          <textarea
            value={question}
            onChange={event => setQuestion(event.target.value)}
            onKeyDown={event => {
              if (event.key === 'Enter' && !event.shiftKey) { event.preventDefault(); handleSend() }
            }}
            placeholder={text('Ayet, konu veya kavram hakkında sorun…', 'Ask about a verse, topic, or concept…')}
            rows={1}
            maxLength={2000}
          />
          <button disabled={streaming || !question.trim()} aria-label={text('Gönder', 'Send')}>{streaming ? '···' : '↑'}</button>
          <small>{text('Yanıtlar eğitim amaçlıdır; fetva niteliği taşımaz.', 'Responses are educational and are not religious rulings.')}</small>
        </form>
      </div>

      {showConsent && (
        <div className={styles.modalBackdrop} role="presentation">
          <div className={styles.consentModal} role="dialog" aria-modal="true" aria-labelledby="consent-title">
            <span className={styles.miniMark}>AI</span>
            <h2 id="consent-title">{text('OpenAI ile veri paylaşımı', 'Data sharing with OpenAI')}</h2>
            <p>{text('Yazdığınız mesaj metni, yanıt üretmek için OpenAI, Inc. ile paylaşılır. Adınız, e-posta adresiniz ve hesap bilgileriniz gönderilmez. İzin vermezseniz mesajınız paylaşılmaz ve sohbet özelliği çalışmaz.', 'Your message text is shared with OpenAI, Inc. to generate a response. Your name, email, and account details are not sent. If you decline, your message is not shared and chat will remain unavailable.')}</p>
            <Link to="/privacy-policy" target="_blank">{text('Gizlilik politikasını incele', 'Read the privacy policy')}</Link>
            <div><button className={styles.secondaryButton} onClick={() => { setShowConsent(false); setPendingQuestion('') }}>{text('Şimdi değil', 'Not now')}</button><button className={styles.authSubmit} onClick={allowAI}>{text('İzin ver ve devam et', 'Allow and continue')}</button></div>
          </div>
        </div>
      )}
    </section>
  )
}

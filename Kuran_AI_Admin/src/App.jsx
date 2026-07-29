import { createContext, createElement, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import {
  Activity, ArrowLeft, ChevronLeft, ChevronRight, CircleDollarSign, CreditCard,
  ExternalLink, Eye, LayoutDashboard, LogOut, Menu, MessageSquareText, MonitorSmartphone,
  RefreshCw, Search, ShieldCheck, Smartphone, Users, X,
} from 'lucide-react'
import { request, TOKEN_KEY } from './api'

const AuthContext = createContext(null)
const RouterContext = createContext(null)

const navigation = [
  { path: '/', label: 'Genel Bakış', icon: LayoutDashboard },
  { path: '/users', label: 'Kullanıcılar', icon: Users },
  { path: '/purchases', label: 'Satın Almalar', icon: CreditCard },
  { path: '/questions', label: 'Sorular', icon: MessageSquareText },
]

function AuthProvider({ children }) {
  const [session, setSession] = useState(() => {
    try { return JSON.parse(sessionStorage.getItem(TOKEN_KEY)) }
    catch { return null }
  })
  const [manager, setManager] = useState(null)
  const [checking, setChecking] = useState(Boolean(session?.accessToken))

  const logout = useCallback(() => {
    sessionStorage.removeItem(TOKEN_KEY)
    setSession(null)
    setManager(null)
  }, [])

  useEffect(() => {
    if (!session?.accessToken) return
    let active = true
    request('/api/manager/auth/me', {}, session.accessToken)
      .then(data => { if (active) setManager(data) })
      .catch(() => { if (active) logout() })
      .finally(() => { if (active) setChecking(false) })
    return () => { active = false }
  }, [session?.accessToken, logout])

  const login = async (email, password) => {
    const next = await request('/api/manager/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    sessionStorage.setItem(TOKEN_KEY, JSON.stringify(next))
    setChecking(true)
    setSession(next)
  }

  const value = useMemo(() => ({
    token: session?.accessToken || '',
    authenticated: Boolean(session?.accessToken),
    manager,
    checking,
    login,
    logout,
  }), [session?.accessToken, manager, checking, logout])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

function useAuth() {
  return useContext(AuthContext)
}

function App() {
  return (
    <RouterProvider>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </RouterProvider>
  )
}

function AppRoutes() {
  const auth = useAuth()
  const pathname = useLocation().pathname
  if (auth.checking) return <FullLoader label="Oturum doğrulanıyor" />
  if (!auth.authenticated) return <LoginPage />
  let page = <OverviewPage />
  if (pathname === '/users') page = <UsersPage />
  else if (/^\/users\/\d+$/.test(pathname)) page = <UserDetailPage />
  else if (pathname === '/purchases') page = <PurchasesPage />
  else if (pathname === '/questions') page = <QuestionsPage />
  return (
    <DashboardShell>
      {page}
    </DashboardShell>
  )
}

function RouterProvider({ children }) {
  const [pathname, setPathname] = useState(() => window.location.pathname)
  useEffect(() => {
    const handlePopState = () => setPathname(window.location.pathname)
    window.addEventListener('popstate', handlePopState)
    return () => window.removeEventListener('popstate', handlePopState)
  }, [])
  const navigate = useCallback((path, options = {}) => {
    window.history[options.replace ? 'replaceState' : 'pushState']({}, '', path)
    setPathname(path)
    window.scrollTo({ top: 0, behavior: 'instant' })
  }, [])
  const value = useMemo(() => ({ pathname, navigate }), [pathname, navigate])
  return <RouterContext.Provider value={value}>{children}</RouterContext.Provider>
}

function useLocation() {
  const { pathname } = useContext(RouterContext)
  return { pathname }
}

function useNavigate() {
  return useContext(RouterContext).navigate
}

function useParams() {
  const parts = useLocation().pathname.split('/').filter(Boolean)
  return { id: parts[1] || '' }
}

function Link({ to, children, onClick, ...props }) {
  const navigate = useNavigate()
  const click = event => {
    onClick?.(event)
    if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return
    event.preventDefault()
    navigate(to)
  }
  return <a href={to} onClick={click} {...props}>{children}</a>
}

function LoginPage() {
  const { login } = useAuth()
  const [form, setForm] = useState({ email: '', password: '' })
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  const submit = async event => {
    event.preventDefault()
    setBusy(true)
    setError('')
    try { await login(form.email.trim(), form.password) }
    catch (err) { setError(err.message === 'Invalid credentials' ? 'E-posta veya parola hatalı.' : err.message) }
    finally { setBusy(false) }
  }

  return (
    <main className="login-page">
      <div className="login-aurora" />
      <section className="login-card">
        <div className="brand-mark large">ق</div>
        <div className="login-heading">
          <span className="eyebrow">KURAN AI OPERATIONS</span>
          <h1>Yönetim paneli</h1>
          <p>Kullanıcı, sohbet ve abonelik operasyonlarına güvenli erişim.</p>
        </div>
        <form onSubmit={submit}>
          <label>
            E-posta
            <input type="email" autoComplete="username" required value={form.email}
              onChange={event => setForm({ ...form, email: event.target.value })} placeholder="manager@company.com" />
          </label>
          <label>
            Parola
            <input type="password" autoComplete="current-password" required value={form.password}
              onChange={event => setForm({ ...form, password: event.target.value })} placeholder="••••••••••••" />
          </label>
          {error && <div className="error-banner">{error}</div>}
          <button className="primary-button login-button" disabled={busy}>
            {busy ? <><RefreshCw size={15} className="spin" /> Giriş yapılıyor</> : <><ShieldCheck size={16} /> Güvenli giriş</>}
          </button>
        </form>
        <p className="security-note">Bu panel yalnızca yetkili yöneticiler içindir. Oturum tarayıcı kapatıldığında sona erer.</p>
      </section>
    </main>
  )
}

function DashboardShell({ children }) {
  const location = useLocation()
  const { manager, logout } = useAuth()
  const [menuOpen, setMenuOpen] = useState(false)

  return (
    <div className="app-shell">
      <aside className={`sidebar ${menuOpen ? 'open' : ''}`}>
        <div className="sidebar-brand">
          <div className="brand-mark">ق</div>
          <div><strong>Kuran AI</strong><span>Operations</span></div>
          <button className="icon-button mobile-close" onClick={() => setMenuOpen(false)} aria-label="Menüyü kapat"><X size={18} /></button>
        </div>
        <nav>
          <p className="nav-label">ÇALIŞMA ALANI</p>
          {navigation.map(item => {
            const active = item.path === '/' ? location.pathname === '/' : location.pathname.startsWith(item.path)
            const Icon = item.icon
            return <Link key={item.path} to={item.path} onClick={() => setMenuOpen(false)} className={active ? 'active' : ''}><Icon size={17} />{item.label}</Link>
          })}
        </nav>
        <div className="sidebar-bottom">
          <div className="status-chip"><span className="status-dot" /> Canlı sistem</div>
          <div className="manager-card">
            <Avatar name={manager?.name || manager?.email || 'M'} />
            <div><strong>{manager?.name || 'Yönetici'}</strong><span>{manager?.email}</span></div>
            <button className="icon-button" onClick={logout} title="Çıkış yap"><LogOut size={16} /></button>
          </div>
        </div>
      </aside>
      {menuOpen && <button className="sidebar-backdrop" onClick={() => setMenuOpen(false)} aria-label="Menüyü kapat" />}
      <div className="workspace">
        <header className="topbar">
          <button className="icon-button mobile-menu" onClick={() => setMenuOpen(true)} aria-label="Menüyü aç"><Menu size={19} /></button>
          <div className="breadcrumb"><span>Workspace</span><b>/</b><strong>{pageName(location.pathname)}</strong></div>
          <a className="site-link" href="https://kuranai.talhacan.com" target="_blank" rel="noreferrer">Canlı site <ExternalLink size={14} /></a>
        </header>
        <div className="page-container">{children}</div>
      </div>
    </div>
  )
}

function pageName(pathname) {
  if (pathname.startsWith('/users/')) return 'Kullanıcı detayı'
  return navigation.find(item => item.path === '/' ? pathname === '/' : pathname.startsWith(item.path))?.label || 'Panel'
}

function useRemote(path, pollInterval = 0) {
  const { token, logout } = useAuth()
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [version, setVersion] = useState(0)

  const reload = useCallback(() => setVersion(value => value + 1), [])

  useEffect(() => {
    let active = true
    const load = (quiet = false) => {
      if (!quiet) setLoading(true)
      setError('')
      request(path, {}, token)
        .then(result => { if (active) setData(result) })
        .catch(err => {
          if (!active) return
          if (err.status === 401) logout()
          else setError(err.message)
        })
        .finally(() => { if (active && !quiet) setLoading(false) })
    }
    load()
    const timer = pollInterval ? window.setInterval(() => load(true), pollInterval) : null
    return () => {
      active = false
      if (timer) window.clearInterval(timer)
    }
  }, [path, pollInterval, token, logout, version])

  return { data, loading, error, reload }
}

function OverviewPage() {
  const { data, loading, error, reload } = useRemote('/api/manager/dashboard/overview?days=14', 15000)
  if (loading) return <PageLoader />
  if (error) return <ErrorState message={error} onRetry={reload} />
  const metrics = data?.metrics || {}
  const maxTrend = Math.max(1, ...(data?.trend || []).map(point => point.questions))

  return (
    <>
      <PageHeader eyebrow="GERÇEK ZAMANLI OPERASYON" title="Genel Bakış"
        description="Kullanıcı etkinliği, platform dağılımı ve abonelik performansını tek görünümde izleyin."
        action={<button className="secondary-button" onClick={reload}><RefreshCw size={14} /> Yenile</button>} />
      <section className="metric-grid">
        <MetricCard label="Toplam kullanıcı" value={number(metrics.totalUsers)} meta="Tüm kayıtlı hesaplar" icon={Users} tone="violet" />
        <MetricCard label="Şu an online" value={number(metrics.onlineUsers)} meta="Son 2 dakikada etkin" icon={Activity} tone="green" live />
        <MetricCard label="Aktif abonelik" value={number(metrics.activeSubscriptions)} meta="Süresi devam eden paketler" icon={CreditCard} tone="blue" />
        <MetricCard label="Bugünkü sorular" value={number(metrics.questionsToday)} meta="Kullanıcı mesajları" icon={MessageSquareText} tone="amber" />
        <MetricCard label="Brüt paket değeri" value={money(metrics.grossRevenueTL)} meta="Kayıtlı abonelik toplamı" icon={CircleDollarSign} tone="rose" />
      </section>

      <section className="dashboard-grid">
        <Panel className="activity-panel" title="Soru etkinliği" subtitle="Son 14 gün">
          <div className="chart-legend"><span><i className="legend-question" /> Sorular</span><span><i className="legend-user" /> Yeni kullanıcılar</span></div>
          <div className="bar-chart">
            {(data?.trend || []).map(point => (
              <div className="bar-column" key={point.day}>
                <div className="bar-stack">
                  <div className="bar questions" style={{ height: `${Math.max(3, (point.questions / maxTrend) * 100)}%` }} title={`${point.questions} soru`} />
                  <div className="bar users" style={{ height: `${Math.max(2, (point.newUsers / maxTrend) * 100)}%` }} title={`${point.newUsers} yeni kullanıcı`} />
                </div>
                <span>{shortDate(point.day)}</span>
              </div>
            ))}
          </div>
        </Panel>
        <Panel title="Platform dağılımı" subtitle="Kullanıcıların son oturumu">
          <Breakdown rows={data?.platforms || []} labels={{ ios: 'iOS', android: 'Android', web: 'Web', unknown: 'Bilinmiyor' }} />
        </Panel>
      </section>

      <section className="dashboard-grid lower">
        <Panel title="Online kullanıcılar" subtitle={`${data?.presenceWindowSeconds || 120} saniyelik canlı pencere`}>
          <div className="compact-list">
            {(data?.onlineUsers || []).length === 0 && <EmptyInline text="Şu anda etkin kullanıcı yok." />}
            {(data?.onlineUsers || []).slice(0, 7).map(user => (
              <Link to={`/users/${user.id}`} className="compact-user" key={user.id}>
                <Avatar name={user.name || user.email} online />
                <div><strong>{user.name || user.username}</strong><span>{user.email}</span></div>
                <Platform platform={user.platform} />
              </Link>
            ))}
          </div>
        </Panel>
        <Panel title="Son satın almalar" subtitle="Apple ve Google kaynakları">
          <div className="compact-list">
            {(data?.recentPurchases || []).length === 0 && <EmptyInline text="Henüz satın alma yok." />}
            {(data?.recentPurchases || []).map(item => (
              <div className="purchase-row" key={item.id}>
                <ProviderIcon provider={item.provider} />
                <div><strong>{item.userName}</strong><span>{item.packageName} · {billingLabel(item.billing)}</span></div>
                <b>{money(item.amountTL)}</b>
              </div>
            ))}
          </div>
        </Panel>
      </section>
    </>
  )
}

function UsersPage() {
  const [filters, setFilters] = useState({ search: '', platform: 'all', status: 'all', page: 1 })
  const query = new URLSearchParams({
    page: String(filters.page), pageSize: '20', search: filters.search, platform: filters.platform, status: filters.status,
  }).toString()
  const { data, loading, error, reload } = useRemote(`/api/manager/users?${query}`)

  return (
    <>
      <PageHeader eyebrow="KİMLİK VE ERİŞİM" title="Kullanıcılar" description="Hesapları, cihaz platformlarını, abonelikleri ve sohbet etkinliğini inceleyin."
        action={<span className="count-badge">{number(data?.total || 0)} kullanıcı</span>} />
      <FilterBar>
        <SearchInput value={filters.search} onChange={value => setFilters({ ...filters, search: value, page: 1 })} placeholder="Ad, e-posta veya kullanıcı adı ara" />
        <Select value={filters.platform} onChange={value => setFilters({ ...filters, platform: value, page: 1 })}
          options={[['all', 'Tüm platformlar'], ['ios', 'iOS'], ['android', 'Android'], ['web', 'Web'], ['unknown', 'Bilinmiyor']]} />
        <Select value={filters.status} onChange={value => setFilters({ ...filters, status: value, page: 1 })}
          options={[['all', 'Tüm durumlar'], ['active', 'Aktif'], ['inactive', 'Pasif']]} />
        <button className="icon-button bordered" onClick={reload} aria-label="Yenile"><RefreshCw size={15} /></button>
      </FilterBar>
      <DataPanel>
        {loading ? <TableLoader /> : error ? <ErrorState message={error} onRetry={reload} compact /> : (
          <>
            <div className="table-scroll">
              <table>
                <thead><tr><th>Kullanıcı</th><th>Platform</th><th>Paket</th><th>Sohbet / Soru</th><th>Son etkinlik</th><th>Durum</th><th /></tr></thead>
                <tbody>
                  {(data?.items || []).map(user => (
                    <tr key={user.id}>
                      <td><UserCell user={user} /></td>
                      <td><Platform platform={user.platform} /></td>
                      <td><div className="stacked-cell"><strong>{user.packageName}</strong><span>{providerLabel(user.provider)}</span></div></td>
                      <td><span className="tabular">{number(user.sessionCount)} / {number(user.questionCount)}</span></td>
                      <td>{user.online ? <LiveLabel /> : <span className="muted">{relativeDate(user.lastSeenAt || user.createdAt)}</span>}</td>
                      <td><StatusBadge status={user.status === 1 ? 'active' : 'inactive'} /></td>
                      <td><Link className="row-action" to={`/users/${user.id}`} aria-label="Detayı aç"><Eye size={16} /></Link></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {(data?.items || []).length === 0 && <EmptyState title="Kullanıcı bulunamadı" text="Filtreleri değiştirerek tekrar deneyin." />}
            <Pagination page={data?.page || 1} totalPages={data?.totalPages || 1}
              onChange={page => setFilters({ ...filters, page })} />
          </>
        )}
      </DataPanel>
    </>
  )
}

function UserDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { data, loading, error, reload } = useRemote(`/api/manager/users/${id}`, 15000)
  const conversations = useRemote(`/api/manager/users/${id}/conversations?page=1&pageSize=50`)
  const [sessionID, setSessionID] = useState(null)
  if (loading) return <PageLoader />
  if (error) return <ErrorState message={error} onRetry={reload} />
  const user = data?.user || {}

  return (
    <>
      <button className="back-button" onClick={() => navigate('/users')}><ArrowLeft size={15} /> Kullanıcılara dön</button>
      <PageHeader eyebrow={`KULLANICI #${user.id}`} title={user.name || user.username}
        description={user.email} action={user.online ? <LiveLabel /> : <StatusBadge status={user.status === 1 ? 'active' : 'inactive'} />} />
      <section className="user-summary">
        <div className="profile-card">
          <Avatar name={user.name || user.email} online={user.online} large />
          <div><h3>{user.name}</h3><p>@{user.username}</p></div>
          <Platform platform={user.platform} />
        </div>
        <MiniStat label="Toplam sohbet" value={number(user.sessionCount)} />
        <MiniStat label="Toplam soru" value={number(user.questionCount)} />
        <MiniStat label="Üyelik" value={user.packageName || 'Ücretsiz'} />
        <MiniStat label="Kayıt tarihi" value={formatDate(user.createdAt)} />
      </section>

      <Panel title="Sohbet geçmişi" subtitle="Bir sohbeti seçerek tüm mesaj akışını inceleyin" noPadding>
        {conversations.loading ? <TableLoader /> : (
          <div className="table-scroll">
            <table>
              <thead><tr><th>Başlangıç</th><th>Kaynak</th><th>Mesaj</th><th>Son güncelleme</th><th /></tr></thead>
              <tbody>
                {(conversations.data?.items || []).map(item => (
                  <tr key={item.id} className="clickable-row" onClick={() => setSessionID(item.id)}>
                    <td><div className="question-preview"><strong>{item.firstMessage || 'Yeni sohbet'}</strong><span>{truncate(item.lastMessage, 70)}</span></div></td>
                    <td><SourceBadge source={item.sourceType} /></td>
                    <td className="tabular">{item.messageCount}</td>
                    <td className="muted">{formatDateTime(item.updatedAt)}</td>
                    <td><button className="row-action"><Eye size={16} /></button></td>
                  </tr>
                ))}
              </tbody>
            </table>
            {(conversations.data?.items || []).length === 0 && <EmptyState title="Sohbet bulunamadı" text="Bu kullanıcı henüz sohbet başlatmamış." />}
          </div>
        )}
      </Panel>

      <section className="dashboard-grid lower">
        <Panel title="Cihazlar ve oturumlar" subtitle="Hassas token değerleri gösterilmez" noPadding>
          <div className="table-scroll">
            <table><thead><tr><th>Platform</th><th>Cihaz</th><th>IP</th><th>Oluşturma</th></tr></thead>
              <tbody>{(data?.devices || []).map(device => <tr key={device.id}>
                <td><Platform platform={device.platform} /></td><td className="muted">{device.deviceInfo || '—'}</td>
                <td className="tabular muted">{device.ipAddress || '—'}</td><td className="muted">{formatDate(device.createdAt)}</td>
              </tr>)}</tbody>
            </table>
          </div>
        </Panel>
        <Panel title="Abonelik geçmişi" subtitle="Mağaza ve paket süreleri" noPadding>
          <div className="table-scroll">
            <table><thead><tr><th>Paket</th><th>Kaynak</th><th>Süre</th><th>Durum</th></tr></thead>
              <tbody>{(data?.subscriptions || []).map(item => <tr key={item.id}>
                <td><div className="stacked-cell"><strong>{item.packageName}</strong><span>{money(item.amountTL)}</span></div></td>
                <td><Provider provider={item.provider} /></td><td>{billingLabel(item.billing)}</td><td><StatusBadge status={item.status} /></td>
              </tr>)}</tbody>
            </table>
          </div>
        </Panel>
      </section>
      {sessionID && <ConversationDrawer sessionID={sessionID} onClose={() => setSessionID(null)} />}
    </>
  )
}

function PurchasesPage() {
  const [filters, setFilters] = useState({ search: '', provider: 'all', status: 'all', page: 1 })
  const query = new URLSearchParams({ ...filters, pageSize: '20' }).toString()
  const { data, loading, error, reload } = useRemote(`/api/manager/purchases?${query}`)
  return (
    <>
      <PageHeader eyebrow="GELİR OPERASYONLARI" title="Satın Almalar" description="Apple App Store ve Google Play aboneliklerini, sürelerini ve durumlarını izleyin."
        action={<span className="count-badge">{number(data?.total || 0)} kayıt</span>} />
      <FilterBar>
        <SearchInput value={filters.search} onChange={value => setFilters({ ...filters, search: value, page: 1 })} placeholder="Kullanıcı veya e-posta ara" />
        <Select value={filters.provider} onChange={value => setFilters({ ...filters, provider: value, page: 1 })}
          options={[['all', 'Tüm kaynaklar'], ['apple', 'App Store'], ['google', 'Google Play']]} />
        <Select value={filters.status} onChange={value => setFilters({ ...filters, status: value, page: 1 })}
          options={[['all', 'Tüm durumlar'], ['active', 'Aktif'], ['expired', 'Süresi dolmuş']]} />
        <button className="icon-button bordered" onClick={reload}><RefreshCw size={15} /></button>
      </FilterBar>
      <DataPanel>
        {loading ? <TableLoader /> : error ? <ErrorState message={error} onRetry={reload} compact /> : <>
          <div className="table-scroll"><table>
            <thead><tr><th>Kullanıcı</th><th>Paket</th><th>Kaynak</th><th>Dönem</th><th>Tutar</th><th>Bitiş</th><th>Durum</th></tr></thead>
            <tbody>{(data?.items || []).map(item => <tr key={item.id}>
              <td><div className="stacked-cell"><strong>{item.userName}</strong><span>{item.userEmail}</span></div></td>
              <td><div className="stacked-cell"><strong>{item.packageName}</strong><span>{item.productId || item.packageSlug}</span></div></td>
              <td><Provider provider={item.provider} /></td><td>{billingLabel(item.billing)}</td>
              <td className="tabular">{money(item.amountTL)}</td><td className="muted">{formatDate(item.expiresAt)}</td>
              <td><StatusBadge status={item.status} /></td>
            </tr>)}</tbody>
          </table></div>
          {(data?.items || []).length === 0 && <EmptyState title="Satın alma bulunamadı" text="Seçili filtrelere uygun kayıt yok." />}
          <Pagination page={data?.page || 1} totalPages={data?.totalPages || 1} onChange={page => setFilters({ ...filters, page })} />
        </>}
      </DataPanel>
    </>
  )
}

function QuestionsPage() {
  const [filters, setFilters] = useState({ search: '', sourceType: 'all', page: 1 })
  const [sessionID, setSessionID] = useState(null)
  const query = new URLSearchParams({ ...filters, pageSize: '20' }).toString()
  const { data, loading, error, reload } = useRemote(`/api/manager/questions?${query}`)
  return (
    <>
      <PageHeader eyebrow="İÇERİK GÖZLEMLEME" title="Kullanıcı Soruları"
        description="Sorulan soruları kaynağı ve verilen yanıtla birlikte inceleyin."
        action={<span className="count-badge">{number(data?.total || 0)} soru</span>} />
      <FilterBar>
        <SearchInput value={filters.search} onChange={value => setFilters({ ...filters, search: value, page: 1 })} placeholder="Soru, kullanıcı veya e-posta ara" />
        <Select value={filters.sourceType} onChange={value => setFilters({ ...filters, sourceType: value, page: 1 })}
          options={[['all', 'Tüm kaynaklar'], ['meal', 'Kur’an Meali'], ['tefsir', 'Tefsir'], ['hadis', 'Hadis']]} />
        <button className="icon-button bordered" onClick={reload}><RefreshCw size={15} /></button>
      </FilterBar>
      <DataPanel>
        {loading ? <TableLoader /> : error ? <ErrorState message={error} onRetry={reload} compact /> : <>
          <div className="table-scroll"><table>
            <thead><tr><th>Soru</th><th>Kullanıcı</th><th>Kaynak</th><th>Tarih</th><th /></tr></thead>
            <tbody>{(data?.items || []).map(item => <tr key={item.id}>
              <td><div className="question-preview"><strong>{truncate(item.content, 105)}</strong><span>{truncate(item.assistantText, 100) || 'Yanıt bekleniyor'}</span></div></td>
              <td><div className="stacked-cell"><strong>{item.userName}</strong><span>{item.userEmail}</span></div></td>
              <td><SourceBadge source={item.sourceType} /></td><td className="muted">{formatDateTime(item.createdAt)}</td>
              <td><button className="row-action" onClick={() => setSessionID(item.sessionId)}><Eye size={16} /></button></td>
            </tr>)}</tbody>
          </table></div>
          {(data?.items || []).length === 0 && <EmptyState title="Soru bulunamadı" text="Filtreleri değiştirerek tekrar deneyin." />}
          <Pagination page={data?.page || 1} totalPages={data?.totalPages || 1} onChange={page => setFilters({ ...filters, page })} />
        </>}
      </DataPanel>
      {sessionID && <ConversationDrawer sessionID={sessionID} onClose={() => setSessionID(null)} />}
    </>
  )
}

function ConversationDrawer({ sessionID, onClose }) {
  const { data, loading, error, reload } = useRemote(`/api/manager/conversations/${sessionID}/messages`)
  useEffect(() => {
    const close = event => { if (event.key === 'Escape') onClose() }
    document.addEventListener('keydown', close)
    return () => document.removeEventListener('keydown', close)
  }, [onClose])

  return (
    <div className="drawer-layer" role="dialog" aria-modal="true" aria-label="Sohbet geçmişi">
      <button className="drawer-backdrop" onClick={onClose} aria-label="Kapat" />
      <aside className="conversation-drawer">
        <header>
          <div><span className="eyebrow">SOHBET #{sessionID}</span><h2>{data?.session?.userName || 'Sohbet geçmişi'}</h2><p>{data?.session?.userEmail}</p></div>
          <button className="icon-button bordered" onClick={onClose}><X size={18} /></button>
        </header>
        <div className="conversation-body">
          {loading ? <PageLoader /> : error ? <ErrorState message={error} onRetry={reload} /> :
            (data?.messages || []).map(message => (
              <article className={`message ${message.role}`} key={message.id}>
                <div className="message-meta"><strong>{message.role === 'user' ? 'Kullanıcı' : 'Kuran AI'}</strong><span>{formatDateTime(message.createdAt)}</span></div>
                <p>{message.content}</p>
                <SourceBadge source={message.sourceType} />
              </article>
            ))}
        </div>
      </aside>
    </div>
  )
}

function PageHeader({ eyebrow, title, description, action }) {
  return <header className="page-header"><div><span className="eyebrow">{eyebrow}</span><h1>{title}</h1><p>{description}</p></div>{action && <div className="page-action">{action}</div>}</header>
}

function MetricCard({ label, value, meta, icon, tone, live }) {
  return <article className={`metric-card ${tone}`}><div className="metric-icon">{createElement(icon, { size: 17 })}</div><div className="metric-label">{label}{live && <span className="pulse-dot" />}</div><strong>{value}</strong><p>{meta}</p></article>
}

function MiniStat({ label, value }) {
  return <article className="mini-stat"><span>{label}</span><strong>{value}</strong></article>
}

function Panel({ title, subtitle, children, className = '', noPadding = false }) {
  return <section className={`panel ${className} ${noPadding ? 'no-padding' : ''}`}><header><div><h2>{title}</h2><p>{subtitle}</p></div></header><div className="panel-content">{children}</div></section>
}

function DataPanel({ children }) {
  return <section className="data-panel">{children}</section>
}

function FilterBar({ children }) {
  return <div className="filter-bar">{children}</div>
}

function SearchInput({ value, onChange, placeholder }) {
  return <label className="search-input"><Search size={15} /><input value={value} onChange={event => onChange(event.target.value)} placeholder={placeholder} /></label>
}

function Select({ value, onChange, options }) {
  return <select value={value} onChange={event => onChange(event.target.value)}>{options.map(([key, label]) => <option key={key} value={key}>{label}</option>)}</select>
}

function Pagination({ page, totalPages, onChange }) {
  if (totalPages <= 1) return null
  return <div className="pagination"><span>Sayfa {page} / {totalPages}</span><div>
    <button disabled={page <= 1} onClick={() => onChange(page - 1)}><ChevronLeft size={15} /></button>
    <button disabled={page >= totalPages} onClick={() => onChange(page + 1)}><ChevronRight size={15} /></button>
  </div></div>
}

function UserCell({ user }) {
  return <div className="user-cell"><Avatar name={user.name || user.email} online={user.online} /><div><strong>{user.name || user.username}</strong><span>{user.email}</span></div></div>
}

function Avatar({ name = '', online = false, large = false }) {
  const initials = name.split(/\s+/).filter(Boolean).slice(0, 2).map(part => part[0]).join('').toUpperCase() || 'K'
  return <span className={`avatar ${large ? 'large' : ''}`}>{initials}{online && <i />}</span>
}

function Platform({ platform }) {
  const normalized = platform || 'unknown'
  const Icon = normalized === 'web' ? MonitorSmartphone : Smartphone
  return <span className={`platform ${normalized}`}><Icon size={13} />{normalized === 'ios' ? 'iOS' : normalized === 'android' ? 'Android' : normalized === 'web' ? 'Web' : 'Bilinmiyor'}</span>
}

function ProviderIcon({ provider }) {
  return <span className={`provider-icon ${provider}`}><CreditCard size={15} /></span>
}

function Provider({ provider }) {
  return <span className={`provider ${provider}`}><ProviderIcon provider={provider} />{providerLabel(provider)}</span>
}

function providerLabel(provider) {
  if (provider === 'apple') return 'App Store'
  if (provider === 'google') return 'Google Play'
  return '—'
}

function StatusBadge({ status }) {
  const label = status === 'active' ? 'Aktif' : status === 'expired' ? 'Süresi dolmuş' : status === 'inactive' ? 'Pasif' : status
  return <span className={`status-badge ${status}`}>{label}</span>
}

function SourceBadge({ source }) {
  const normalized = source === 'kuran' ? 'meal' : source
  return <span className={`source-badge ${normalized}`}>{normalized === 'meal' ? 'Kur’an Meali' : normalized === 'tefsir' ? 'Tefsir' : normalized === 'hadis' ? 'Hadis' : normalized || 'Genel'}</span>
}

function LiveLabel() {
  return <span className="live-label"><i /> Online</span>
}

function Breakdown({ rows, labels }) {
  const total = Math.max(1, rows.reduce((sum, row) => sum + row.count, 0))
  return <div className="breakdown">{rows.map(row => <div key={row.key}>
    <div className="breakdown-label"><span><i className={`breakdown-dot ${row.key}`} />{labels[row.key] || row.key}</span><b>{number(row.count)}</b></div>
    <div className="progress"><i className={row.key} style={{ width: `${(row.count / total) * 100}%` }} /></div>
  </div>)}</div>
}

function FullLoader({ label }) {
  return <div className="full-loader"><div className="brand-mark large">ق</div><RefreshCw className="spin" size={18} /><span>{label}</span></div>
}

function PageLoader() {
  return <div className="page-loader"><RefreshCw className="spin" size={18} /><span>Veriler yükleniyor</span></div>
}

function TableLoader() {
  return <div className="table-loader">{[1, 2, 3, 4, 5].map(item => <i key={item} />)}</div>
}

function ErrorState({ message, onRetry, compact = false }) {
  return <div className={`error-state ${compact ? 'compact' : ''}`}><Activity size={24} /><h3>Veri yüklenemedi</h3><p>{message}</p><button className="secondary-button" onClick={onRetry}><RefreshCw size={14} /> Tekrar dene</button></div>
}

function EmptyState({ title, text }) {
  return <div className="empty-state"><Search size={24} /><h3>{title}</h3><p>{text}</p></div>
}

function EmptyInline({ text }) {
  return <div className="empty-inline">{text}</div>
}

function number(value) {
  return new Intl.NumberFormat('tr-TR').format(Number(value || 0))
}

function money(value) {
  return new Intl.NumberFormat('tr-TR', { style: 'currency', currency: 'TRY', maximumFractionDigits: 2 }).format(Number(value || 0))
}

function formatDate(value) {
  if (!value) return '—'
  return new Intl.DateTimeFormat('tr-TR', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(value))
}

function formatDateTime(value) {
  if (!value) return '—'
  return new Intl.DateTimeFormat('tr-TR', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' }).format(new Date(value))
}

function relativeDate(value) {
  if (!value) return 'Bilinmiyor'
  const minutes = Math.round((new Date(value).getTime() - Date.now()) / 60000)
  if (Math.abs(minutes) < 60) return new Intl.RelativeTimeFormat('tr', { numeric: 'auto' }).format(minutes, 'minute')
  const hours = Math.round(minutes / 60)
  if (Math.abs(hours) < 24) return new Intl.RelativeTimeFormat('tr', { numeric: 'auto' }).format(hours, 'hour')
  return formatDate(value)
}

function shortDate(value) {
  return new Intl.DateTimeFormat('tr-TR', { day: '2-digit', month: 'short' }).format(new Date(value))
}

function billingLabel(value) {
  return value === 'monthly' ? 'Aylık' : value === 'yearly' ? 'Yıllık' : value || '—'
}

function truncate(value = '', max = 80) {
  return value.length > max ? `${value.slice(0, max).trim()}…` : value
}

export default App

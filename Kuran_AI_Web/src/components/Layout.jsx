import { Outlet, Link, useLocation } from 'react-router-dom'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Layout.module.css'

export default function Layout() {
  const { pathname } = useLocation()
  const { language, setLanguage, text } = useLanguage()
  const isChatRoute = pathname === '/app'

  const navItems = [
    { to: '/app', label: text('Sohbet', 'Chat', 'المحادثة') },
    { to: '/packages', label: text('Paketler', 'Plans', 'الباقات') },
    { to: '/support', label: text('Destek', 'Support', 'الدعم') },
    { to: '/contact', label: text('İletişim', 'Contact', 'اتصل بنا') },
  ]

  return (
    <div className={`${styles.wrapper} ${isChatRoute ? styles.chatWrapper : ''}`}>
      <header className={styles.header}>
        <div className={styles.headerInner}>
          <Link to="/" className={styles.logo} aria-label={text('Kuran AI ana sayfa', 'Kuran AI home', 'الصفحة الرئيسية لـ Kuran AI')}>
            <span className={styles.logoMark}>ق</span>
            <span className={styles.logoText}>Kuran AI</span>
          </Link>

          <div className={styles.headerActions}>
            <nav className={styles.nav} aria-label={text('Ana menü', 'Main navigation', 'التنقل الرئيسي')}>
              {navItems.map(item => (
                <Link key={item.to} to={item.to} className={pathname === item.to ? styles.active : ''}>
                  {item.label}
                </Link>
              ))}
            </nav>
            <div className={styles.languageSwitch} aria-label={text('Dil seçimi', 'Language selection', 'اختيار اللغة')}>
              <button
                type="button"
                className={language === 'tr' ? styles.languageActive : ''}
                onClick={() => setLanguage('tr')}
                aria-pressed={language === 'tr'}
              >TR</button>
              <span>/</span>
              <button
                type="button"
                className={language === 'en' ? styles.languageActive : ''}
                onClick={() => setLanguage('en')}
                aria-pressed={language === 'en'}
              >EN</button>
              <span>/</span>
              <button
                type="button"
                className={language === 'ar' ? styles.languageActive : ''}
                onClick={() => setLanguage('ar')}
                aria-pressed={language === 'ar'}
              >AR</button>
            </div>
          </div>
        </div>
      </header>

      <main className={pathname === '/' ? styles.mainLanding : isChatRoute ? styles.mainApp : pathname === '/packages' ? styles.mainPackages : styles.mainPage}>
        <Outlet />
      </main>

      {!isChatRoute && (
        <footer className={styles.footer}>
          <div className={styles.footerInner}>
            <div className={styles.footerLinks}>
              <Link to="/privacy-policy">{text('Gizlilik Politikası', 'Privacy Policy', 'سياسة الخصوصية')}</Link>
              <span className={styles.dot}>·</span>
              <Link to="/terms">{text('Kullanım Koşulları', 'Terms of Use', 'شروط الاستخدام')}</Link>
              <span className={styles.dot}>·</span>
              <Link to="/packages">{text('Paketler', 'Plans', 'الباقات')}</Link>
              <span className={styles.dot}>·</span>
              <Link to="/support">{text('Destek', 'Support', 'الدعم')}</Link>
              <span className={styles.dot}>·</span>
              <Link to="/contact">{text('İletişim', 'Contact', 'اتصل بنا')}</Link>
            </div>
            <p className={styles.copyright}>© {new Date().getFullYear()} Kuran AI. {text('Tüm hakları saklıdır.', 'All rights reserved.', 'جميع الحقوق محفوظة.')}</p>
          </div>
        </footer>
      )}
    </div>
  )
}

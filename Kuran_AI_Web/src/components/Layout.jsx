import { Outlet, Link, useLocation } from 'react-router-dom'
import { useLanguage } from '../i18n/LanguageContext'
import styles from './Layout.module.css'

export default function Layout() {
  const { pathname } = useLocation()
  const { language, setLanguage, text } = useLanguage()

  const navItems = [
    { to: '/app', label: text('Sohbet', 'Chat') },
    { to: '/support', label: text('Destek', 'Support') },
    { to: '/contact', label: text('İletişim', 'Contact') },
  ]

  return (
    <div className={styles.wrapper}>
      <header className={styles.header}>
        <div className={styles.headerInner}>
          <Link to="/" className={styles.logo} aria-label={text('Kuran AI ana sayfa', 'Kuran AI home')}>
            <span className={styles.logoMark}>ق</span>
            <span className={styles.logoText}>Kuran AI</span>
          </Link>

          <div className={styles.headerActions}>
            <nav className={styles.nav} aria-label={text('Ana menü', 'Main navigation')}>
              {navItems.map(item => (
                <Link key={item.to} to={item.to} className={pathname === item.to ? styles.active : ''}>
                  {item.label}
                </Link>
              ))}
            </nav>
            <div className={styles.languageSwitch} aria-label={text('Dil seçimi', 'Language selection')}>
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
            </div>
          </div>
        </div>
      </header>

      <main className={pathname === '/' ? styles.mainLanding : pathname === '/app' ? styles.mainApp : styles.mainPage}>
        <Outlet />
      </main>

      <footer className={styles.footer}>
        <div className={styles.footerInner}>
          <div className={styles.footerLinks}>
            <Link to="/privacy-policy">{text('Gizlilik Politikası', 'Privacy Policy')}</Link>
            <span className={styles.dot}>·</span>
            <Link to="/terms">{text('Kullanım Koşulları', 'Terms of Use')}</Link>
            <span className={styles.dot}>·</span>
            <Link to="/support">{text('Destek', 'Support')}</Link>
            <span className={styles.dot}>·</span>
            <Link to="/contact">{text('İletişim', 'Contact')}</Link>
          </div>
          <p className={styles.copyright}>© {new Date().getFullYear()} Kuran AI. {text('Tüm hakları saklıdır.', 'All rights reserved.')}</p>
        </div>
      </footer>
    </div>
  )
}

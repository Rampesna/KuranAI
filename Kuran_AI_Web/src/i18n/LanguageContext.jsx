import { createContext, useContext, useEffect, useMemo, useState } from 'react'

const LanguageContext = createContext(null)

export function LanguageProvider({ children }) {
  const [language, setLanguageState] = useState(() => localStorage.getItem('kuranai_language') || 'tr')

  const setLanguage = (nextLanguage) => {
    const next = nextLanguage === 'en' ? 'en' : 'tr'
    localStorage.setItem('kuranai_language', next)
    setLanguageState(next)
  }

  useEffect(() => {
    document.documentElement.lang = language
    document.title = language === 'tr'
      ? 'Kuran AI — Kur’an, Tefsir ve Hadis Asistanı'
      : 'Kuran AI — Quran, Tafsir and Hadith Assistant'
  }, [language])

  const value = useMemo(() => ({
    language,
    setLanguage,
    isTurkish: language === 'tr',
    text: (tr, en) => language === 'tr' ? tr : en,
  }), [language])

  return <LanguageContext.Provider value={value}>{children}</LanguageContext.Provider>
}

// The hook intentionally lives beside its provider so consumers share one context instance.
// eslint-disable-next-line react-refresh/only-export-components
export function useLanguage() {
  const context = useContext(LanguageContext)
  if (!context) throw new Error('useLanguage must be used inside LanguageProvider')
  return context
}

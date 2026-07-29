import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { arabicText } from './ar'

const LanguageContext = createContext(null)

export function LanguageProvider({ children }) {
  const [language, setLanguageState] = useState(() => {
    const stored = localStorage.getItem('kuranai_language')
    return ['tr', 'en', 'ar'].includes(stored) ? stored : 'tr'
  })

  const setLanguage = (nextLanguage) => {
    const next = ['tr', 'en', 'ar'].includes(nextLanguage) ? nextLanguage : 'tr'
    localStorage.setItem('kuranai_language', next)
    setLanguageState(next)
  }

  useEffect(() => {
    document.documentElement.lang = language
    document.documentElement.dir = language === 'ar' ? 'rtl' : 'ltr'
    document.title = language === 'tr'
      ? 'Kuran AI — Kur’an, Tefsir ve Hadis Asistanı'
      : language === 'ar'
        ? 'Kuran AI — مساعد القرآن والتفسير والحديث'
        : 'Kuran AI — Quran, Tafsir and Hadith Assistant'
  }, [language])

  const value = useMemo(() => ({
    language,
    setLanguage,
    isTurkish: language === 'tr',
    isArabic: language === 'ar',
    text: (tr, en, ar) => language === 'tr' ? tr : language === 'ar' ? (ar || arabicText(tr, en)) : en,
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

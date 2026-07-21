import SwiftUI

enum SourceFilter: String, CaseIterable, Identifiable {
    case kuran, tefsir, hadis

    var id: String { rawValue }

    func title(for language: Language) -> String {
        switch self {
        case .kuran: return language.localized("Kur'an", "Qur'an")
        case .tefsir: return language.localized("Tefsir", "Tafsir")
        case .hadis: return language.localized("Hadis", "Hadith")
        }
    }

    func caption(for language: Language) -> String {
        switch self {
        case .kuran: return language.localized("Mealler ve bağlamdan cevaplar.", "Answers grounded in translations & context.")
        case .tefsir: return language.localized("Klasik ve modern tefsir özetleri.", "Summaries from classical & modern exegesis.")
        case .hadis: return language.localized("Sahih hadis referansları.", "References from authenticated hadith sources.")
        }
    }
}

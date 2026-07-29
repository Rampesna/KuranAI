import SwiftUI

enum Language: String, CaseIterable, Identifiable {
    case turkish = "tr"
    case english = "en"
    case arabic = "ar"

    var id: String { rawValue }

    var displayName: String {
        switch self {
        case .turkish: return "Türkçe"
        case .english: return "English"
        case .arabic: return "العربية"
        }
    }

    var localeIdentifier: String {
        switch self {
        case .turkish: return "tr_TR"
        case .english: return "en_US"
        case .arabic: return "ar_SA"
        }
    }

    var speechLanguage: String {
        switch self {
        case .turkish: return "tr-TR"
        case .english: return "en-US"
        case .arabic: return "ar-SA"
        }
    }

    var layoutDirection: LayoutDirection {
        self == .arabic ? .rightToLeft : .leftToRight
    }
}

extension Language {
    func localized(_ tr: String, _ en: String, _ ar: String? = nil) -> String {
        switch self {
        case .turkish: return tr
        case .english: return en
        case .arabic: return ar ?? ArabicTranslations.resolve(tr: tr, fallback: en)
        }
    }
}

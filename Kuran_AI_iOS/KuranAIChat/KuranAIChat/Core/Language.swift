import SwiftUI

enum Language: String, CaseIterable, Identifiable {
    case turkish = "tr"
    case english = "en"

    var id: String { rawValue }

    var displayName: String {
        switch self {
        case .turkish: return "Türkçe"
        case .english: return "English"
        }
    }
}

extension Language {
    func localized(_ tr: String, _ en: String) -> String {
        switch self {
        case .turkish: return tr
        case .english: return en
        }
    }
}

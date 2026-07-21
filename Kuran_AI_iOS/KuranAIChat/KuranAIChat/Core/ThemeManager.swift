import Combine
import SwiftUI

// MARK: - Theme Mode

enum ThemeMode: String, CaseIterable, Identifiable {
    case system
    case light
    case dark

    var id: String { rawValue }

    func displayName(_ language: Language) -> String {
        switch self {
        case .system: return language.localized("Sistem", "System")
        case .light:  return language.localized("Açık", "Light")
        case .dark:   return language.localized("Koyu", "Dark")
        }
    }

    var icon: String {
        switch self {
        case .system: return "gear"
        case .light:  return "sun.max.fill"
        case .dark:   return "moon.fill"
        }
    }
}

// MARK: - Theme Manager

final class ThemeManager: ObservableObject {
    static let shared = ThemeManager()

    private static let key = "kuranai_theme_mode"

    @Published var themeMode: ThemeMode {
        didSet {
            UserDefaults.standard.set(themeMode.rawValue, forKey: Self.key)
        }
    }

    /// Root view'a `.preferredColorScheme(...)` olarak uygulanır.
    /// `nil` → sistem ayarına uy.
    var colorScheme: ColorScheme? {
        switch themeMode {
        case .system: return nil
        case .light:  return .light
        case .dark:   return .dark
        }
    }

    private init() {
        let raw = UserDefaults.standard.string(forKey: Self.key) ?? ThemeMode.system.rawValue
        self.themeMode = ThemeMode(rawValue: raw) ?? .system
    }
}

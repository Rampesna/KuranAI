import SwiftUI

enum Theme {
    // MARK: - Gold Accents (adaptive – light modda daha koyu amber, okunurluğu artırır)
    static let neonPrimary = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor(red: 1.0, green: 0.86, blue: 0.3, alpha: 1)
            : UIColor(red: 0.62, green: 0.46, blue: 0.0, alpha: 1) // Koyu amber
    })

    static let neonSecondary = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor(red: 0.96, green: 0.63, blue: 0.08, alpha: 1)
            : UIColor(red: 0.52, green: 0.36, blue: 0.0, alpha: 1) // Derin amber
    })

    static let neon = LinearGradient(
        colors: [neonPrimary, neonSecondary],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )

    // MARK: - Adaptive Colors

    /// Ana arka plan (dark: sıcak siyah, light: krem)
    static let background = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor(red: 0.05, green: 0.03, blue: 0.01, alpha: 1)
            : UIColor(red: 0.98, green: 0.965, blue: 0.935, alpha: 1)
    })

    /// Birincil metin rengi
    static let textPrimary = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? .white
            : UIColor(red: 0.12, green: 0.10, blue: 0.06, alpha: 1)
    })

    /// İkincil / sönük metin rengi
    static let textMuted = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor.white.withAlphaComponent(0.7)
            : UIColor(red: 0.42, green: 0.37, blue: 0.28, alpha: 1)
    })

    /// Cam efekti arka plan
    static let glassColor = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor.white.withAlphaComponent(0.08)
            : UIColor(red: 0.55, green: 0.45, blue: 0.25, alpha: 0.06)
    })

    /// Pill / küçük kart arka planı
    static let pillColor = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor.white.withAlphaComponent(0.12)
            : UIColor(red: 0.55, green: 0.45, blue: 0.25, alpha: 0.08)
    })

    /// Kenar çizgisi
    static let border = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor.white.withAlphaComponent(0.18)
            : UIColor(red: 0.55, green: 0.45, blue: 0.25, alpha: 0.15)
    })

    /// Kart arka planı (satır arka planları, inputlar vb.)
    static let cardBackground = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor.white.withAlphaComponent(0.06)
            : UIColor.white
    })

    /// Sidebar arka planı
    static let sidebarBackground = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor(red: 0.06, green: 0.05, blue: 0.04, alpha: 1)
            : UIColor(red: 0.97, green: 0.955, blue: 0.93, alpha: 1)
    })

    /// Overlay karartma
    static let dimOverlay = Color(UIColor { trait in
        trait.userInterfaceStyle == .dark
            ? UIColor.black.withAlphaComponent(0.4)
            : UIColor.black.withAlphaComponent(0.25)
    })

    // MARK: - Gradient Shortcuts

    static let glassGradient = solidGradient(glassColor)
    static let pillGradient = solidGradient(pillColor)

    static let whiteCardGradient = solidGradient(
        Color(UIColor { trait in
            trait.userInterfaceStyle == .dark
                ? UIColor.white.withAlphaComponent(0.95)
                : UIColor.white
        })
    )

    static let alertGradient = solidGradient(Color.red.opacity(0.9))

    static func solidGradient(_ color: Color) -> LinearGradient {
        LinearGradient(
            colors: [color, color],
            startPoint: .topLeading,
            endPoint: .bottomTrailing
        )
    }
}

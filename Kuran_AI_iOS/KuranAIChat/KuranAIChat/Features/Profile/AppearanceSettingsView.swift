import SwiftUI

struct AppearanceSettingsView: View {
    let language: Language
    let onBack: () -> Void

    @ObservedObject private var themeManager = ThemeManager.shared

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            ScrollView {
                VStack(alignment: .leading, spacing: 20) {
                    Text(language.localized(
                        "Uygulamanın görünümünü tercihlerinize göre ayarlayın.",
                        "Customize the appearance of the app to your preference."
                    ))
                    .font(.system(size: 13))
                    .foregroundColor(Theme.textMuted)

                    VStack(spacing: 1) {
                        ForEach(ThemeMode.allCases) { mode in
                            themeRow(mode: mode)
                        }
                    }
                }
                .padding(20)
            }
        }
        .navigationTitle(language.localized("Görünüm", "Appearance"))
        .navigationBarTitleDisplayMode(.inline)
    }

    private func themeRow(mode: ThemeMode) -> some View {
        Button {
            withAnimation(.easeInOut(duration: 0.25)) {
                themeManager.themeMode = mode
            }
        } label: {
            HStack(spacing: 14) {
                Image(systemName: mode.icon)
                    .font(.system(size: 18))
                    .foregroundColor(themeManager.themeMode == mode ? Theme.neonPrimary : Theme.textMuted)
                    .frame(width: 28, alignment: .center)

                VStack(alignment: .leading, spacing: 2) {
                    Text(mode.displayName(language))
                        .font(.system(size: 15, weight: .medium))
                        .foregroundColor(Theme.textPrimary)
                    Text(themeSubtitle(mode))
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textMuted)
                }

                Spacer()

                if themeManager.themeMode == mode {
                    Image(systemName: "checkmark.circle.fill")
                        .font(.system(size: 20))
                        .foregroundColor(Theme.neonPrimary)
                } else {
                    Circle()
                        .stroke(Theme.border, lineWidth: 1.5)
                        .frame(width: 20, height: 20)
                }
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 14)
            .background(
                RoundedRectangle(cornerRadius: 10)
                    .fill(Theme.cardBackground)
            )
        }
        .buttonStyle(.plain)
    }

    private func themeSubtitle(_ mode: ThemeMode) -> String {
        switch mode {
        case .system:
            return language.localized(
                "Telefonunuzun ayarlarına göre otomatik değişir",
                "Automatically matches your device settings"
            )
        case .light:
            return language.localized(
                "Her zaman açık tema kullanır",
                "Always uses light theme"
            )
        case .dark:
            return language.localized(
                "Her zaman koyu tema kullanır",
                "Always uses dark theme"
            )
        }
    }
}

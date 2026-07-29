import SwiftUI

struct OnboardingView: View {
    @Binding var language: Language
    @Binding var selectedSource: SourceFilter
    let isAuthenticated: Bool
    let onProfileTap: () -> Void
    let onPrimaryAction: (String?) -> Void
    @State private var animateGlow = false

    var body: some View {
        VStack(alignment: .leading, spacing: 24) {
            HStack {
                Button {
                    onProfileTap()
                } label: {
                    Image(systemName: "person.circle.fill")
                        .font(.system(size: 28))
                        .foregroundStyle(Theme.neon)
                }
                Spacer()
                LanguagePicker(selectedLanguage: $language)
            }
            header
            heroCard
            SourceSelector(language: language, selected: $selectedSource)
            HomeCarousel(language: language, onPrimaryAction: onPrimaryAction)
            Spacer()
        }
        .padding(30)
        .gesture(
            DragGesture(minimumDistance: 50, coordinateSpace: .local)
                .onEnded { value in
                    let horizontal = value.translation.width
                    let vertical = abs(value.translation.height)
                    // Yatay hareket dikeyden büyükse swipe say
                    guard abs(horizontal) > vertical else { return }
                    let opensChat = language == .arabic ? horizontal > 50 : horizontal < -50
                    let opensProfile = language == .arabic ? horizontal < -50 : horizontal > 50
                    if opensChat {
                        onPrimaryAction(nil)
                    } else if opensProfile {
                        onProfileTap()
                    }
                }
        )
    }

    private var header: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(language.localized("KURAN • AI", "QURAN • AI"))
                .font(.system(size: 13, weight: .semibold, design: .rounded))
                .foregroundColor(Theme.textMuted)
                .kerning(4)
            Text(language.localized("Vahy ile konuş.", "Speak with\nRevelation."))
                .font(.system(size: 42, weight: .bold, design: .rounded))
                .foregroundColor(Theme.textPrimary)
            Text(language.localized("Kur'an, tefsir ve hadis kaynaklarını tek asistanda birleştiren ses odaklı deneyim.", "Voice-first experience unifying Quran, tafsir & hadith data."))
                .font(.system(size: 15))
                .foregroundColor(Theme.textMuted)
        }
    }

    private var heroCard: some View {
        Button {
            onPrimaryAction(nil)
        } label: {
            ZStack {
                RoundedRectangle(cornerRadius: 28)
                    .fill(Theme.glassGradient)
                    .overlay(
                        RoundedRectangle(cornerRadius: 28)
                            .stroke(Theme.border, lineWidth: 1)
                    )
                    .shadow(color: Color.black.opacity(0.4), radius: 24, y: 12)

                VStack(spacing: 14) {
                    ZStack {
                        Circle()
                            .stroke(Theme.border, lineWidth: 2)
                            .frame(width: 110, height: 110)
                            .scaleEffect(animateGlow ? 1.08 : 0.95)
                            .opacity(animateGlow ? 0.2 : 0.4)
                            .animation(.easeInOut(duration: 1.6).repeatForever(autoreverses: true), value: animateGlow)
                        Circle()
                            .fill(Theme.neon)
                            .frame(width: 76, height: 76)
                            .shadow(color: Color.yellow.opacity(0.55), radius: 20, y: 8)
                            .overlay(
                                Image(systemName: "mic.fill")
                                    .font(.system(size: 26, weight: .semibold))
                                    .foregroundColor(.black)
                            )
                    }
                    .onAppear { animateGlow = true }

                    Text(selectedSource.title(for: language))
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundColor(Theme.textPrimary)
                    Text(selectedSource.caption(for: language))
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textMuted)
                        .multilineTextAlignment(.center)
                        .padding(.horizontal, 24)
                    Text(language.localized("Basılı tut ve konuş ya da yazmaya başla", "Hold to speak or start typing"))
                        .font(.system(size: 13, weight: .medium))
                        .foregroundColor(Theme.neonPrimary.opacity(0.9))
                }
                .padding(22)
            }
            .frame(maxWidth: .infinity)
            .frame(height: 260)
        }
        .buttonStyle(.plain)
    }
}

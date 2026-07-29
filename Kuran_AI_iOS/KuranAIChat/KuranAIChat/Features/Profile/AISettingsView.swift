import SwiftUI

struct AISettingsView: View {
    let language: Language
    let onBack: () -> Void

    @AppStorage(TTSVoiceStorage.key) private var storedVoiceId: String = "ash"

    private var roundedBg: some View {
        RoundedRectangle(cornerRadius: 10).fill(Theme.cardBackground)
    }

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            ScrollView {
                VStack(alignment: .leading, spacing: 20) {
                    Text(language.localized("Kur'an cevaplarının sesli okunmasında kullanılacak sesi seçin.", "Choose the voice for reading Quran answers aloud."))
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textMuted)

                    VStack(alignment: .leading, spacing: 8) {
                        Text(language.localized("Sesli okuma sesi", "TTS Voice"))
                            .font(.system(size: 13, weight: .medium))
                            .foregroundColor(Theme.textMuted)

                        Picker("", selection: $storedVoiceId) {
                            ForEach(TTSVoiceOption.all) { voice in
                                Text(voice.displayName(for: language)).tag(voice.id)
                            }
                        }
                        .pickerStyle(.menu)
                        .tint(Theme.textPrimary)
                        .font(.system(size: 16))
                        .padding(.horizontal, 16)
                        .padding(.vertical, 14)
                        .frame(minHeight: 48)
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .background(roundedBg)
                    }
                }
                .padding(20)
            }
        }
        .navigationTitle(language.localized("AI ayarları", "AI settings"))
        .navigationBarTitleDisplayMode(.inline)
    }
}

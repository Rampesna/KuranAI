import Foundation

/// OpenAI TTS ses seçenekleri. Arka planda OpenAI `voice` parametresi kullanılır.
struct TTSVoiceOption: Identifiable, Hashable {
    let id: String
    let nameTr: String
    let nameEn: String

    func displayName(for language: Language) -> String {
        language.localized(nameTr, nameEn)
    }

    /// Varsayılan: Ash
    static let `default` = TTSVoiceOption(id: "ash", nameTr: "Ash (Varsayılan)", nameEn: "Ash (Default)")

    static let all: [TTSVoiceOption] = [
        TTSVoiceOption(id: "ash", nameTr: "Ash (Varsayılan)", nameEn: "Ash (Default)"),
        TTSVoiceOption(id: "onyx", nameTr: "Olgun erkek (Onyx)", nameEn: "Deep male (Onyx)"),
        TTSVoiceOption(id: "echo", nameTr: "Erkek (Echo)", nameEn: "Male (Echo)"),
        TTSVoiceOption(id: "alloy", nameTr: "Nötr (Alloy)", nameEn: "Neutral (Alloy)"),
        TTSVoiceOption(id: "fable", nameTr: "İngiliz aksanlı (Fable)", nameEn: "British (Fable)"),
        TTSVoiceOption(id: "nova", nameTr: "Kadın (Nova)", nameEn: "Female (Nova)"),
        TTSVoiceOption(id: "shimmer", nameTr: "Kadın yumuşak (Shimmer)", nameEn: "Female soft (Shimmer)"),
    ]
}

/// Kullanıcının seçtiği TTS sesini UserDefaults'ta tutar.
enum TTSVoiceStorage {
    static let key = "tts_voice"

    static var selectedVoice: TTSVoiceOption {
        get {
            guard let id = UserDefaults.standard.string(forKey: key),
                  let v = TTSVoiceOption.all.first(where: { $0.id == id }) else {
                return .default
            }
            return v
        }
        set {
            UserDefaults.standard.set(newValue.id, forKey: key)
        }
    }

    static var selectedVoiceId: String {
        selectedVoice.id
    }
}

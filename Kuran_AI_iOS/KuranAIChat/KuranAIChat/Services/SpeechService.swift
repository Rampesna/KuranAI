import AVFoundation
import Combine
import Foundation

/// LLM cevaplarını sesli okur. Önce OpenAI TTS (backend) dener, başarısız olursa yerel AVSpeechSynthesizer kullanır.
/// Uzun metinler parçalara bölünerek 1–2 saniyede okumaya başlar.
final class SpeechService: NSObject, ObservableObject {
    static let shared = SpeechService()

    private let synthesizer = AVSpeechSynthesizer()
    private var audioPlayer: AVAudioPlayer?
    private var playbackContinuation: CheckedContinuation<Void, Never>?
    private var chunkedCancelled = false

    @Published private(set) var isSpeaking = false
    @Published private(set) var isLoading = false
    @Published var limitExceededInfo: APIError.LimitExceededInfo?

    override init() {
        super.init()
        synthesizer.delegate = self
    }

    /// Metni sesli okur. Uzun metinler parçalanır, ilk parça 1–2 saniyede çalar.
    func speak(_ text: String, language: String = "tr-TR", voice: String? = nil) {
        stop()
        let plain = stripMarkdown(text)
        guard !plain.isEmpty else { return }

        let voiceId = voice ?? TTSVoiceStorage.selectedVoiceId

        Task { @MainActor in
            await speakAsync(plain, language: language, voice: voiceId)
        }
    }

    private func speakAsync(_ plain: String, language: String, voice: String) async {
        isLoading = true
        chunkedCancelled = false
        defer { isLoading = false }

        let chunks = splitIntoParagraphs(plain)

        if chunks.count <= 1 {
            await playSingle(plain, voice: voice, language: language)
            return
        }

        // Parçalı çalma: ilk parça hızlı gelir, hemen başlar. Sonraki parça oynarken bir sonraki prefetch edilir.
        var nextPrefetch: Task<Data?, Never>?
        for (i, chunk) in chunks.enumerated() {
            if chunkedCancelled { nextPrefetch?.cancel(); break }
            let trimmed = chunk.trimmingCharacters(in: .whitespaces)
            guard !trimmed.isEmpty else { continue }

            let audioData: Data?
            if let pf = nextPrefetch, let data = await pf.value {
                audioData = data
            } else {
                do {
                    audioData = try await ChatService.shared.fetchTTSAudio(text: trimmed, voice: voice)
                } catch let e as APIError {
                    if case .limitExceeded(let info) = e {
                        await MainActor.run { limitExceededInfo = info }
                    }
                    return
                } catch {
                    audioData = nil
                }
            }
            nextPrefetch = i + 1 < chunks.count ? Task {
                let next = chunks[i + 1].trimmingCharacters(in: .whitespaces)
                guard !next.isEmpty else { return nil }
                return try? await ChatService.shared.fetchTTSAudio(text: next, voice: voice)
            } : nil

            guard let audioData, !chunkedCancelled else { break }
            guard let player = try? AVAudioPlayer(data: audioData) else { continue }

            await MainActor.run {
                audioPlayer = player
                player.delegate = self
                try? AVAudioSession.sharedInstance().setCategory(.playback, mode: .default)
                player.prepareToPlay()
                player.play()
                lastSpokenText = plain
                isSpeaking = true
                isLoading = false
            }

            await withCheckedContinuation { (cont: CheckedContinuation<Void, Never>) in
                playbackContinuation = cont
            }
        }

        await MainActor.run {
            audioPlayer = nil
            playbackContinuation = nil
            if !chunkedCancelled {
                isSpeaking = false
                lastSpokenText = nil
            }
        }
    }

    private func playSingle(_ plain: String, voice: String, language: String) async {
        do {
            let audioData = try await ChatService.shared.fetchTTSAudio(text: plain, voice: voice)
            if let player = try? AVAudioPlayer(data: audioData) {
                await MainActor.run {
                    audioPlayer = player
                    player.delegate = self
                    try? AVAudioSession.sharedInstance().setCategory(.playback, mode: .default)
                    player.prepareToPlay()
                    player.play()
                    lastSpokenText = plain
                    isSpeaking = true
                    isLoading = false
                }
                return
            }
        } catch let e as APIError {
            if case .limitExceeded(let info) = e {
                await MainActor.run { limitExceededInfo = info }
            }
            return
        } catch {}

        await MainActor.run {
            audioPlayer = nil
            let utterance = AVSpeechUtterance(string: plain)
            utterance.voice = AVSpeechSynthesisVoice(language: language)
            utterance.rate = AVSpeechUtteranceDefaultSpeechRate
            utterance.volume = 1.0
            synthesizer.speak(utterance)
            lastSpokenText = plain
            isSpeaking = true
        }
    }

    func stop() {
        chunkedCancelled = true
        playbackContinuation?.resume()
        playbackContinuation = nil
        if synthesizer.isSpeaking {
            synthesizer.stopSpeaking(at: .immediate)
        }
        audioPlayer?.stop()
        audioPlayer = nil
        isSpeaking = false
        lastSpokenText = nil
    }

    private var lastSpokenText: String?
    func toggleSpeak(_ text: String, language: String = "tr-TR", voice: String? = nil) {
        let plain = stripMarkdown(text)
        guard !plain.isEmpty else { return }
        if isSpeaking, lastSpokenText == plain {
            stop()
            lastSpokenText = nil
        } else {
            lastSpokenText = plain
            speak(text, language: language, voice: voice)
        }
    }

    /// Paragraflara böler (\n\n ile ayrılmış bloklar). Her paragraf ayrı API isteğine gider; cümle ortasından bölme olmaz.
    private func splitIntoParagraphs(_ text: String) -> [String] {
        let paragraphs = text
            .components(separatedBy: "\n\n")
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
        return paragraphs.isEmpty ? [text] : paragraphs
    }

    private func stripMarkdown(_ s: String) -> String {
        var out = s
        out = out.replacingOccurrences(of: "**", with: "")
        out = out.replacingOccurrences(of: "__", with: "")
        let headingPattern = try? NSRegularExpression(pattern: "#+\\s*", options: [])
        if let re = headingPattern {
            let range = NSRange(out.startIndex..., in: out)
            out = re.stringByReplacingMatches(in: out, options: [], range: range, withTemplate: "")
        }
        out = out.replacingOccurrences(of: "> ", with: "")
        out = out.replacingOccurrences(of: ">", with: "")
        out = out.replacingOccurrences(of: "`", with: "")
        out = out.replacingOccurrences(of: "- ", with: "")
        return out.trimmingCharacters(in: .whitespacesAndNewlines)
    }
}

extension SpeechService: AVSpeechSynthesizerDelegate {
    func speechSynthesizer(_ synthesizer: AVSpeechSynthesizer, didFinish utterance: AVSpeechUtterance) {
        DispatchQueue.main.async {
            self.isSpeaking = false
            self.lastSpokenText = nil
        }
    }

    func speechSynthesizer(_ synthesizer: AVSpeechSynthesizer, didCancel utterance: AVSpeechUtterance) {
        DispatchQueue.main.async {
            self.isSpeaking = false
            self.lastSpokenText = nil
        }
    }
}

extension SpeechService: AVAudioPlayerDelegate {
    func audioPlayerDidFinishPlaying(_ player: AVAudioPlayer, successfully flag: Bool) {
        DispatchQueue.main.async {
            self.playbackContinuation?.resume()
            self.playbackContinuation = nil
            if self.chunkedCancelled {
                self.audioPlayer = nil
                self.isSpeaking = false
                self.lastSpokenText = nil
            } else {
                self.audioPlayer = nil
            }
        }
    }
}

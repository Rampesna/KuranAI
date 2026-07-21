import AVFoundation
import Foundation

/// Ses kaydı yapar; m4a formatında döner.
@MainActor
final class AudioRecorderService: NSObject {
    static let shared = AudioRecorderService()

    private var recorder: AVAudioRecorder?
    private var recordingURL: URL?

    private override init() {
        super.init()
    }

    /// Mikrofon iznini ister ve hazırlık yapar.
    func requestPermissionAndPrepare() async throws {
        try await withCheckedThrowingContinuation { (continuation: CheckedContinuation<Void, Error>) in
            AVAudioApplication.requestRecordPermission { granted in
                if granted {
                    continuation.resume()
                } else {
                    continuation.resume(throwing: AudioRecorderError.permissionDenied)
                }
            }
        }
    }

    /// Kaydı başlatır.
    func startRecording() throws {
        let session = AVAudioSession.sharedInstance()
        try session.setCategory(.playAndRecord, mode: .default, options: [.defaultToSpeaker, .allowBluetoothHFP])
        try session.setActive(true)

        let tempDir = FileManager.default.temporaryDirectory
        let filename = "recording_\(UUID().uuidString).m4a"
        recordingURL = tempDir.appendingPathComponent(filename)

        guard let url = recordingURL else { throw AudioRecorderError.invalidURL }

        let settings: [String: Any] = [
            AVFormatIDKey: Int(kAudioFormatMPEG4AAC),
            AVSampleRateKey: 44100,
            AVNumberOfChannelsKey: 1,
            AVEncoderAudioQualityKey: AVAudioQuality.high.rawValue
        ]

        recorder = try AVAudioRecorder(url: url, settings: settings)
        recorder?.record()
    }

    /// Kaydı durdurur ve ses verisini döner.
    func stopRecording() throws -> Data {
        recorder?.stop()
        recorder = nil

        guard let url = recordingURL else { throw AudioRecorderError.noRecording }
        defer { recordingURL = nil }

        let data = try Data(contentsOf: url)
        try? FileManager.default.removeItem(at: url)
        return data
    }

    /// Şu an kayıt yapılıyor mu?
    var isRecording: Bool {
        recorder?.isRecording ?? false
    }
}

enum AudioRecorderError: LocalizedError {
    case permissionDenied
    case invalidURL
    case noRecording

    var errorDescription: String? {
        switch self {
        case .permissionDenied:
            return "Mikrofon izni verilmedi"
        case .invalidURL:
            return "Geçersiz kayıt yolu"
        case .noRecording:
            return "Kayıt bulunamadı"
        }
    }
}

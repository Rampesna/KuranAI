import AVFoundation
import Foundation

/// Ses kaydı yapar; m4a formatında döner.
@MainActor
final class AudioRecorderService: NSObject {
    static let shared = AudioRecorderService()

    private var recorder: AVAudioRecorder?
    private var recordingURL: URL?
    private var meterTimer: Timer?
    private var recordingStartedAt: Date?
    private var activityDetector = VoiceActivityDetector()
    private var autoStopHandler: ((AudioRecordingStopReason) -> Void)?

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
    func startRecording(onAutoStop: @escaping (AudioRecordingStopReason) -> Void) throws {
        cancelRecording()
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
        recorder?.isMeteringEnabled = true
        guard recorder?.record() == true else {
            recorder = nil
            recordingURL = nil
            throw AudioRecorderError.couldNotStart
        }

        activityDetector.reset()
        recordingStartedAt = Date()
        autoStopHandler = onAutoStop
        meterTimer = Timer.scheduledTimer(withTimeInterval: 0.1, repeats: true) { [weak self] _ in
            MainActor.assumeIsolated {
                self?.sampleAudioLevel()
            }
        }
    }

    /// Kaydı durdurur ve ses verisini döner.
    func stopRecording() throws -> Data {
        stopMetering()
        recorder?.stop()
        recorder = nil

        guard let url = recordingURL else { throw AudioRecorderError.noRecording }
        defer { recordingURL = nil }

        let data = try Data(contentsOf: url)
        try? FileManager.default.removeItem(at: url)
        return data
    }

    /// Kaydı yüklemeden iptal eder ve geçici dosyayı temizler.
    func cancelRecording() {
        stopMetering()
        recorder?.stop()
        recorder = nil
        if let recordingURL {
            try? FileManager.default.removeItem(at: recordingURL)
        }
        recordingURL = nil
    }

    /// Şu an kayıt yapılıyor mu?
    var isRecording: Bool {
        recorder?.isRecording ?? false
    }

    private func sampleAudioLevel() {
        guard let recorder, recorder.isRecording, let recordingStartedAt else { return }
        recorder.updateMeters()
        let elapsed = Date().timeIntervalSince(recordingStartedAt)
        guard let reason = activityDetector.observe(
            averagePowerDB: recorder.averagePower(forChannel: 0),
            elapsed: elapsed
        ) else { return }

        let handler = autoStopHandler
        stopMetering()
        handler?(reason)
    }

    private func stopMetering() {
        meterTimer?.invalidate()
        meterTimer = nil
        recordingStartedAt = nil
        autoStopHandler = nil
    }
}

enum AudioRecordingStopReason: Equatable {
    case silence
    case noSpeech
    case maximumDuration
}

struct VoiceActivityDetector {
    private let speechThresholdDB: Float = -38
    private let speechConfirmationSamples = 2
    private let warmupDuration: TimeInterval = 0.35
    private let silenceDuration: TimeInterval = 1.4
    private let noSpeechTimeout: TimeInterval = 8
    private let maximumDuration: TimeInterval = 60

    private var speechSamples = 0
    private var heardSpeech = false
    private var silenceStartedAt: TimeInterval?

    mutating func reset() {
        speechSamples = 0
        heardSpeech = false
        silenceStartedAt = nil
    }

    mutating func observe(averagePowerDB: Float, elapsed: TimeInterval) -> AudioRecordingStopReason? {
        if elapsed >= maximumDuration {
            return .maximumDuration
        }

        let hasVoice = elapsed >= warmupDuration && averagePowerDB >= speechThresholdDB
        if !heardSpeech {
            speechSamples = hasVoice ? speechSamples + 1 : 0
            heardSpeech = speechSamples >= speechConfirmationSamples
            if !heardSpeech && elapsed >= noSpeechTimeout {
                return .noSpeech
            }
            return nil
        }

        if hasVoice {
            silenceStartedAt = nil
            return nil
        }
        if silenceStartedAt == nil {
            silenceStartedAt = elapsed
        }
        if let silenceStartedAt, elapsed - silenceStartedAt >= silenceDuration {
            return .silence
        }
        return nil
    }
}

enum AudioRecorderError: LocalizedError {
    case permissionDenied
    case invalidURL
    case noRecording
    case couldNotStart

    var errorDescription: String? {
        switch self {
        case .permissionDenied:
            return "Mikrofon izni verilmedi"
        case .invalidURL:
            return "Geçersiz kayıt yolu"
        case .noRecording:
            return "Kayıt bulunamadı"
        case .couldNotStart:
            return "Ses kaydı başlatılamadı"
        }
    }
}

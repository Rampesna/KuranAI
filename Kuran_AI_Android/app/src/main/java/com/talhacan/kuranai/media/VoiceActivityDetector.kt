package com.talhacan.kuranai.media

import kotlin.math.log10

enum class RecordingStopReason {
    Silence,
    NoSpeech,
    MaximumDuration,
}

class VoiceActivityDetector(
    private val speechThresholdDb: Double = -38.0,
    private val speechConfirmationSamples: Int = 2,
    private val warmupMs: Long = 350,
    private val silenceDurationMs: Long = 1_400,
    private val noSpeechTimeoutMs: Long = 8_000,
    private val maximumDurationMs: Long = 60_000,
) {
    private var speechSamples = 0
    private var heardSpeech = false
    private var silenceStartedAtMs: Long? = null

    fun reset() {
        speechSamples = 0
        heardSpeech = false
        silenceStartedAtMs = null
    }

    fun observe(amplitude: Int, elapsedMs: Long): RecordingStopReason? {
        if (elapsedMs >= maximumDurationMs) {
            return RecordingStopReason.MaximumDuration
        }

        val hasVoice = elapsedMs >= warmupMs && amplitudeToDb(amplitude) >= speechThresholdDb
        if (!heardSpeech) {
            speechSamples = if (hasVoice) speechSamples + 1 else 0
            heardSpeech = speechSamples >= speechConfirmationSamples
            if (!heardSpeech && elapsedMs >= noSpeechTimeoutMs) {
                return RecordingStopReason.NoSpeech
            }
            return null
        }

        if (hasVoice) {
            silenceStartedAtMs = null
            return null
        }
        if (silenceStartedAtMs == null) {
            silenceStartedAtMs = elapsedMs
        }
        if (elapsedMs - checkNotNull(silenceStartedAtMs) >= silenceDurationMs) {
            return RecordingStopReason.Silence
        }
        return null
    }

    private fun amplitudeToDb(amplitude: Int): Double {
        if (amplitude <= 0) return -160.0
        val normalized = amplitude.coerceAtMost(MAX_AMPLITUDE).toDouble() / MAX_AMPLITUDE
        return 20.0 * log10(normalized)
    }

    private companion object {
        const val MAX_AMPLITUDE = 32_767
    }
}

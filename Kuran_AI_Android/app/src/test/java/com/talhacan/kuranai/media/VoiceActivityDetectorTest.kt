package com.talhacan.kuranai.media

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Test

class VoiceActivityDetectorTest {
    @Test
    fun stopsAfterSpeechFollowedBySilence() {
        val detector = VoiceActivityDetector()

        assertNull(detector.observe(amplitude = 20_000, elapsedMs = 400))
        assertNull(detector.observe(amplitude = 20_000, elapsedMs = 500))
        assertNull(detector.observe(amplitude = 0, elapsedMs = 700))
        assertNull(detector.observe(amplitude = 0, elapsedMs = 2_000))
        assertEquals(
            RecordingStopReason.Silence,
            detector.observe(amplitude = 0, elapsedMs = 2_100),
        )
    }

    @Test
    fun stopsWhenNoSpeechIsDetected() {
        val detector = VoiceActivityDetector()

        assertNull(detector.observe(amplitude = 0, elapsedMs = 7_900))
        assertEquals(
            RecordingStopReason.NoSpeech,
            detector.observe(amplitude = 0, elapsedMs = 8_000),
        )
    }

    @Test
    fun stopsAtMaximumDurationEvenWhileSpeaking() {
        val detector = VoiceActivityDetector()

        assertEquals(
            RecordingStopReason.MaximumDuration,
            detector.observe(amplitude = 20_000, elapsedMs = 60_000),
        )
    }

    @Test
    fun briefNoiseDoesNotCountAsSpeech() {
        val detector = VoiceActivityDetector()

        assertNull(detector.observe(amplitude = 20_000, elapsedMs = 500))
        assertNull(detector.observe(amplitude = 0, elapsedMs = 600))
        assertEquals(
            RecordingStopReason.NoSpeech,
            detector.observe(amplitude = 0, elapsedMs = 8_000),
        )
    }
}

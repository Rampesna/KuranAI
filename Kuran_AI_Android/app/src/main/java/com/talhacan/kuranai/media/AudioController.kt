package com.talhacan.kuranai.media

import android.content.Context
import android.media.MediaPlayer
import android.media.MediaRecorder
import android.os.Build
import android.os.Handler
import android.os.Looper
import android.os.SystemClock
import java.io.File

class AudioController(private val context: Context) {
    private var recorder: MediaRecorder? = null
    private var recordingFile: File? = null
    private var player: MediaPlayer? = null
    private val meterHandler = Handler(Looper.getMainLooper())
    private val activityDetector = VoiceActivityDetector()
    private var recordingStartedAtMs = 0L
    private var autoStopHandler: ((RecordingStopReason) -> Unit)? = null
    private val meterRunnable = object : Runnable {
        override fun run() {
            val activeRecorder = recorder ?: return
            val amplitude = runCatching { activeRecorder.maxAmplitude }.getOrDefault(0)
            val reason = activityDetector.observe(
                amplitude = amplitude,
                elapsedMs = SystemClock.elapsedRealtime() - recordingStartedAtMs,
            )
            if (reason != null) {
                val callback = autoStopHandler
                stopMonitoring()
                callback?.invoke(reason)
            } else {
                meterHandler.postDelayed(this, METER_INTERVAL_MS)
            }
        }
    }

    fun startRecording(onAutoStop: (RecordingStopReason) -> Unit): File {
        stopPlayback()
        cancelRecording()
        val target = File(context.cacheDir, "question-${System.currentTimeMillis()}.m4a")
        @Suppress("DEPRECATION")
        val nextRecorder = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            MediaRecorder(context)
        } else {
            MediaRecorder()
        }
        try {
            nextRecorder.apply {
                setAudioSource(MediaRecorder.AudioSource.MIC)
                setOutputFormat(MediaRecorder.OutputFormat.MPEG_4)
                setAudioEncoder(MediaRecorder.AudioEncoder.AAC)
                setAudioEncodingBitRate(128_000)
                setAudioSamplingRate(44_100)
                setOutputFile(target.absolutePath)
                prepare()
                start()
            }
        } catch (error: Exception) {
            nextRecorder.release()
            target.delete()
            throw error
        }
        recorder = nextRecorder
        recordingFile = target
        activityDetector.reset()
        recordingStartedAtMs = SystemClock.elapsedRealtime()
        autoStopHandler = onAutoStop
        meterHandler.postDelayed(meterRunnable, METER_INTERVAL_MS)
        return target
    }

    fun stopRecording(): File {
        val target = recordingFile ?: error("Recording has not started")
        stopMonitoring()
        recorder?.runCatching { stop() }
        recorder?.release()
        recorder = null
        recordingFile = null
        return target
    }

    fun cancelRecording() {
        stopMonitoring()
        recorder?.runCatching { stop() }
        recorder?.release()
        recorder = null
        recordingFile?.delete()
        recordingFile = null
    }

    fun playMp3(bytes: ByteArray) {
        stopPlayback()
        val target = File(context.cacheDir, "answer-${System.currentTimeMillis()}.mp3")
        target.writeBytes(bytes)
        player = MediaPlayer().apply {
            setDataSource(target.absolutePath)
            setOnCompletionListener {
                it.release()
                if (player === it) player = null
                target.delete()
            }
            prepare()
            start()
        }
    }

    fun stopPlayback() {
        player?.runCatching { stop() }
        player?.release()
        player = null
    }

    fun close() {
        cancelRecording()
        stopPlayback()
    }

    private fun stopMonitoring() {
        meterHandler.removeCallbacks(meterRunnable)
        autoStopHandler = null
        recordingStartedAtMs = 0L
    }

    private companion object {
        const val METER_INTERVAL_MS = 100L
    }
}

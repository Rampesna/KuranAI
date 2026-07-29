package com.talhacan.kuranai.media

import android.content.Context
import android.media.MediaPlayer
import android.media.MediaRecorder
import android.os.Build
import java.io.File

class AudioController(private val context: Context) {
    private var recorder: MediaRecorder? = null
    private var recordingFile: File? = null
    private var player: MediaPlayer? = null

    fun startRecording(): File {
        stopPlayback()
        val target = File(context.cacheDir, "question-${System.currentTimeMillis()}.m4a")
        @Suppress("DEPRECATION")
        val nextRecorder = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            MediaRecorder(context)
        } else {
            MediaRecorder()
        }
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
        recorder = nextRecorder
        recordingFile = target
        return target
    }

    fun stopRecording(): File {
        val target = recordingFile ?: error("Recording has not started")
        recorder?.runCatching { stop() }
        recorder?.release()
        recorder = null
        recordingFile = null
        return target
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
}

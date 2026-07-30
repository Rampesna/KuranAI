package com.talhacan.kuranai.review

import android.app.Activity
import android.content.Context
import com.google.android.play.core.review.ReviewManagerFactory
import com.talhacan.kuranai.BuildConfig
import java.util.concurrent.TimeUnit

/**
 * Requests a Play review only after the user has received meaningful value.
 *
 * Google Play can silently suppress the dialog because of its own quota. The
 * app therefore never promises that a dialog will appear and never interrupts
 * the normal chat flow when the request fails.
 */
class InAppReviewCoordinator(context: Context) {
    private val appContext = context.applicationContext
    private val preferences = appContext.getSharedPreferences(
        "kuranai_review_prompt",
        Context.MODE_PRIVATE,
    )
    private val manager = ReviewManagerFactory.create(appContext)

    fun recordSuccessfulAnswer(activity: Activity) {
        val now = System.currentTimeMillis()
        val firstValueAt = preferences.getLong(KEY_FIRST_VALUE_AT, 0L).let { existing ->
            if (existing > 0L) {
                existing
            } else {
                preferences.edit().putLong(KEY_FIRST_VALUE_AT, now).apply()
                now
            }
        }
        val successfulAnswers = preferences.getInt(KEY_SUCCESSFUL_ANSWERS, 0) + 1
        preferences.edit().putInt(KEY_SUCCESSFUL_ANSWERS, successfulAnswers).apply()

        val lastPromptAt = preferences.getLong(KEY_LAST_PROMPT_AT, 0L)
        val lastPromptVersion = preferences.getString(KEY_LAST_PROMPT_VERSION, null)
        val enoughExperience = successfulAnswers >= REQUIRED_SUCCESSFUL_ANSWERS &&
            now - firstValueAt >= MINIMUM_EXPERIENCE_MILLIS
        val outsideCooldown = lastPromptAt == 0L || now - lastPromptAt >= PROMPT_COOLDOWN_MILLIS
        val newVersion = lastPromptVersion != BuildConfig.VERSION_NAME

        if (!enoughExperience || !outsideCooldown || !newVersion) return

        manager.requestReviewFlow().addOnCompleteListener { request ->
            if (!request.isSuccessful || activity.isFinishing || activity.isDestroyed) {
                return@addOnCompleteListener
            }

            preferences.edit()
                .putLong(KEY_LAST_PROMPT_AT, System.currentTimeMillis())
                .putString(KEY_LAST_PROMPT_VERSION, BuildConfig.VERSION_NAME)
                .apply()

            manager.launchReviewFlow(activity, request.result)
        }
    }

    private companion object {
        const val KEY_FIRST_VALUE_AT = "first_value_at"
        const val KEY_SUCCESSFUL_ANSWERS = "successful_answers"
        const val KEY_LAST_PROMPT_AT = "last_prompt_at"
        const val KEY_LAST_PROMPT_VERSION = "last_prompt_version"
        const val REQUIRED_SUCCESSFUL_ANSWERS = 3
        val MINIMUM_EXPERIENCE_MILLIS = TimeUnit.DAYS.toMillis(2)
        val PROMPT_COOLDOWN_MILLIS = TimeUnit.DAYS.toMillis(120)
    }
}

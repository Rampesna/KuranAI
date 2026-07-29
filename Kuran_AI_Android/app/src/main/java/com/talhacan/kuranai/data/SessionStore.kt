package com.talhacan.kuranai.data

import android.content.Context
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.AuthTokens

@Suppress("DEPRECATION")
class SessionStore(context: Context) {
    private val masterKey = MasterKey.Builder(context)
        .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
        .build()

    private val preferences = EncryptedSharedPreferences.create(
        context,
        "kuranai_secure_preferences",
        masterKey,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM,
    )

    fun tokens(): AuthTokens? {
        val accessToken = preferences.getString(KEY_ACCESS_TOKEN, null) ?: return null
        val refreshToken = preferences.getString(KEY_REFRESH_TOKEN, null) ?: return null
        return AuthTokens(accessToken, refreshToken)
    }

    fun saveTokens(tokens: AuthTokens) {
        preferences.edit()
            .putString(KEY_ACCESS_TOKEN, tokens.accessToken)
            .putString(KEY_REFRESH_TOKEN, tokens.refreshToken)
            .apply()
    }

    fun clearTokens() {
        preferences.edit()
            .remove(KEY_ACCESS_TOKEN)
            .remove(KEY_REFRESH_TOKEN)
            .apply()
    }

    var language: AppLanguage
        get() = AppLanguage.from(preferences.getString(KEY_LANGUAGE, null))
        set(value) {
            preferences.edit().putString(KEY_LANGUAGE, value.code).apply()
        }

    var darkTheme: Boolean
        get() = preferences.getBoolean(KEY_DARK_THEME, true)
        set(value) {
            preferences.edit().putBoolean(KEY_DARK_THEME, value).apply()
        }

    var aiConsent: Boolean
        get() = preferences.getBoolean(KEY_AI_CONSENT, false)
        set(value) {
            preferences.edit().putBoolean(KEY_AI_CONSENT, value).apply()
        }

    companion object {
        private const val KEY_ACCESS_TOKEN = "access_token"
        private const val KEY_REFRESH_TOKEN = "refresh_token"
        private const val KEY_LANGUAGE = "language"
        private const val KEY_DARK_THEME = "dark_theme"
        private const val KEY_AI_CONSENT = "ai_consent"
    }
}

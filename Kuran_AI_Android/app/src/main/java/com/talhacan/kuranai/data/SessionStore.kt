package com.talhacan.kuranai.data

import android.content.Context
import android.content.SharedPreferences
import android.util.Log
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.AuthTokens
import java.security.KeyStore

@Suppress("DEPRECATION")
class SessionStore(context: Context) {
    private data class Backend(
        val preferences: SharedPreferences,
        val isSecure: Boolean,
    )

    private val appContext = context.applicationContext
    private val backendLock = Any()

    @Volatile
    private var backend = createBackendWithRecovery()

    @Volatile
    private var inMemoryTokens: AuthTokens? = null

    fun tokens(): AuthTokens? {
        val storedTokens = read(null) { current ->
            if (!current.isSecure) {
                return@read inMemoryTokens
            }
            val accessToken =
                current.preferences.getString(KEY_ACCESS_TOKEN, null) ?: return@read null
            val refreshToken =
                current.preferences.getString(KEY_REFRESH_TOKEN, null) ?: return@read null
            AuthTokens(accessToken, refreshToken)
        }
        if (storedTokens != null) {
            inMemoryTokens = storedTokens
        }
        return storedTokens ?: inMemoryTokens
    }

    fun saveTokens(tokens: AuthTokens) {
        inMemoryTokens = tokens
        write { current ->
            if (current.isSecure) {
                current.preferences.edit()
                    .putString(KEY_ACCESS_TOKEN, tokens.accessToken)
                    .putString(KEY_REFRESH_TOKEN, tokens.refreshToken)
                    .apply()
            }
        }
    }

    fun clearTokens() {
        inMemoryTokens = null
        write { current ->
            if (current.isSecure) {
                current.preferences.edit()
                    .remove(KEY_ACCESS_TOKEN)
                    .remove(KEY_REFRESH_TOKEN)
                    .apply()
            }
        }
    }

    var language: AppLanguage
        get() = AppLanguage.from(
            read(null) { current ->
                current.preferences.getString(KEY_LANGUAGE, null)
            },
        )
        set(value) {
            write { current ->
                current.preferences.edit().putString(KEY_LANGUAGE, value.code).apply()
            }
        }

    var darkTheme: Boolean
        get() = read(true) { current ->
            current.preferences.getBoolean(KEY_DARK_THEME, true)
        }
        set(value) {
            write { current ->
                current.preferences.edit().putBoolean(KEY_DARK_THEME, value).apply()
            }
        }

    var aiConsent: Boolean
        get() = read(false) { current ->
            current.preferences.getBoolean(KEY_AI_CONSENT, false)
        }
        set(value) {
            write { current ->
                current.preferences.edit().putBoolean(KEY_AI_CONSENT, value).apply()
            }
        }

    private fun createBackendWithRecovery(): Backend {
        return try {
            createSecureBackend()
        } catch (firstError: Exception) {
            Log.e(
                TAG,
                "Encrypted preferences could not be opened; resetting unreadable storage.",
                firstError,
            )
            resetEncryptedStorage()
            try {
                createSecureBackend()
            } catch (secondError: Exception) {
                Log.e(
                    TAG,
                    "Encrypted preferences are unavailable; using a token-free fallback.",
                    secondError,
                )
                Backend(
                    preferences = appContext.getSharedPreferences(
                        RECOVERY_PREFERENCES_NAME,
                        Context.MODE_PRIVATE,
                    ),
                    isSecure = false,
                )
            }
        }
    }

    private fun createSecureBackend(): Backend {
        val masterKey = MasterKey.Builder(appContext)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build()
        val preferences = EncryptedSharedPreferences.create(
            appContext,
            SECURE_PREFERENCES_NAME,
            masterKey,
            EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
            EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM,
        )

        // EncryptedSharedPreferences can open successfully and fail only when the first
        // encrypted value is read. Validate every entry before application startup continues.
        preferences.all
        return Backend(preferences = preferences, isSecure = true)
    }

    private fun recoverBackend(failedBackend: Backend, error: Exception): Backend {
        synchronized(backendLock) {
            if (backend !== failedBackend) {
                return backend
            }
            Log.e(TAG, "Session storage became unreadable; recreating it.", error)
            resetEncryptedStorage()
            backend = createBackendWithRecovery()
            return backend
        }
    }

    private fun resetEncryptedStorage() {
        appContext.deleteSharedPreferences(SECURE_PREFERENCES_NAME)
        try {
            val keyStore = KeyStore.getInstance(ANDROID_KEYSTORE).apply {
                load(null)
            }
            if (keyStore.containsAlias(MasterKey.DEFAULT_MASTER_KEY_ALIAS)) {
                keyStore.deleteEntry(MasterKey.DEFAULT_MASTER_KEY_ALIAS)
            }
        } catch (error: Exception) {
            Log.w(TAG, "The obsolete Android Keystore entry could not be removed.", error)
        }
    }

    private fun <T> read(fallback: T, operation: (Backend) -> T): T {
        val current = backend
        return try {
            operation(current)
        } catch (firstError: Exception) {
            val recovered = recoverBackend(current, firstError)
            try {
                operation(recovered)
            } catch (secondError: Exception) {
                Log.e(TAG, "Session storage read failed after recovery.", secondError)
                fallback
            }
        }
    }

    private fun write(operation: (Backend) -> Unit) {
        val current = backend
        try {
            operation(current)
        } catch (firstError: Exception) {
            val recovered = recoverBackend(current, firstError)
            try {
                operation(recovered)
            } catch (secondError: Exception) {
                Log.e(TAG, "Session storage write failed after recovery.", secondError)
            }
        }
    }

    companion object {
        private const val TAG = "SessionStore"
        private const val ANDROID_KEYSTORE = "AndroidKeyStore"
        private const val SECURE_PREFERENCES_NAME = "kuranai_secure_preferences"
        private const val RECOVERY_PREFERENCES_NAME = "kuranai_recovery_preferences"
        private const val KEY_ACCESS_TOKEN = "access_token"
        private const val KEY_REFRESH_TOKEN = "refresh_token"
        private const val KEY_LANGUAGE = "language"
        private const val KEY_DARK_THEME = "dark_theme"
        private const val KEY_AI_CONSENT = "ai_consent"
    }
}

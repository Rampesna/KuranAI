package com.talhacan.kuranai

import android.Manifest
import android.app.Activity
import android.app.Application
import android.content.pm.PackageManager
import android.location.Location
import android.location.LocationManager
import android.net.Uri
import androidx.browser.customtabs.CustomTabsIntent
import androidx.core.content.ContextCompat
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.talhacan.kuranai.billing.BillingStore
import com.talhacan.kuranai.data.ApiException
import com.talhacan.kuranai.media.AudioController
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.AppPage
import com.talhacan.kuranai.model.AppState
import com.talhacan.kuranai.model.ChatMessage
import com.talhacan.kuranai.model.DailyVerse
import com.talhacan.kuranai.model.SourceType
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.collectLatest
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.File
import java.util.Calendar
import java.util.UUID

class KuranAIViewModel(application: Application) : AndroidViewModel(application) {
    private val container = application as KuranAIApplication
    private val sessionStore = container.sessionStore
    private val api = container.apiClient
    private val identity = container.identityManager
    private val billing = container.billingStore
    private val audio = AudioController(application)
    private val _state = MutableStateFlow(
        AppState(
            language = sessionStore.language,
            darkTheme = sessionStore.darkTheme,
        ),
    )
    val state: StateFlow<AppState> = _state.asStateFlow()

    val billingStore: BillingStore
        get() = billing

    val hasAiConsent: Boolean
        get() = sessionStore.aiConsent

    init {
        billing.setLanguage(sessionStore.language)
        viewModelScope.launch {
            billing.products.collectLatest { products ->
                _state.value = _state.value.copy(storeProducts = products)
            }
        }
        viewModelScope.launch {
            billing.message.collectLatest { message ->
                if (!message.isNullOrBlank()) {
                    _state.value = _state.value.copy(notice = message)
                }
            }
        }
        viewModelScope.launch {
            container.entitlementEvents.collectLatest { refreshSubscription() }
        }
        restoreSession()
    }

    fun clearStatus() {
        billing.clearMessage()
        _state.value = _state.value.copy(notice = null, error = null)
    }

    fun setLanguage(language: AppLanguage) {
        sessionStore.language = language
        billing.setLanguage(language)
        _state.value = _state.value.copy(language = language)
        viewModelScope.launch {
            loadVerse()
            loadPrayerTimes()
        }
    }

    fun setDarkTheme(enabled: Boolean) {
        sessionStore.darkTheme = enabled
        _state.value = _state.value.copy(darkTheme = enabled)
    }

    fun setPage(page: AppPage) {
        _state.value = _state.value.copy(page = page)
        when (page) {
            AppPage.Plans -> viewModelScope.launch { loadPackagesAndSubscription() }
            AppPage.Profile -> viewModelScope.launch { loadProfileDetails() }
            AppPage.Chat -> viewModelScope.launch { loadSessions() }
            AppPage.Home -> Unit
        }
    }

    fun setSource(source: SourceType) {
        _state.value = _state.value.copy(
            source = source,
            messages = emptyList(),
            currentSessionId = null,
        )
    }

    fun newChat(initialQuestion: String? = null) {
        _state.value = _state.value.copy(
            page = AppPage.Chat,
            messages = emptyList(),
            currentSessionId = null,
        )
        if (!initialQuestion.isNullOrBlank()) sendMessage(initialQuestion)
    }

    fun selectSession(sessionId: Long) {
        viewModelScope.launch {
            runTask {
                val messages = api.messages(sessionId)
                _state.value = _state.value.copy(
                    page = AppPage.Chat,
                    currentSessionId = sessionId,
                    messages = messages,
                )
            }
        }
    }

    fun grantAiConsent() {
        sessionStore.aiConsent = true
    }

    fun loginEmail(credential: String, password: String) {
        viewModelScope.launch {
            runTask {
                api.login(credential.trim(), password)
                onAuthenticated()
            }
        }
    }

    fun register(name: String, username: String, email: String, password: String) {
        viewModelScope.launch {
            runTask {
                api.register(name.trim(), username.trim(), email.trim(), password)
                api.login(email.trim(), password)
                onAuthenticated()
            }
        }
    }

    fun signInGoogle(activity: Activity) {
        viewModelScope.launch {
            runTask {
                identity.signInWithGoogle(activity, _state.value.language)
                onAuthenticated()
            }
        }
    }

    fun signInPasskey(activity: Activity) {
        viewModelScope.launch {
            runTask {
                identity.signInWithPasskey(activity, _state.value.language)
                onAuthenticated()
            }
        }
    }

    fun openAppleSignIn(activity: Activity) {
        val signInUrl = Uri.parse(BuildConfig.APPLE_SIGN_IN_URL)
            .buildUpon()
            .appendQueryParameter("lang", _state.value.language.code)
            .build()
        CustomTabsIntent.Builder()
            .setShowTitle(true)
            .build()
            .launchUrl(activity, signInUrl)
    }

    fun completeAppleSignIn(exchangeCode: String) {
        viewModelScope.launch {
            runTask {
                api.finishMobileAuth(exchangeCode)
                onAuthenticated()
            }
        }
    }

    fun logout(activity: Activity) {
        viewModelScope.launch {
            runCatching { api.logout() }
            identity.clear(activity)
            billing.attachUser(null)
            _state.value = AppState(
                language = sessionStore.language,
                darkTheme = sessionStore.darkTheme,
                restoringSession = false,
            )
        }
    }

    fun createPasskey(activity: Activity) {
        viewModelScope.launch {
            runTask {
                identity.createPasskey(activity, _state.value.language)
                _state.value = _state.value.copy(
                    passkeyCount = api.passkeyCount(),
                    notice = localized(
                        "Geçiş anahtarınız hazır.",
                        "Your passkey is ready.",
                        "مفتاح المرور جاهز.",
                    ),
                )
            }
        }
    }

    fun deletePasskeys() {
        viewModelScope.launch {
            runTask {
                api.deletePasskeys()
                _state.value = _state.value.copy(
                    passkeyCount = 0,
                    notice = localized(
                        "Geçiş anahtarları kaldırıldı.",
                        "Passkeys removed.",
                        "تمت إزالة مفاتيح المرور.",
                    ),
                )
            }
        }
    }

    fun updateUsername(username: String) {
        viewModelScope.launch {
            runTask {
                _state.value = _state.value.copy(
                    user = api.updateUsername(username.trim()),
                    notice = localized("Kaydedildi.", "Saved.", "تم الحفظ."),
                )
            }
        }
    }

    fun updatePassword(currentPassword: String, newPassword: String) {
        viewModelScope.launch {
            runTask {
                api.updatePassword(currentPassword, newPassword)
                _state.value = _state.value.copy(
                    notice = localized(
                        "Şifre güncellendi.",
                        "Password updated.",
                        "تم تحديث كلمة المرور.",
                    ),
                )
            }
        }
    }

    fun uploadAvatar(source: Uri) {
        viewModelScope.launch {
            runTask {
                val temp = withContext(Dispatchers.IO) {
                    val target = File(getApplication<Application>().cacheDir, "avatar-upload.jpg")
                    getApplication<Application>().contentResolver.openInputStream(source).use { input ->
                        requireNotNull(input) { "Photo could not be read" }
                        target.outputStream().use { output ->
                            input.copyTo(output)
                        }
                    }
                    target
                }
                _state.value = _state.value.copy(
                    user = api.uploadAvatar(temp),
                    notice = localized(
                        "Fotoğraf güncellendi.",
                        "Photo updated.",
                        "تم تحديث الصورة.",
                    ),
                )
                temp.delete()
            }
        }
    }

    fun deleteAvatar() {
        viewModelScope.launch {
            runTask {
                val user = api.deleteAvatar() ?: api.me()
                _state.value = _state.value.copy(user = user)
            }
        }
    }

    fun deleteAccount() {
        viewModelScope.launch {
            runTask {
                api.deleteAccount()
                billing.attachUser(null)
                _state.value = AppState(
                    language = sessionStore.language,
                    darkTheme = sessionStore.darkTheme,
                    restoringSession = false,
                )
            }
        }
    }

    fun sendMessage(content: String) {
        val trimmed = content.trim()
        if (trimmed.isBlank() || _state.value.sending) return
        if (!sessionStore.aiConsent) {
            _state.value = _state.value.copy(
                error = localized(
                    "Devam etmek için AI veri paylaşımına izin verin.",
                    "Allow AI data sharing to continue.",
                    "اسمح بمشاركة البيانات مع الذكاء الاصطناعي للمتابعة.",
                ),
            )
            return
        }

        val assistantId = UUID.randomUUID().toString()
        _state.value = _state.value.copy(
            sending = true,
            error = null,
            messages = _state.value.messages +
                ChatMessage(UUID.randomUUID().toString(), "user", trimmed) +
                ChatMessage(assistantId, "assistant", "", streaming = true),
        )
        viewModelScope.launch {
            try {
                val resolvedSession = api.streamMessage(
                    _state.value.currentSessionId,
                    _state.value.source,
                    trimmed,
                    onStarted = { id ->
                        _state.value = _state.value.copy(currentSessionId = id)
                    },
                    onChunk = { chunk ->
                        _state.value = _state.value.copy(
                            messages = _state.value.messages.map {
                                if (it.id == assistantId) it.copy(content = it.content + chunk) else it
                            },
                        )
                    },
                )
                _state.value = _state.value.copy(
                    sending = false,
                    currentSessionId = resolvedSession,
                    messages = _state.value.messages.map {
                        if (it.id == assistantId) it.copy(streaming = false) else it
                    },
                )
                loadSessions()
                refreshSubscription()
            } catch (error: Exception) {
                _state.value = _state.value.copy(
                    sending = false,
                    messages = _state.value.messages.filterNot {
                        it.id == assistantId && it.content.isBlank()
                    },
                    error = friendlyError(error),
                )
            }
        }
    }

    fun toggleRecording() {
        if (!sessionStore.aiConsent) {
            _state.value = _state.value.copy(
                error = localized(
                    "Sesli soru için AI veri paylaşımına izin verin.",
                    "Allow AI data sharing for voice questions.",
                    "اسمح بمشاركة البيانات للأسئلة الصوتية.",
                ),
            )
            return
        }
        if (_state.value.recording) {
            val file = runCatching(audio::stopRecording).getOrElse {
                _state.value = _state.value.copy(recording = false, error = it.localizedMessage)
                return
            }
            _state.value = _state.value.copy(recording = false, loading = true)
            viewModelScope.launch {
                try {
                    val transcription = api.transcribe(file)
                    sendMessage(transcription)
                } catch (error: Exception) {
                    _state.value = _state.value.copy(error = friendlyError(error))
                } finally {
                    file.delete()
                    _state.value = _state.value.copy(loading = false)
                }
            }
        } else {
            runCatching { audio.startRecording() }
                .onSuccess { _state.value = _state.value.copy(recording = true, error = null) }
                .onFailure { _state.value = _state.value.copy(error = it.localizedMessage) }
        }
    }

    fun speak(text: String, voice: String = "ash") {
        viewModelScope.launch {
            runTask {
                audio.playMp3(api.tts(text, voice))
                refreshSubscription()
            }
        }
    }

    fun refreshHome() {
        viewModelScope.launch {
            loadVerse()
            loadPrayerTimes()
        }
    }

    private fun restoreSession() {
        viewModelScope.launch {
            if (sessionStore.tokens() == null) {
                _state.value = _state.value.copy(restoringSession = false)
                loadVerse()
                loadPrayerTimes()
                return@launch
            }
            runCatching { onAuthenticated() }
                .onFailure {
                    sessionStore.clearTokens()
                    _state.value = _state.value.copy(
                        authenticated = false,
                        restoringSession = false,
                    )
                    loadVerse()
                    loadPrayerTimes()
                }
        }
    }

    private suspend fun onAuthenticated() {
        val user = api.me()
        _state.value = _state.value.copy(
            authenticated = true,
            restoringSession = false,
            user = user,
            error = null,
        )
        billing.attachUser(user.id)
        loadDashboard()
    }

    private suspend fun loadDashboard() {
        val verse = viewModelScope.async { api.verse(_state.value.language) }
        val packages = viewModelScope.async { api.packages() }
        val overview = viewModelScope.async { api.subscriptionOverview() }
        val sessions = viewModelScope.async { api.sessions() }
        val passkeys = viewModelScope.async { api.passkeyCount() }
        _state.value = _state.value.copy(
            dailyVerse = verse.await() ?: fallbackVerse(),
            packages = packages.await(),
            subscription = overview.await(),
            sessions = sessions.await(),
            passkeyCount = passkeys.await(),
        )
        loadPrayerTimes()
    }

    private suspend fun loadProfileDetails() {
        runCatching {
            _state.value = _state.value.copy(
                user = api.me(),
                passkeyCount = api.passkeyCount(),
            )
        }
    }

    private suspend fun loadPackagesAndSubscription() {
        runCatching {
            _state.value = _state.value.copy(
                packages = api.packages(),
                subscription = api.subscriptionOverview(),
            )
        }.onFailure {
            _state.value = _state.value.copy(error = friendlyError(it))
        }
    }

    private suspend fun refreshSubscription() {
        runCatching {
            _state.value = _state.value.copy(subscription = api.subscriptionOverview())
        }
    }

    private suspend fun loadSessions() {
        runCatching {
            _state.value = _state.value.copy(sessions = api.sessions())
        }
    }

    private suspend fun loadVerse() {
        _state.value = _state.value.copy(
            dailyVerse = api.verse(_state.value.language) ?: fallbackVerse(),
        )
    }

    private suspend fun loadPrayerTimes() {
        val location = lastKnownLocation()
        runCatching {
            api.prayerTimes(location?.latitude, location?.longitude, _state.value.language)
        }.onSuccess { (place, times) ->
            _state.value = _state.value.copy(prayerPlace = place, prayerTimes = times)
        }
    }

    private fun lastKnownLocation(): Location? {
        val application = getApplication<Application>()
        if (
            ContextCompat.checkSelfPermission(application, Manifest.permission.ACCESS_FINE_LOCATION) !=
            PackageManager.PERMISSION_GRANTED &&
            ContextCompat.checkSelfPermission(application, Manifest.permission.ACCESS_COARSE_LOCATION) !=
            PackageManager.PERMISSION_GRANTED
        ) {
            return null
        }
        val manager = application.getSystemService(LocationManager::class.java)
        return manager.getProviders(true)
            .mapNotNull { provider -> runCatching { manager.getLastKnownLocation(provider) }.getOrNull() }
            .maxByOrNull { it.time }
    }

    private fun fallbackVerse(): DailyVerse {
        val verses = listOf(
            DailyVerse(
                "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
                localized(
                    "Bismillâhirrahmânirrahîm",
                    "In the name of Allah, the Entirely Merciful, the Especially Merciful.",
                    "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
                ),
                localized("Fatiha 1", "Al-Fatiha 1", "الفاتحة ١"),
            ),
            DailyVerse(
                "قُلْ هُوَ اللَّهُ أَحَدٌ",
                localized(
                    "De ki: O Allah birdir.",
                    "Say: He is Allah, the One.",
                    "قُلْ هُوَ اللَّهُ أَحَدٌ.",
                ),
                localized("İhlas 1", "Al-Ikhlas 1", "الإخلاص ١"),
            ),
        )
        val day = Calendar.getInstance().get(Calendar.DAY_OF_YEAR)
        return verses[day % verses.size]
    }

    private suspend fun runTask(block: suspend () -> Unit) {
        _state.value = _state.value.copy(loading = true, error = null)
        try {
            block()
        } catch (error: Exception) {
            _state.value = _state.value.copy(error = friendlyError(error))
        } finally {
            _state.value = _state.value.copy(loading = false)
        }
    }

    private fun friendlyError(error: Throwable): String =
        if (error is ApiException) error.message else error.localizedMessage ?: "Unknown error"

    private fun localized(tr: String, en: String, ar: String): String = when (_state.value.language) {
        AppLanguage.Turkish -> tr
        AppLanguage.English -> en
        AppLanguage.Arabic -> ar
    }
}

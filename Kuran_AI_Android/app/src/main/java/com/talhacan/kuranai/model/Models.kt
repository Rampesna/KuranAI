package com.talhacan.kuranai.model

enum class AppLanguage(val code: String, val displayName: String) {
    Turkish("tr", "Türkçe"),
    English("en", "English"),
    Arabic("ar", "العربية");

    companion object {
        fun from(code: String?): AppLanguage = entries.firstOrNull { it.code == code } ?: Turkish
    }
}

enum class AppPage {
    Home,
    Chat,
    Profile,
    Plans,
}

enum class SourceType(val apiValue: String, val mark: String) {
    Quran("meal", "ق"),
    Tafsir("tefsir", "ت"),
    Hadith("hadis", "ح"),
}

data class AuthTokens(
    val accessToken: String,
    val refreshToken: String,
)

data class User(
    val id: Long,
    val email: String,
    val name: String,
    val username: String,
    val avatarUrl: String?,
    val hasPassword: Boolean,
)

data class ChatMessage(
    val id: String,
    val role: String,
    val content: String,
    val streaming: Boolean = false,
)

data class ChatSession(
    val id: Long,
    val firstMessage: String,
    val updatedAt: String,
)

data class DailyVerse(
    val arabicText: String,
    val translation: String,
    val reference: String,
)

data class PrayerTime(
    val name: String,
    val time: String,
)

data class SubscriptionPackage(
    val id: Long,
    val name: String,
    val nameEn: String,
    val slug: String,
    val dailyQuestionLimit: Int?,
    val dailyTtsLimit: Int?,
    val monthlyPriceTl: Double,
    val yearlyPriceTl: Double,
) {
    fun localizedName(language: AppLanguage): String = when (language) {
        AppLanguage.Turkish -> name
        AppLanguage.English -> nameEn
        AppLanguage.Arabic -> when (slug.lowercase()) {
            "bronze" -> "برونزية"
            "silver" -> "فضية"
            "gold" -> "ذهبية"
            "diamond" -> "ماسية"
            else -> nameEn
        }
    }
}

data class SubscriptionOverview(
    val hasPackage: Boolean = false,
    val currentPackage: SubscriptionPackage? = null,
    val expiresAt: String? = null,
    val messageUsed: Int = 0,
    val messageLimit: Int = 3,
    val ttsUsed: Int = 0,
    val ttsLimit: Int = 0,
)

data class StoreProduct(
    val productId: String,
    val formattedPrice: String,
)

data class AppState(
    val language: AppLanguage = AppLanguage.Turkish,
    val darkTheme: Boolean = true,
    val authenticated: Boolean = false,
    val restoringSession: Boolean = true,
    val page: AppPage = AppPage.Home,
    val user: User? = null,
    val source: SourceType = SourceType.Quran,
    val messages: List<ChatMessage> = emptyList(),
    val sessions: List<ChatSession> = emptyList(),
    val currentSessionId: Long? = null,
    val dailyVerse: DailyVerse? = null,
    val prayerTimes: List<PrayerTime> = emptyList(),
    val prayerPlace: String = "",
    val packages: List<SubscriptionPackage> = emptyList(),
    val subscription: SubscriptionOverview = SubscriptionOverview(),
    val storeProducts: Map<String, StoreProduct> = emptyMap(),
    val passkeyCount: Int = 0,
    val loading: Boolean = false,
    val sending: Boolean = false,
    val recording: Boolean = false,
    val notice: String? = null,
    val error: String? = null,
)

data class PasskeyOptions(
    val transactionId: String,
    val optionsJson: String,
)

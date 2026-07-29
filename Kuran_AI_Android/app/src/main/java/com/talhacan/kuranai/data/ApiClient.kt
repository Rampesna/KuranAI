package com.talhacan.kuranai.data

import com.google.gson.JsonElement
import com.google.gson.JsonNull
import com.google.gson.JsonObject
import com.google.gson.JsonParser
import com.talhacan.kuranai.BuildConfig
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.AuthTokens
import com.talhacan.kuranai.model.ChatMessage
import com.talhacan.kuranai.model.ChatSession
import com.talhacan.kuranai.model.DailyVerse
import com.talhacan.kuranai.model.PasskeyOptions
import com.talhacan.kuranai.model.PrayerTime
import com.talhacan.kuranai.model.SourceType
import com.talhacan.kuranai.model.SubscriptionOverview
import com.talhacan.kuranai.model.SubscriptionPackage
import com.talhacan.kuranai.model.User
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import kotlinx.coroutines.withContext
import okhttp3.Headers
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import okhttp3.logging.HttpLoggingInterceptor
import okio.BufferedSource
import java.io.File
import java.io.IOException
import java.util.UUID
import java.util.concurrent.TimeUnit

class ApiException(
    val statusCode: Int,
    override val message: String,
    val payload: JsonObject? = null,
) : IOException(message)

class ApiClient(private val sessionStore: SessionStore) {
    private val jsonType = "application/json; charset=utf-8".toMediaType()
    private val refreshMutex = Mutex()
    private val client = OkHttpClient.Builder()
        .connectTimeout(20, TimeUnit.SECONDS)
        .readTimeout(120, TimeUnit.SECONDS)
        .writeTimeout(60, TimeUnit.SECONDS)
        .addInterceptor(
            HttpLoggingInterceptor().apply {
                level = if (BuildConfig.DEBUG) {
                    HttpLoggingInterceptor.Level.BASIC
                } else {
                    HttpLoggingInterceptor.Level.NONE
                }
            },
        )
        .build()

    suspend fun login(credential: String, password: String): AuthTokens =
        authRequest(
            "/api/user/auth/login",
            jsonObject("credential" to credential, "password" to password),
        )

    suspend fun register(name: String, username: String, email: String, password: String) {
        requestData(
            "/api/user/auth/register",
            method = "POST",
            body = jsonObject(
                "name" to name,
                "username" to username,
                "email" to email,
                "password" to password,
            ),
            authenticated = false,
        )
    }

    suspend fun loginGoogle(idToken: String): AuthTokens =
        authRequest("/api/user/auth/login/google", jsonObject("idToken" to idToken))

    suspend fun loginApple(identityToken: String): AuthTokens =
        authRequest("/api/user/auth/login/apple", jsonObject("identityToken" to identityToken))

    suspend fun finishMobileAuth(code: String): AuthTokens =
        authRequest("/api/user/auth/login/mobile/finish", jsonObject("code" to code))

    suspend fun logout() {
        runCatching { requestData("/api/user/auth/logout", "POST") }
        sessionStore.clearTokens()
    }

    suspend fun me(): User {
        val data = requestData("/api/user/auth/me")
        return parseUser(data.asJsonObject)
    }

    suspend fun updateUsername(username: String): User {
        val data = requestData(
            "/api/user/profile/username",
            "PATCH",
            jsonObject("username" to username),
        )
        return parseUser(data.asJsonObject)
    }

    suspend fun updatePassword(currentPassword: String, newPassword: String) {
        requestData(
            "/api/user/profile/password",
            "PATCH",
            jsonObject(
                "currentPassword" to currentPassword,
                "newPassword" to newPassword,
            ),
        )
    }

    suspend fun deleteAccount() {
        requestData("/api/user/profile/account", "DELETE")
        sessionStore.clearTokens()
    }

    suspend fun uploadAvatar(file: File): User {
        val body = MultipartBody.Builder()
            .setType(MultipartBody.FORM)
            .addFormDataPart(
                "avatar",
                "avatar.jpg",
                file.asRequestBody("image/jpeg".toMediaType()),
            )
            .build()
        val data = requestData("/api/user/profile/avatar", "POST", requestBody = body)
        return parseUser(data.asJsonObject["user"].asJsonObject)
    }

    suspend fun deleteAvatar(): User? {
        val data = requestData("/api/user/profile/avatar", "DELETE")
        return if (data.isJsonObject) parseUser(data.asJsonObject) else null
    }

    suspend fun beginPasskeyLogin(): PasskeyOptions =
        parsePasskeyOptions(
            requestData(
                "/api/user/auth/passkey/login/options",
                method = "POST",
                authenticated = false,
            ),
        )

    suspend fun finishPasskeyLogin(transactionId: String, credentialJson: String): AuthTokens {
        val data = requestData(
            "/api/user/auth/passkey/login/finish",
            method = "POST",
            requestBody = credentialJson.toRequestBody(jsonType),
            authenticated = false,
            headers = Headers.headersOf("X-Passkey-Transaction", transactionId),
        )
        return parseTokens(data.asJsonObject).also(sessionStore::saveTokens)
    }

    suspend fun beginPasskeyRegistration(): PasskeyOptions =
        parsePasskeyOptions(
            requestData(
                "/api/user/auth/passkey/register/options",
                method = "POST",
            ),
        )

    suspend fun finishPasskeyRegistration(transactionId: String, credentialJson: String) {
        requestData(
            "/api/user/auth/passkey/register/finish",
            method = "POST",
            requestBody = credentialJson.toRequestBody(jsonType),
            headers = Headers.headersOf("X-Passkey-Transaction", transactionId),
        )
    }

    suspend fun passkeyCount(): Int =
        requestData("/api/user/auth/passkey").asJsonObject.int("count")

    suspend fun deletePasskeys() {
        requestData("/api/user/auth/passkey", "DELETE")
    }

    suspend fun sessions(): List<ChatSession> {
        val data = requestData("/api/v1/chat/sessions?pageIndex=0&pageSize=30").asJsonObject
        return data.array("items").map { item ->
            item.asJsonObject.let {
                ChatSession(
                    id = it.long("id"),
                    firstMessage = it.stringOrNull("firstMessage").orEmpty(),
                    updatedAt = it.stringOrNull("updatedAt").orEmpty(),
                )
            }
        }
    }

    suspend fun messages(sessionId: Long): List<ChatMessage> {
        val data = requestData(
            "/api/v1/chat/sessions/$sessionId/messages?pageIndex=0&pageSize=100",
        ).asJsonObject
        return data.array("items").map { item ->
            item.asJsonObject.let {
                ChatMessage(
                    id = it.long("id").toString(),
                    role = it.string("role"),
                    content = it.string("content"),
                )
            }
        }
    }

    suspend fun streamMessage(
        sessionId: Long?,
        source: SourceType,
        content: String,
        onStarted: suspend (Long) -> Unit,
        onChunk: suspend (String) -> Unit,
    ): Long = withContext(Dispatchers.IO) {
        val requestBody = JsonObject().apply {
            if (sessionId != null) addProperty("sessionId", sessionId)
            addProperty("sourceType", source.apiValue)
            addProperty("content", content)
        }
        val callRequest = buildRequest(
            "/api/v1/chat/message",
            "POST",
            requestBody.toString().toRequestBody(jsonType),
            authenticated = true,
            extraHeaders = Headers.headersOf("Accept", "text/event-stream"),
        )
        val response = client.newCall(callRequest).execute()
        if (response.code == 401 && refresh()) {
            response.close()
            return@withContext streamMessage(sessionId, source, content, onStarted, onChunk)
        }
        if (!response.isSuccessful) {
            val raw = response.body?.string().orEmpty()
            response.close()
            throw parseError(response.code, raw)
        }
        val sourceBuffer = response.body?.source() ?: throw ApiException(500, "Empty response stream")
        var resolvedSessionId = sessionId ?: 0L
        readServerEvents(sourceBuffer) { event ->
            val payload = JsonParser.parseString(event).asJsonObject
            when (payload.string("type")) {
                "started" -> {
                    resolvedSessionId = payload.long("sessionId")
                    onStarted(resolvedSessionId)
                }
                "chunk" -> onChunk(payload.stringOrNull("content").orEmpty())
                "done" -> {
                    resolvedSessionId = payload.long("sessionId")
                    onStarted(resolvedSessionId)
                }
                "error" -> throw ApiException(500, payload.stringOrNull("message") ?: "Stream error")
            }
        }
        response.close()
        resolvedSessionId
    }

    suspend fun transcribe(audio: File): String {
        val body = MultipartBody.Builder()
            .setType(MultipartBody.FORM)
            .addFormDataPart(
                "audio",
                "question.m4a",
                audio.asRequestBody("audio/mp4".toMediaType()),
            )
            .build()
        val data = requestData("/api/v1/chat/transcribe", "POST", requestBody = body)
        return data.asJsonObject.string("text")
    }

    suspend fun tts(text: String, voice: String): ByteArray = withContext(Dispatchers.IO) {
        val request = buildRequest(
            "/api/v1/chat/tts",
            "POST",
            jsonObject("text" to text, "voice" to voice).toRequestBody(jsonType),
            authenticated = true,
            extraHeaders = Headers.headersOf("Accept", "audio/mpeg"),
        )
        val response = client.newCall(request).execute()
        if (!response.isSuccessful) {
            val raw = response.body?.string().orEmpty()
            response.close()
            throw parseError(response.code, raw)
        }
        response.body?.bytes() ?: throw ApiException(500, "Empty audio response")
    }

    suspend fun verse(language: AppLanguage): DailyVerse? {
        val data = runCatching {
            requestData(
                "/api/verse-of-day?lang=${language.code}",
                authenticated = false,
            ).asJsonObject
        }.getOrNull() ?: return null
        val translations = data.array("translations")
        val translation = translations.firstOrNull()?.asJsonObject
        val arabicText = data.stringOrNull("arabicText").orEmpty()
        if (arabicText.isEmpty()) return null
        return DailyVerse(
            arabicText = arabicText,
            translation = translation?.stringOrNull("translationText")
                ?: if (language == AppLanguage.Arabic) arabicText else "",
            reference = buildString {
                append(translation?.stringOrNull("sourceName") ?: "")
                if (isNotEmpty()) append(" ")
                append(data.int("ayahNumber"))
            },
        )
    }

    suspend fun prayerTimes(
        latitude: Double?,
        longitude: Double?,
        language: AppLanguage,
    ): Pair<String, List<PrayerTime>> {
        val query = if (latitude != null && longitude != null) {
            "?lat=$latitude&lng=$longitude"
        } else {
            ""
        }
        val data = requestData("/api/v1/prayer-times$query", authenticated = false).asJsonObject
        val place = data.objectOrNull("place")
        val times = data.objectOrNull("times") ?: JsonObject()
        val names = when (language) {
            AppLanguage.Turkish -> listOf("İmsak", "Güneş", "Öğle", "İkindi", "Akşam", "Yatsı")
            AppLanguage.English -> listOf("Fajr", "Sunrise", "Dhuhr", "Asr", "Maghrib", "Isha")
            AppLanguage.Arabic -> listOf("الفجر", "الشروق", "الظهر", "العصر", "المغرب", "العشاء")
        }
        val keys = listOf("imsak", "gunes", "ogle", "ikindi", "aksam", "yatsi")
        return (place?.stringOrNull("name").orEmpty()) to keys.mapIndexed { index, key ->
            PrayerTime(names[index], times.stringOrNull(key) ?: "--:--")
        }
    }

    suspend fun packages(): List<SubscriptionPackage> {
        val data = requestData("/api/v1/subscription/packages").asJsonObject
        return data.array("packages").map { parsePackage(it.asJsonObject) }
    }

    suspend fun subscriptionOverview(): SubscriptionOverview {
        val data = requestData("/api/v1/subscription/me").asJsonObject
        val subscription = data.objectOrNull("subscription")
        val usage = data.objectOrNull("usage")
        val currentPackage = subscription?.objectOrNull("package")?.let(::parsePackage)
        return SubscriptionOverview(
            hasPackage = subscription?.boolean("hasPackage") ?: false,
            currentPackage = currentPackage,
            expiresAt = subscription?.stringOrNull("expiresAt"),
            messageUsed = usage?.int("messageUsed") ?: 0,
            messageLimit = usage?.int("messageLimit") ?: 3,
            ttsUsed = usage?.int("ttsUsed") ?: 0,
            ttsLimit = usage?.int("ttsLimit") ?: 0,
        )
    }

    suspend fun syncGooglePurchase(
        productId: String,
        purchaseToken: String,
    ) {
        requestData(
            "/api/v1/subscription/google/sync",
            "POST",
            jsonObject(
                "packageName" to "com.talhacan.kuranai",
                "productId" to productId,
                "purchaseToken" to purchaseToken,
            ),
        )
    }

    private suspend fun authRequest(path: String, body: String): AuthTokens {
        val data = requestData(
            path,
            method = "POST",
            requestBody = body.toRequestBody(jsonType),
            authenticated = false,
        )
        return parseTokens(data.asJsonObject).also(sessionStore::saveTokens)
    }

    private fun parseTokens(data: JsonObject): AuthTokens = AuthTokens(
        accessToken = data.string("accessToken"),
        refreshToken = data.string("refreshToken"),
    )

    private suspend fun requestData(
        path: String,
        method: String = "GET",
        body: String? = null,
        requestBody: RequestBody? = body?.toRequestBody(jsonType),
        authenticated: Boolean = true,
        headers: Headers = Headers.Builder().build(),
        retryAfterRefresh: Boolean = true,
    ): JsonElement = withContext(Dispatchers.IO) {
        val response = client.newCall(
            buildRequest(path, method, requestBody, authenticated, headers),
        ).execute()
        val raw = response.body?.string().orEmpty()
        val statusCode = response.code
        response.close()

        if (statusCode == 401 && authenticated && retryAfterRefresh && refresh()) {
            return@withContext requestData(
                path,
                method,
                body,
                requestBody,
                authenticated,
                headers,
                retryAfterRefresh = false,
            )
        }
        if (statusCode !in 200..299) throw parseError(statusCode, raw)
        if (raw.isBlank()) return@withContext JsonObject()

        val root = JsonParser.parseString(raw).asJsonObject
        if (!root.boolean("isSuccess")) {
            throw ApiException(statusCode, root.stringOrNull("message") ?: "Request failed", root)
        }
        root["data"] ?: JsonObject()
    }

    private suspend fun refresh(): Boolean = refreshMutex.withLock {
        val current = sessionStore.tokens() ?: return@withLock false
        return@withLock runCatching {
            val data = requestData(
                "/api/user/auth/refresh",
                "POST",
                jsonObject("refreshToken" to current.refreshToken),
                authenticated = false,
            )
            sessionStore.saveTokens(parseTokens(data.asJsonObject))
            true
        }.getOrElse {
            sessionStore.clearTokens()
            false
        }
    }

    private fun buildRequest(
        path: String,
        method: String,
        body: RequestBody?,
        authenticated: Boolean,
        extraHeaders: Headers,
    ): Request {
        val builder = Request.Builder()
            .url(BuildConfig.API_BASE_URL.trimEnd('/') + path)
            .header("Accept", "application/json")
            .header("X-Device-Type", "android")
            .headers(mergeHeaders(extraHeaders))
        if (authenticated) {
            sessionStore.tokens()?.accessToken?.let {
                builder.header("Authorization", "Bearer $it")
            }
        }
        builder.method(method, if (method == "GET" || method == "DELETE") null else body ?: EMPTY_BODY)
        return builder.build()
    }

    private fun mergeHeaders(headers: Headers): Headers {
        val builder = Headers.Builder()
            .add("Accept", "application/json")
            .add("X-Device-Type", "android")
        for (index in 0 until headers.size) {
            builder.set(headers.name(index), headers.value(index))
        }
        return builder.build()
    }

    private fun parsePasskeyOptions(data: JsonElement): PasskeyOptions {
        val value = data.asJsonObject
        return PasskeyOptions(
            transactionId = value.string("transactionId"),
            optionsJson = value["options"].toString(),
        )
    }

    private fun parseUser(data: JsonObject): User = User(
        id = data.long("id"),
        email = data.stringOrNull("email").orEmpty(),
        name = data.stringOrNull("name").orEmpty(),
        username = data.stringOrNull("username").orEmpty(),
        avatarUrl = data.stringOrNull("avatarUrl")?.let {
            if (it.startsWith("http")) it else BuildConfig.API_BASE_URL.trimEnd('/') + it
        },
        hasPassword = data.boolean("hasPassword"),
    )

    private fun parsePackage(data: JsonObject): SubscriptionPackage = SubscriptionPackage(
        id = data.long("id"),
        name = data.stringOrNull("name").orEmpty(),
        nameEn = data.stringOrNull("nameEn").orEmpty(),
        slug = data.stringOrNull("slug").orEmpty(),
        dailyQuestionLimit = data.intOrNull("dailyQuestionLimit"),
        dailyTtsLimit = data.intOrNull("dailyTTSLimit"),
        monthlyPriceTl = data.doubleOrNull("monthlyPriceTL") ?: 0.0,
        yearlyPriceTl = data.doubleOrNull("yearlyPriceTL") ?: 0.0,
    )

    private suspend fun readServerEvents(source: BufferedSource, onEvent: suspend (String) -> Unit) {
        var dataLine: String? = null
        while (!source.exhausted()) {
            val line = source.readUtf8Line() ?: break
            if (line.startsWith("data:")) {
                dataLine = line.removePrefix("data:").trim()
            }
            if (line.isBlank() && !dataLine.isNullOrBlank()) {
                onEvent(dataLine!!)
                dataLine = null
            }
        }
    }

    private fun parseError(statusCode: Int, raw: String): ApiException {
        val root = runCatching { JsonParser.parseString(raw).asJsonObject }.getOrNull()
        return ApiException(
            statusCode,
            root?.stringOrNull("message") ?: "HTTP $statusCode",
            root,
        )
    }

    companion object {
        private val EMPTY_BODY = ByteArray(0).toRequestBody(null)

        private fun jsonObject(vararg values: Pair<String, Any?>): String =
            JsonObject().apply {
                values.forEach { (key, value) ->
                    when (value) {
                        null -> add(key, JsonNull.INSTANCE)
                        is String -> addProperty(key, value)
                        is Number -> addProperty(key, value)
                        is Boolean -> addProperty(key, value)
                        else -> addProperty(key, value.toString())
                    }
                }
            }.toString()
    }
}

private fun JsonObject.string(name: String): String =
    stringOrNull(name) ?: throw ApiException(500, "Missing $name")

private fun JsonObject.stringOrNull(name: String): String? =
    get(name)?.takeUnless { it.isJsonNull }?.asString

private fun JsonObject.long(name: String): Long =
    get(name)?.takeUnless { it.isJsonNull }?.asLong ?: 0L

private fun JsonObject.int(name: String): Int = intOrNull(name) ?: 0

private fun JsonObject.intOrNull(name: String): Int? =
    get(name)?.takeUnless { it.isJsonNull }?.asInt

private fun JsonObject.doubleOrNull(name: String): Double? =
    get(name)?.takeUnless { it.isJsonNull }?.asDouble

private fun JsonObject.boolean(name: String): Boolean =
    get(name)?.takeUnless { it.isJsonNull }?.asBoolean ?: false

private fun JsonObject.array(name: String): List<JsonElement> =
    getAsJsonArray(name)?.toList().orEmpty()

private fun JsonObject.objectOrNull(name: String): JsonObject? =
    get(name)?.takeIf { it.isJsonObject }?.asJsonObject

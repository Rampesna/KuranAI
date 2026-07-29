package com.talhacan.kuranai.auth

import android.app.Activity
import androidx.credentials.ClearCredentialStateRequest
import androidx.credentials.CreatePublicKeyCredentialRequest
import androidx.credentials.CreatePublicKeyCredentialResponse
import androidx.credentials.CredentialManager
import androidx.credentials.CustomCredential
import androidx.credentials.GetCredentialRequest
import androidx.credentials.GetPublicKeyCredentialOption
import androidx.credentials.PublicKeyCredential
import androidx.credentials.exceptions.NoCredentialException
import com.google.android.libraries.identity.googleid.GetGoogleIdOption
import com.google.android.libraries.identity.googleid.GoogleIdTokenCredential
import com.talhacan.kuranai.BuildConfig
import com.talhacan.kuranai.data.ApiClient
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.AuthTokens
import com.talhacan.kuranai.ui.t

class IdentityManager(private val apiClient: ApiClient) {
    suspend fun signInWithGoogle(activity: Activity, language: AppLanguage): AuthTokens {
        val googleOption = GetGoogleIdOption.Builder()
            .setServerClientId(BuildConfig.GOOGLE_WEB_CLIENT_ID)
            .setFilterByAuthorizedAccounts(false)
            .setAutoSelectEnabled(false)
            .build()
        val response = try {
            CredentialManager.create(activity).getCredential(
                activity,
                GetCredentialRequest.Builder()
                    .addCredentialOption(googleOption)
                    .build(),
            )
        } catch (error: NoCredentialException) {
            throw IllegalStateException(
                language.t(
                    "Bu cihazda kullanılabilir bir Google hesabı bulunamadı.",
                    "No Google account is available on this device.",
                    "لا يتوفر حساب Google على هذا الجهاز.",
                ),
                error,
            )
        }
        val credential = response.credential
        require(
            credential is CustomCredential &&
                credential.type == GoogleIdTokenCredential.TYPE_GOOGLE_ID_TOKEN_CREDENTIAL,
        ) { "Google credential was not returned" }
        val idToken = GoogleIdTokenCredential.createFrom(credential.data).idToken
        return apiClient.loginGoogle(idToken)
    }

    suspend fun signInWithPasskey(activity: Activity, language: AppLanguage): AuthTokens {
        val challenge = apiClient.beginPasskeyLogin()
        val response = try {
            CredentialManager.create(activity).getCredential(
                activity,
                GetCredentialRequest.Builder()
                    .addCredentialOption(GetPublicKeyCredentialOption(challenge.optionsJson))
                    .build(),
            )
        } catch (error: NoCredentialException) {
            throw IllegalStateException(
                language.t(
                    "Bu cihazda Kuran AI için kayıtlı bir geçiş anahtarı bulunamadı.",
                    "No passkey saved for Kuran AI was found on this device.",
                    "لم يتم العثور على مفتاح مرور محفوظ لتطبيق Kuran AI على هذا الجهاز.",
                ),
                error,
            )
        }
        val credential = response.credential as? PublicKeyCredential
            ?: error(
                language.t(
                    "Geçiş anahtarı bilgisi alınamadı.",
                    "The passkey credential was not returned.",
                    "تعذّر الحصول على بيانات مفتاح المرور.",
                ),
            )
        return apiClient.finishPasskeyLogin(
            challenge.transactionId,
            credential.authenticationResponseJson,
        )
    }

    suspend fun createPasskey(activity: Activity, language: AppLanguage) {
        val challenge = apiClient.beginPasskeyRegistration()
        val response = CredentialManager.create(activity).createCredential(
            activity,
            CreatePublicKeyCredentialRequest(challenge.optionsJson),
        )
        val credential = response as? CreatePublicKeyCredentialResponse
            ?: error(
                language.t(
                    "Geçiş anahtarı oluşturulamadı.",
                    "The passkey could not be created.",
                    "تعذّر إنشاء مفتاح المرور.",
                ),
            )
        apiClient.finishPasskeyRegistration(
            challenge.transactionId,
            credential.registrationResponseJson,
        )
    }

    suspend fun clear(activity: Activity) {
        runCatching {
            CredentialManager.create(activity).clearCredentialState(
                ClearCredentialStateRequest(),
            )
        }
    }
}

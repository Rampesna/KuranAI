package com.talhacan.kuranai

import android.app.Application
import com.talhacan.kuranai.auth.IdentityManager
import com.talhacan.kuranai.billing.BillingStore
import com.talhacan.kuranai.data.ApiClient
import com.talhacan.kuranai.data.SessionStore
import kotlinx.coroutines.flow.MutableSharedFlow

class KuranAIApplication : Application() {
    lateinit var sessionStore: SessionStore
        private set
    lateinit var apiClient: ApiClient
        private set
    lateinit var identityManager: IdentityManager
        private set
    lateinit var billingStore: BillingStore
        private set

    val entitlementEvents = MutableSharedFlow<Unit>(extraBufferCapacity = 1)

    override fun onCreate() {
        super.onCreate()
        sessionStore = SessionStore(this)
        apiClient = ApiClient(sessionStore)
        identityManager = IdentityManager(apiClient)
        billingStore = BillingStore(this, apiClient) {
            entitlementEvents.emit(Unit)
        }
    }
}

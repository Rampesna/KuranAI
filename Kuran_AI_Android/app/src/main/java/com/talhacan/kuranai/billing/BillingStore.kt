package com.talhacan.kuranai.billing

import android.app.Activity
import android.app.Application
import android.util.Log
import com.android.billingclient.api.AcknowledgePurchaseParams
import com.android.billingclient.api.BillingClient
import com.android.billingclient.api.BillingClientStateListener
import com.android.billingclient.api.BillingFlowParams
import com.android.billingclient.api.BillingResult
import com.android.billingclient.api.PendingPurchasesParams
import com.android.billingclient.api.ProductDetails
import com.android.billingclient.api.Purchase
import com.android.billingclient.api.PurchasesUpdatedListener
import com.android.billingclient.api.QueryProductDetailsParams
import com.android.billingclient.api.QueryProductDetailsResult
import com.android.billingclient.api.QueryPurchasesParams
import com.android.billingclient.api.acknowledgePurchase
import com.android.billingclient.api.queryPurchasesAsync
import com.talhacan.kuranai.data.ApiClient
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.StoreProduct
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.suspendCancellableCoroutine
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import kotlin.coroutines.resume

class BillingStore(
    application: Application,
    private val apiClient: ApiClient,
    private val onEntitlementChanged: suspend () -> Unit,
) : PurchasesUpdatedListener {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.Main)
    private val syncMutex = Mutex()
    private var currentUserId: Long? = null
    private var language = AppLanguage.Turkish
    private val _products = MutableStateFlow<Map<String, StoreProduct>>(emptyMap())
    val products: StateFlow<Map<String, StoreProduct>> = _products.asStateFlow()
    private val _message = MutableStateFlow<String?>(null)
    val message: StateFlow<String?> = _message.asStateFlow()
    private var detailsById: Map<String, ProductDetails> = emptyMap()

    private val billingClient = BillingClient.newBuilder(application)
        .setListener(this)
        .enableAutoServiceReconnection()
        .enablePendingPurchases(
            PendingPurchasesParams.newBuilder().enableOneTimeProducts().build(),
        )
        .build()

    init {
        connect()
    }

    fun attachUser(userId: Long?) {
        currentUserId = userId
        if (userId != null) {
            scope.launch {
                ensureConnected()
                queryOwned().forEach { purchase ->
                    syncAndAcknowledge(purchase)
                }
            }
        }
    }

    fun setLanguage(language: AppLanguage) {
        this.language = language
    }

    fun clearMessage() {
        _message.value = null
    }

    fun purchase(activity: Activity, productId: String) {
        scope.launch {
            val userId = currentUserId
            if (userId == null) {
                _message.value = localized(
                    "Satın almak için giriş yapın.",
                    "Sign in to make a purchase.",
                    "سجّل الدخول لإتمام عملية الشراء.",
                )
                return@launch
            }
            ensureConnected()
            if (detailsById[productId] == null) refreshProducts()
            val details = detailsById[productId]
            val offerToken = details?.subscriptionOfferDetails?.firstOrNull()?.offerToken
            if (details == null || offerToken.isNullOrBlank()) {
                _message.value = localized(
                    "Bu abonelik Google Play'de henüz kullanılamıyor.",
                    "This subscription is not available on Google Play yet.",
                    "هذا الاشتراك غير متاح على Google Play بعد.",
                )
                return@launch
            }
            val params = BillingFlowParams.newBuilder()
                .setProductDetailsParamsList(
                    listOf(
                        BillingFlowParams.ProductDetailsParams.newBuilder()
                            .setProductDetails(details)
                            .setOfferToken(offerToken)
                            .build(),
                    ),
                )
                .setObfuscatedAccountId(userId.toString())
                .build()
            val result = billingClient.launchBillingFlow(activity, params)
            if (result.responseCode != BillingClient.BillingResponseCode.OK) {
                _message.value = billingMessage(result)
            }
        }
    }

    fun restore() {
        scope.launch {
            val purchases = queryOwned()
            if (purchases.isEmpty()) {
                _message.value = localized(
                    "Bu Google hesabında geri yüklenecek abonelik bulunamadı.",
                    "No subscription was found to restore for this Google account.",
                    "لم يتم العثور على اشتراك لاستعادته في حساب Google هذا.",
                )
                return@launch
            }
            purchases.forEach { purchase ->
                syncAndAcknowledge(purchase)
            }
            _message.value = localized(
                "Satın almalar doğrulandı.",
                "Purchases verified.",
                "تم التحقق من المشتريات.",
            )
        }
    }

    override fun onPurchasesUpdated(result: BillingResult, purchases: MutableList<Purchase>?) {
        scope.launch {
            when (result.responseCode) {
                BillingClient.BillingResponseCode.OK -> {
                    purchases.orEmpty().forEach { purchase ->
                        syncAndAcknowledge(purchase)
                    }
                    _message.value = localized(
                        "Aboneliğiniz etkinleştirildi.",
                        "Your subscription is active.",
                        "تم تفعيل اشتراكك.",
                    )
                }
                BillingClient.BillingResponseCode.ITEM_ALREADY_OWNED -> restore()
                BillingClient.BillingResponseCode.USER_CANCELED -> Unit
                else -> _message.value = billingMessage(result)
            }
        }
    }

    private fun connect() {
        if (billingClient.isReady) {
            scope.launch { refreshProducts() }
            return
        }
        billingClient.startConnection(object : BillingClientStateListener {
            override fun onBillingSetupFinished(result: BillingResult) {
                if (result.responseCode == BillingClient.BillingResponseCode.OK) {
                    scope.launch {
                        refreshProducts()
                        queryOwned().forEach { purchase ->
                            syncAndAcknowledge(purchase)
                        }
                    }
                } else {
                    Log.w(TAG, "Billing setup: ${result.debugMessage}")
                }
            }

            override fun onBillingServiceDisconnected() = Unit
        })
    }

    private suspend fun ensureConnected() {
        if (billingClient.isReady) return
        connect()
        repeat(30) {
            if (billingClient.isReady) return
            delay(100)
        }
    }

    private suspend fun refreshProducts() {
        if (!billingClient.isReady) return
        val params = QueryProductDetailsParams.newBuilder()
            .setProductList(
                PRODUCT_IDS.map { id ->
                    QueryProductDetailsParams.Product.newBuilder()
                        .setProductId(id)
                        .setProductType(BillingClient.ProductType.SUBS)
                        .build()
                },
            )
            .build()
        val (result, details) = queryProducts(params)
        if (result.responseCode != BillingClient.BillingResponseCode.OK) {
            Log.w(TAG, "Product query: ${result.debugMessage}")
            return
        }
        detailsById = details.productDetailsList.associateBy { it.productId }
        _products.value = detailsById.mapValues { (id, product) ->
            StoreProduct(
                id,
                product.subscriptionOfferDetails
                    ?.firstOrNull()
                    ?.pricingPhases
                    ?.pricingPhaseList
                    ?.firstOrNull()
                    ?.formattedPrice
                    .orEmpty(),
            )
        }
    }

    private suspend fun queryProducts(
        params: QueryProductDetailsParams,
    ): Pair<BillingResult, QueryProductDetailsResult> = suspendCancellableCoroutine { continuation ->
        billingClient.queryProductDetailsAsync(params) { result, details ->
            if (continuation.isActive) continuation.resume(result to details)
        }
    }

    private suspend fun queryOwned(): List<Purchase> {
        if (!billingClient.isReady) return emptyList()
        val result = billingClient.queryPurchasesAsync(
            QueryPurchasesParams.newBuilder()
                .setProductType(BillingClient.ProductType.SUBS)
                .build(),
        )
        return if (result.billingResult.responseCode == BillingClient.BillingResponseCode.OK) {
            result.purchasesList.filter { it.purchaseState == Purchase.PurchaseState.PURCHASED }
        } else {
            emptyList()
        }
    }

    private suspend fun syncAndAcknowledge(purchase: Purchase) = syncMutex.withLock {
        if (currentUserId == null || purchase.purchaseState != Purchase.PurchaseState.PURCHASED) {
            return@withLock
        }
        val productId = purchase.products.firstOrNull().orEmpty()
        if (productId.isBlank() || purchase.purchaseToken.isBlank()) return@withLock
        try {
            apiClient.syncGooglePurchase(productId, purchase.purchaseToken)
            if (!purchase.isAcknowledged) {
                val result = billingClient.acknowledgePurchase(
                    AcknowledgePurchaseParams.newBuilder()
                        .setPurchaseToken(purchase.purchaseToken)
                        .build(),
                )
                if (result.responseCode != BillingClient.BillingResponseCode.OK) {
                    Log.w(TAG, "Acknowledge: ${result.debugMessage}")
                }
            }
            onEntitlementChanged()
        } catch (error: Exception) {
            _message.value = error.localizedMessage ?: localized(
                "Satın alma doğrulanamadı.",
                "The purchase could not be verified.",
                "تعذّر التحقق من عملية الشراء.",
            )
        }
    }

    private fun billingMessage(result: BillingResult): String =
        "Google Play: ${result.debugMessage.ifBlank {
            localized(
                "kod ${result.responseCode}",
                "code ${result.responseCode}",
                "الرمز ${result.responseCode}",
            )
        }}"

    private fun localized(tr: String, en: String, ar: String): String = when (language) {
        AppLanguage.Turkish -> tr
        AppLanguage.English -> en
        AppLanguage.Arabic -> ar
    }

    companion object {
        private const val TAG = "KuranAIBilling"
        val PRODUCT_IDS = listOf(
            "kuranai.bronze.monthly",
            "kuranai.bronze.yearly",
            "kuranai.silver.monthly",
            "kuranai.silver.yearly",
            "kuranai.gold.monthly",
            "kuranai.gold.yearly",
            "kuranai.diamond.monthly",
            "kuranai.diamond.yearly",
        )
    }
}

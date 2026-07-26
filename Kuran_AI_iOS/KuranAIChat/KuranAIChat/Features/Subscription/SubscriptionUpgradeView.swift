import StoreKit
import SwiftUI

private let collapsedCardHeight: CGFloat = 64
private let cardCornerRadius: CGFloat = 20

struct SubscriptionUpgradeView: View {
    private static let pendingReceiptSyncKey = "kuranai.pendingReceiptSync"

    private static let fallbackPackages: [SubscriptionPackageDTO] = [
        .init(id: 1, name: "Bronze", nameEn: "Bronze", slug: "bronze", dailyQuestionLimit: 100, dailyTTSLimit: 5, monthlyPriceTL: 29.99, yearlyPriceTL: 299.99),
        .init(id: 2, name: "Gümüş", nameEn: "Silver", slug: "silver", dailyQuestionLimit: 250, dailyTTSLimit: 50, monthlyPriceTL: 49.99, yearlyPriceTL: 499.99),
        .init(id: 3, name: "Altın", nameEn: "Gold", slug: "gold", dailyQuestionLimit: nil, dailyTTSLimit: 100, monthlyPriceTL: 79.99, yearlyPriceTL: 799.99),
        .init(id: 4, name: "Elmas", nameEn: "Diamond", slug: "diamond", dailyQuestionLimit: nil, dailyTTSLimit: nil, monthlyPriceTL: 149.99, yearlyPriceTL: 1499.99)
    ]

    private static let storeProductIDs = fallbackPackages.flatMap { package in
        ["kuranai.\(package.slug).monthly", "kuranai.\(package.slug).yearly"]
    }

    let language: Language
    let limitInfo: APIError.LimitExceededInfo?
    let onDismiss: () -> Void
    let onPurchaseSuccess: () -> Void

    @State private var packages: [SubscriptionPackageDTO] = []
    @State private var storeProducts: [String: Product] = [:]
    @State private var isLoading = true
    @State private var errorMessage: String?
    @State private var purchasingProductId: String?
    @State private var selectedIndex: Int = 0
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationStack {
            ZStack {
                Theme.background.ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 0) {
                        headerSection

                        if isLoading {
                            loadingSection
                        } else {
                            stackedCardsSection
                        }

                        if let err = errorMessage {
                            errorBanner(err)
                        }

                        restoreButton

                        subscriptionLegalLinks
                    }
                }
            }
            .navigationBarHidden(true)
            .overlay(alignment: .topLeading) {
                Button {
                    onDismiss()
                    dismiss()
                } label: {
                    Image(systemName: "chevron.left")
                        .font(.system(size: 18, weight: .semibold))
                        .foregroundColor(Theme.textPrimary)
                        .frame(width: 44, height: 44)
                }
                .padding(.top, 8)
            }
        }
        .task {
            await syncPendingPurchaseIfNeeded()
            await loadPackages()
        }
        .onChange(of: packages.count) { _, count in
            if count > 0 { selectedIndex = count - 1 }
        }
    }

    private var headerSection: some View {
        VStack(spacing: 8) {
            if limitInfo != nil {
                limitBanner
            }
                Text(language.localized("Sınırsız Erişim", "Get Unlimited Access"))
                .font(.system(size: 28, weight: .bold))
                .foregroundColor(Theme.textPrimary)
                .padding(.top, limitInfo != nil ? 8 : 24)
            Text(language.localized("Paket seçin ve özelliklerin kilidini açın", "Choose a plan and unlock all features"))
                .font(.system(size: 15))
                .foregroundColor(Theme.textMuted)
        }
        .frame(maxWidth: .infinity)
        .padding(.horizontal, 24)
        .padding(.bottom, 20)
    }

    private var limitBanner: some View {
        HStack(spacing: 12) {
            Image(systemName: "exclamationmark.triangle.fill")
                .font(.system(size: 20))
                .foregroundColor(Theme.neonPrimary)
            VStack(alignment: .leading, spacing: 2) {
                if let info = limitInfo {
                    Text(info.message)
                        .font(.subheadline.weight(.semibold))
                        .foregroundColor(Theme.textPrimary)
                }
                Text(language.localized("Daha fazla kullanım için bir paket seçin.", "Choose a package for more usage."))
                    .font(.caption)
                    .foregroundColor(Theme.textMuted)
            }
            Spacer(minLength: 0)
        }
        .padding(16)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(Theme.glassColor)
                .overlay(RoundedRectangle(cornerRadius: 14).stroke(Theme.neonPrimary.opacity(0.3), lineWidth: 1))
        )
    }

    private var loadingSection: some View {
        VStack(spacing: 16) {
            ProgressView().tint(Theme.textPrimary).scaleEffect(1.2)
            Text(language.localized("Paketler yükleniyor...", "Loading packages..."))
                .font(.subheadline)
                .foregroundColor(Theme.textMuted)
        }
        .frame(maxWidth: .infinity)
        .padding(48)
    }

    private var stackedCardsSection: some View {
        VStack(spacing: 12) {
            ForEach(Array(packages.enumerated()), id: \.element.id) { index, pkg in
                StackableSubscriptionCard(
                    package: pkg,
                    monthlyProduct: storeProducts[productID(for: pkg, period: "monthly")],
                    yearlyProduct: storeProducts[productID(for: pkg, period: "yearly")],
                    language: language,
                    isExpanded: index == selectedIndex,
                    isLast: index == packages.count - 1,
                    isPurchasing: purchasingProductId != nil
                ) { product in
                    Task { await purchase(product: product) }
                } onTap: {
                    withAnimation(.spring(response: 0.35, dampingFraction: 0.85)) {
                        selectedIndex = index
                    }
                }
            }
        }
        .padding(.horizontal, 20)
        .padding(.bottom, 24)
    }

    private func errorBanner(_ message: String) -> some View {
        VStack(spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: "exclamationmark.circle.fill")
                    .foregroundColor(.red.opacity(0.9))
                Text(message)
                    .font(.caption)
                    .foregroundColor(.red.opacity(0.9))
                Spacer(minLength: 0)
            }

            Button {
                Task { await loadPackages() }
            } label: {
                Label(
                    language.localized("Tekrar Dene", "Try Again"),
                    systemImage: "arrow.clockwise"
                )
                .font(.caption.weight(.semibold))
                .foregroundColor(Theme.textPrimary)
                .frame(maxWidth: .infinity)
                .padding(.vertical, 10)
                .background(Theme.glassColor)
                .clipShape(RoundedRectangle(cornerRadius: 10))
            }
            .buttonStyle(.plain)
        }
        .padding(16)
        .background(RoundedRectangle(cornerRadius: 10).fill(Color.red.opacity(0.15)))
        .padding(.horizontal, 20)
        .padding(.bottom, 16)
    }

    private var restoreButton: some View {
        Button {
            Task { await restorePurchases() }
        } label: {
            HStack(spacing: 8) {
                Image(systemName: "arrow.clockwise")
                    .font(.system(size: 13, weight: .medium))
                Text(language.localized("Satın Almaları Geri Yükle", "Restore Purchases"))
                    .font(.subheadline.weight(.medium))
            }
            .foregroundColor(Theme.neonPrimary.opacity(0.9))
        }
        .padding(.bottom, 16)
    }

    /// Apple Guideline 3.1.2: Subscription ekranında Privacy Policy ve Terms of Use linkleri zorunlu.
    private var subscriptionLegalLinks: some View {
        VStack(spacing: 8) {
            Text(language.localized(
                "Abonelikler otomatik olarak yenilenir. İstediğiniz zaman Ayarlar > Apple Kimliği > Abonelikler üzerinden iptal edebilirsiniz.",
                "Subscriptions auto-renew. You can cancel anytime in Settings > Apple ID > Subscriptions."
            ))
            .font(.system(size: 11))
            .foregroundColor(Theme.textMuted)
            .multilineTextAlignment(.center)
            .padding(.horizontal, 24)

            HStack(spacing: 16) {
                Link(language.localized("Gizlilik Politikası", "Privacy Policy"),
                     destination: URL(string: "https://kuranai.talhacan.com/privacy-policy")!)
                Text("·").foregroundColor(Theme.textMuted.opacity(0.5))
                Link(language.localized("Kullanım Koşulları", "Terms of Use"),
                     destination: URL(string: "https://kuranai.talhacan.com/terms")!)
            }
            .font(.system(size: 12, weight: .medium))
            .foregroundColor(Theme.neonPrimary.opacity(0.8))
        }
        .padding(.bottom, 40)
    }

    @MainActor
    private func loadPackages() async {
        isLoading = true
        errorMessage = nil
        defer { isLoading = false }

        // Paket API'si kısa süreli erişilemez olsa bile StoreKit ürünlerini yüklemeyi
        // denemeye devam et. App Store ürünleri hiçbir zaman backend yanıtına bağlı olmamalı.
        let remotePackages = try? await SubscriptionService.shared.fetchPackages()
        packages = normalizedPackages(remotePackages)

        do {
            let fetchedProducts = try await fetchStoreProductsWithRetry()
            storeProducts = Dictionary(uniqueKeysWithValues: fetchedProducts.map { ($0.id, $0) })

            if fetchedProducts.count != Self.storeProductIDs.count {
                errorMessage = language.localized(
                    "Abonelik seçenekleri App Store'dan tam olarak yüklenemedi. Lütfen bağlantınızı kontrol edip tekrar deneyin.",
                    "Subscription options could not be fully loaded from the App Store. Check your connection and try again."
                )
            }
        } catch {
            storeProducts = [:]
            errorMessage = language.localized(
                "Abonelik seçenekleri App Store'dan yüklenemedi. Lütfen tekrar deneyin.",
                "Subscription options could not be loaded from the App Store. Please try again."
            )
        }
    }

    private func normalizedPackages(_ remotePackages: [SubscriptionPackageDTO]?) -> [SubscriptionPackageDTO] {
        guard let remotePackages, !remotePackages.isEmpty else {
            return Self.fallbackPackages
        }

        let bySlug = Dictionary(uniqueKeysWithValues: remotePackages.map { ($0.slug.lowercased(), $0) })
        return Self.fallbackPackages.map { bySlug[$0.slug] ?? $0 }
    }

    private func fetchStoreProductsWithRetry() async throws -> [Product] {
        var lastError: Error?

        for attempt in 0..<3 {
            do {
                let products = try await Product.products(for: Self.storeProductIDs)
                if !products.isEmpty || attempt == 2 {
                    return products
                }
            } catch {
                lastError = error
                if attempt == 2 { throw error }
            }

            try await Task.sleep(for: .milliseconds(600 * (attempt + 1)))
        }

        if let lastError { throw lastError }
        return []
    }

    private func productID(for package: SubscriptionPackageDTO, period: String) -> String {
        "kuranai.\(package.slug.lowercased()).\(period)"
    }

    private func purchase(product: Product) async {
        purchasingProductId = product.id
        errorMessage = nil
        defer { purchasingProductId = nil }

        do {
            let result = try await product.purchase()
            switch result {
            case .success(let verification):
                let transaction: StoreKit.Transaction = try checkVerified(verification)

                // StoreKit 2 may update the app receipt shortly after a successful
                // transaction. Finish the verified transaction first, then retry
                // backend synchronization instead of reporting a completed purchase
                // as failed because the receipt was not available immediately.
                UserDefaults.standard.set(true, forKey: Self.pendingReceiptSyncKey)
                await transaction.finish()

                do {
                    try await syncReceiptToBackendWithRetry()
                } catch {
                    // The App Store purchase is complete. Keep the pending marker so
                    // the next subscription-screen visit can synchronize it again.
                }

                onPurchaseSuccess()
                dismiss()
            case .userCancelled:
                break
            case .pending:
                errorMessage = language.localized("Ödeme bekleniyor.", "Payment pending.")
            @unknown default:
                errorMessage = language.localized("Bilinmeyen hata.", "Unknown error.")
            }
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func checkVerified<T>(_ result: VerificationResult<T>) throws -> T {
        switch result {
        case .unverified:
            throw StoreError.failedVerification
        case .verified(let safe):
            return safe
        }
    }

    private func getAppReceiptBase64() -> String? {
        guard let url = Bundle.main.appStoreReceiptURL,
              let data = try? Data(contentsOf: url) else { return nil }
        return data.base64EncodedString()
    }

    private func syncReceiptToBackendWithRetry(maxAttempts: Int = 6) async throws {
        var lastError: Error?

        for attempt in 0..<maxAttempts {
            if let receipt = getAppReceiptBase64(), !receipt.isEmpty {
                do {
                    _ = try await SubscriptionService.shared.purchase(receiptData: receipt)
                    UserDefaults.standard.removeObject(forKey: Self.pendingReceiptSyncKey)
                    return
                } catch {
                    lastError = error
                }
            } else {
                lastError = StoreError.receiptUnavailable
            }

            if attempt < maxAttempts - 1 {
                try await Task.sleep(for: .milliseconds(500 * (attempt + 1)))
            }
        }

        throw lastError ?? StoreError.receiptUnavailable
    }

    private func syncPendingPurchaseIfNeeded() async {
        guard UserDefaults.standard.bool(forKey: Self.pendingReceiptSyncKey) else { return }
        try? await syncReceiptToBackendWithRetry(maxAttempts: 3)
    }

    private func restorePurchases() async {
        do {
            try await AppStore.sync()
            UserDefaults.standard.set(true, forKey: Self.pendingReceiptSyncKey)
            try await syncReceiptToBackendWithRetry()
            onPurchaseSuccess()
            dismiss()
        } catch {
            errorMessage = language.localized(
                "Satın alma bilgileri geri yüklenemedi. Lütfen tekrar deneyin.",
                "Purchase information could not be restored. Please try again."
            )
        }
    }
}

private enum StoreError: LocalizedError {
    case failedVerification
    case receiptUnavailable

    var errorDescription: String? {
        switch self {
        case .failedVerification:
            return "The App Store transaction could not be verified."
        case .receiptUnavailable:
            return "The App Store receipt is not available yet. Please try restoring your purchases."
        }
    }
}

// MARK: - Stackable Card (kart üstünde kart)

private struct StackableSubscriptionCard: View {
    let package: SubscriptionPackageDTO
    let monthlyProduct: Product?
    let yearlyProduct: Product?
    let language: Language
    let isExpanded: Bool
    let isLast: Bool
    let isPurchasing: Bool
    let onPurchase: (Product) -> Void
    let onTap: () -> Void

    private var yearlySavePercent: Int {
        guard let monthlyProduct, let yearlyProduct else { return 0 }
        let monthly12 = NSDecimalNumber(decimal: monthlyProduct.price).doubleValue * 12
        let yearlyPrice = NSDecimalNumber(decimal: yearlyProduct.price).doubleValue
        guard monthly12 > 0 else { return 0 }
        let save = ((monthly12 - yearlyPrice) / monthly12) * 100
        return max(0, Int(save.rounded()))
    }

    private var cardStyle: CardStyle {
        if isExpanded && isLast { return .highlighted }
        if !isExpanded { return .collapsed }
        return .expanded
    }

    private enum CardStyle {
        case collapsed
        case expanded
        case highlighted
    }

    var body: some View {
        VStack(spacing: 0) {
            collapsedHeader
                .contentShape(Rectangle())
                .onTapGesture { onTap() }

            if isExpanded {
                expandedContent
            }
        }
        .frame(maxWidth: .infinity, alignment: .topLeading)
        .background(backgroundView)
        .clipShape(RoundedRectangle(cornerRadius: cardCornerRadius))
        .overlay(
            RoundedRectangle(cornerRadius: cardCornerRadius)
                .stroke(borderColor, lineWidth: isExpanded && isLast ? 1.5 : 1)
        )
        .shadow(color: .black.opacity(isExpanded ? 0.3 : 0.15), radius: isExpanded ? 16 : 8, y: 4)
    }

    private var collapsedHeader: some View {
        HStack {
            Text(language == .turkish ? package.name : package.nameEn)
                .font(.system(size: 17, weight: .bold))
                .foregroundColor(headerTextColor)
                .lineLimit(1)

            Spacer(minLength: 12)

            Text(monthlyProduct?.displayPrice ?? "—")
                .font(.system(size: 16, weight: .bold))
                .foregroundColor(headerTextColor)
            Text("/ \(language.localized("ay", "mo"))")
                .font(.system(size: 13, weight: .medium))
                .foregroundColor(headerTextColor.opacity(0.85))
        }
        .padding(.horizontal, 20)
        .frame(height: collapsedCardHeight)
        .contentShape(Rectangle())
    }

    @ViewBuilder
    private var expandedContent: some View {
        VStack(alignment: .leading, spacing: 0) {
            if isLast {
                Text(language.localized("Önerilen", "Recommended"))
                    .font(.system(size: 13, weight: .bold))
                    .foregroundColor(.black)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 4)
                    .background(Theme.neonPrimary)
                    .cornerRadius(6)
                    .padding(.horizontal, 20)
                    .padding(.bottom, 10)
            }

            // Pakete özel özellikler (sadece bu kartta)
            VStack(alignment: .leading, spacing: 8) {
                HStack(spacing: 8) {
                    featureBadge(
                        package.dailyQuestionLimit.map { "\($0)" } ?? "∞",
                        suffix: language.localized("soru/gün", "q/day"),
                        isUnlimited: package.dailyQuestionLimit == nil
                    )
                    Text("•")
                        .font(.caption)
                        .foregroundColor(detailTextColor)
                    featureBadge(
                        package.dailyTTSLimit.map { "\($0)" } ?? "∞",
                        suffix: language.localized("TTS/gün", "TTS/day"),
                        isUnlimited: package.dailyTTSLimit == nil
                    )
                }
            }
            .padding(.horizontal, 20)
            .padding(.bottom, 16)

            VStack(spacing: 10) {
                chooseRow(
                    product: monthlyProduct,
                    period: language.localized("ay", "mo"),
                    isPrimary: false
                )
                chooseRow(
                    product: yearlyProduct,
                    period: language.localized("yıl", "yr"),
                    isPrimary: true,
                    savePercent: yearlySavePercent
                )
            }
            .padding(.horizontal, 20)
            .padding(.bottom, 20)
        }
    }

    private var headerTextColor: Color {
        if isExpanded { return Theme.textPrimary }
        return package.slug.lowercased() == "bronze" ? .black : Theme.textPrimary
    }

    private var detailTextColor: Color {
        isLast ? Theme.textMuted : Theme.textMuted
    }

    private var backgroundColor: Color {
        if isExpanded && isLast {
            return Theme.cardBackground
        }
        if isExpanded {
            return Theme.glassColor
        }
        if package.slug.lowercased() == "bronze" {
            return Theme.neonPrimary.opacity(0.95)
        }
        return Theme.cardBackground
    }

    private var borderColor: Color {
        if isExpanded && isLast {
            return Theme.neonPrimary.opacity(0.4)
        }
        return Theme.border.opacity(0.5)
    }

    private var backgroundView: some View {
        RoundedRectangle(cornerRadius: cardCornerRadius)
            .fill(backgroundColor)
    }

    private func featureBadge(_ value: String, suffix: String, isUnlimited: Bool) -> some View {
        Text("\(value) \(suffix)")
            .font(.system(size: 13, weight: .medium))
            .foregroundColor(isUnlimited ? Theme.neonPrimary : detailTextColor)
    }

    private func chooseRow(product: Product?, period: String, isPrimary: Bool, savePercent: Int = 0) -> some View {
        Button {
            guard let product else { return }
            onPurchase(product)
        } label: {
            HStack(spacing: 12) {
                VStack(alignment: .leading, spacing: 2) {
                    if isPurchasing {
                        ProgressView()
                            .tint(chooseButtonTextColor(isPrimary: isPrimary))
                            .scaleEffect(0.85)
                    } else {
                        Text(product?.displayPrice ?? language.localized("Kullanılamıyor", "Unavailable"))
                            .font(.system(size: 17, weight: .bold))
                            .foregroundColor(isPrimary ? Theme.neonPrimary : Theme.textPrimary)
                        HStack(spacing: 6) {
                            Text("/ \(period)")
                                .font(.system(size: 13, weight: .medium))
                                .foregroundColor((isPrimary ? Theme.neonPrimary : Theme.textPrimary).opacity(0.8))
                            if savePercent > 0 && isPrimary {
                                Text(language.localized("%\(savePercent) tasarruf", "Save %\(savePercent)"))
                                    .font(.system(size: 13, weight: .bold))
                                    .foregroundColor(Theme.neonPrimary)
                            }
                        }
                    }
                }

                Spacer(minLength: 8)

                Text(product == nil
                     ? language.localized("Yok", "Unavailable")
                     : language.localized("Seç", "Choose"))
                    .font(.system(size: 14, weight: .bold))
                    .foregroundColor(chooseButtonTextColor(isPrimary: isPrimary))
                    .padding(.horizontal, 16)
                    .padding(.vertical, 10)
                    .background(chooseButtonBackground(isPrimary: isPrimary))
                    .cornerRadius(12)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .frame(maxWidth: .infinity)
            .background(chooseRowBackground(isPrimary: isPrimary))
            .cornerRadius(14)
        }
        .buttonStyle(.plain)
        .disabled(isPurchasing || product == nil)
        .opacity(product == nil ? 0.6 : 1)
    }

    private func chooseButtonTextColor(isPrimary: Bool) -> Color {
        isPrimary ? .black : Theme.textPrimary
    }

    private func chooseButtonBackground(isPrimary: Bool) -> Color {
        isPrimary ? Theme.neonPrimary : Theme.glassColor
    }

    private func chooseRowBackground(isPrimary: Bool) -> Color {
        Theme.glassColor.opacity(isPrimary ? 1 : 0.6)
    }

}

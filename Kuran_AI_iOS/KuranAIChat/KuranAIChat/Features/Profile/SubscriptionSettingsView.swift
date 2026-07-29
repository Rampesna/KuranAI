import SwiftUI

struct SubscriptionSettingsView: View {
    let language: Language
    let onBack: () -> Void

    @State private var subscriptionInfo: SubscriptionInfo?
    @State private var usageInfo: UsageInfo?
    @State private var packages: [SubscriptionPackageDTO] = []
    @State private var isLoading = true
    @State private var errorMessage: String?
    @State private var showUpgradeSheet = false

    private var roundedBg: some View {
        RoundedRectangle(cornerRadius: 10).fill(Theme.cardBackground)
    }

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            if isLoading {
                ProgressView().tint(Theme.textPrimary)
            } else {
                ScrollView {
                    VStack(alignment: .leading, spacing: 24) {
                        currentPackageCard
                        usageCard
                        openPackagesButton
                    }
                    .padding(20)
                }
            }
        }
        .navigationTitle(language.localized("Abonelik ayarları", "Subscription settings"))
        .navigationBarTitleDisplayMode(.inline)
        .task { await load() }
        .sheet(isPresented: $showUpgradeSheet) {
            SubscriptionUpgradeView(
                language: language,
                limitInfo: nil,
                onDismiss: { showUpgradeSheet = false },
                onPurchaseSuccess: {
                    showUpgradeSheet = false
                    Task { await load() }
                }
            )
        }
    }

    private var currentPackageCard: some View {
        Button {
            showUpgradeSheet = true
        } label: {
            VStack(alignment: .leading, spacing: 12) {
                HStack {
                    Image(systemName: "crown.fill")
                        .font(.system(size: 16))
                        .foregroundColor(Theme.neonPrimary)
                    Text(language.localized("Mevcut paket", "Current package"))
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundColor(Theme.textPrimary)
                }

                if let sub = subscriptionInfo, sub.hasPackage, let pkg = sub.package {
                    HStack {
                        Text(language.localized(pkg.name, pkg.nameEn))
                            .font(.system(size: 17, weight: .semibold))
                            .foregroundColor(Theme.textPrimary)
                        Spacer()
                        if let expires = sub.expiresAt, !expires.isEmpty {
                            Text(language.localized("Bitiş", "Expires") + ": \(formatDate(expires))")
                                .font(.system(size: 13))
                                .foregroundColor(Theme.textMuted)
                        }
                        Image(systemName: "chevron.right")
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundColor(Theme.textMuted)
                    }
                    .padding(.horizontal, 16)
                    .padding(.vertical, 14)
                    .frame(minHeight: 48)
                    .background(roundedBg)
                    .cornerRadius(10)
                } else {
                    HStack {
                        Text(language.localized("Ücretsiz", "Free"))
                            .font(.system(size: 17, weight: .semibold))
                            .foregroundColor(Theme.textMuted)
                        Spacer()
                        Text(language.localized("Paket seçmek için dokun", "Tap to choose package"))
                            .font(.system(size: 13))
                            .foregroundColor(Theme.neonPrimary.opacity(0.9))
                        Image(systemName: "chevron.right")
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundColor(Theme.textMuted)
                    }
                    .padding(.horizontal, 16)
                    .padding(.vertical, 14)
                    .frame(minHeight: 48)
                    .background(roundedBg)
                    .cornerRadius(10)
                }
            }
        }
        .buttonStyle(.plain)
    }

    private var usageCard: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(language.localized("Kullanım", "Usage"))
                .font(.system(size: 14, weight: .semibold))
                .foregroundColor(Theme.textPrimary)

            if let usage = usageInfo {
                usageRow(
                    title: language.localized("Soru", "Questions"),
                    used: usage.messageUsed,
                    limit: usage.messageLimit
                )
                usageRow(
                    title: language.localized("Sesli okuma", "TTS"),
                    used: usage.ttsUsed,
                    limit: usage.ttsLimit
                )
            }
        }
    }

    private func usageRow(title: String, used: Int, limit: Int) -> some View {
        HStack {
            Text(title)
                .font(.system(size: 14))
                .foregroundColor(Theme.textMuted)
            Spacer()
            Text("\(used) / \(limit == -1 ? "∞" : "\(limit)")")
                .font(.system(size: 13, weight: .medium))
                .foregroundColor(Theme.textPrimary)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 14)
        .frame(minHeight: 48)
        .background(roundedBg)
        .cornerRadius(10)
    }

    private var openPackagesButton: some View {
        Button {
            showUpgradeSheet = true
        } label: {
            HStack(spacing: 14) {
                ZStack {
                    RoundedRectangle(cornerRadius: 12)
                        .fill(Theme.neonPrimary.opacity(0.2))
                        .frame(width: 48, height: 48)
                    Image(systemName: "sparkles")
                        .font(.system(size: 22, weight: .medium))
                        .foregroundColor(Theme.neonPrimary)
                }
                VStack(alignment: .leading, spacing: 4) {
                    Text(language.localized("Paketleri görüntüle", "View packages"))
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundColor(Theme.textPrimary)
                    Text(language.localized("Seç ve satın al", "Select and purchase"))
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textMuted)
                }
                Spacer()
                Image(systemName: "chevron.right")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundColor(Theme.textMuted)
            }
            .padding(16)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Theme.glassColor)
                    .overlay(
                        RoundedRectangle(cornerRadius: 16)
                            .stroke(Theme.neonPrimary.opacity(0.25), lineWidth: 1)
                    )
            )
        }
        .buttonStyle(.plain)
    }

    private func formatDate(_ iso: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: iso) ?? ISO8601DateFormatter().date(from: iso) {
            let out = DateFormatter()
            out.dateStyle = .medium
            out.locale = Locale(identifier: language.localeIdentifier)
            return out.string(from: date)
        }
        return iso
    }

    private func load() async {
        isLoading = true
        errorMessage = nil
        defer { isLoading = false }
        do {
            let (sub, usage) = try await SubscriptionService.shared.fetchMe()
            subscriptionInfo = sub
            usageInfo = usage
            packages = (try? await SubscriptionService.shared.fetchPackages()) ?? []
        } catch {
            errorMessage = error.localizedDescription
            packages = (try? await SubscriptionService.shared.fetchPackages()) ?? []
        }
    }
}

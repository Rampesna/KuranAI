import Foundation

final class SubscriptionService {
    static let shared = SubscriptionService()

    private init() {}

    func fetchMe() async throws -> (subscription: SubscriptionInfo, usage: UsageInfo) {
        let response: SubscriptionMeResponse = try await APIClient.shared.request(.subscriptionMe)
        guard response.isSuccess, let data = response.data else {
            throw APIError.serverError(response.message ?? "Invalid response")
        }
        return (data.subscription, data.usage)
    }

    func fetchPackages() async throws -> [SubscriptionPackageDTO] {
        let response: SubscriptionPackagesResponse = try await APIClient.shared.request(.subscriptionPackages)
        guard response.isSuccess, let data = response.data else {
            throw APIError.serverError(response.message ?? "Invalid response")
        }
        return data.packages
    }

    func purchase(receiptData: String) async throws -> PurchaseSubscriptionInfo {
        let body = PurchaseRequest(receiptData: receiptData)
        let response: PurchaseResponse = try await APIClient.shared.request(.subscriptionPurchase, body: body)
        guard response.isSuccess, let data = response.data, let sub = data.subscription else {
            throw APIError.serverError(response.message ?? "Purchase failed")
        }
        return sub
    }
}

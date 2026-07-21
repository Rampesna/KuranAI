import Foundation

struct SubscriptionMeResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: SubscriptionMeData?
}

struct SubscriptionMeData: Decodable {
    let subscription: SubscriptionInfo
    let usage: UsageInfo
}

struct SubscriptionInfo: Decodable {
    let hasPackage: Bool
    let packageId: UInt?
    let package: SubscriptionPackageDTO?
    let expiresAt: String?
    let isMonthly: Bool?

    enum CodingKeys: String, CodingKey {
        case hasPackage, packageId, package, expiresAt, isMonthly
    }
}

struct UsageInfo: Decodable {
    let messageUsed: Int
    let messageLimit: Int
    let ttsUsed: Int
    let ttsLimit: Int

    enum CodingKeys: String, CodingKey {
        case messageUsed, messageLimit, ttsUsed, ttsLimit
    }
}

struct SubscriptionPackageDTO: Decodable, Identifiable {
    let id: UInt
    let name: String
    let nameEn: String
    let slug: String
    let dailyQuestionLimit: Int?
    let dailyTTSLimit: Int?
    let monthlyPriceTL: Double
    let yearlyPriceTL: Double

    enum CodingKeys: String, CodingKey {
        case id, name, nameEn, slug, dailyQuestionLimit, dailyTTSLimit, monthlyPriceTL, yearlyPriceTL
    }
}

struct SubscriptionPackagesResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: SubscriptionPackagesData?
}

struct SubscriptionPackagesData: Decodable {
    let packages: [SubscriptionPackageDTO]
}

struct PurchaseRequest: Encodable {
    let receiptData: String
}

struct PurchaseResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: PurchaseData?
}

struct PurchaseData: Decodable {
    let subscription: PurchaseSubscriptionInfo?
}

struct PurchaseSubscriptionInfo: Decodable {
    let id: UInt
    let packageId: UInt
    let expiresAt: String
}

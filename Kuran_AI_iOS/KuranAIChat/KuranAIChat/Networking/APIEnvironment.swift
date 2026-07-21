import Foundation

enum APIEnvironment: Sendable {
    case production

    private static let productionBaseURL = "https://kuranai-api.talhacan.com"

    nonisolated static var current: APIEnvironment { .production }

    nonisolated var baseURL: String { Self.productionBaseURL }

    nonisolated var timeout: TimeInterval { 30 }
}

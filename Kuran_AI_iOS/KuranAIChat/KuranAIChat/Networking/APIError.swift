import Foundation

enum APIError: Error, LocalizedError {
    case invalidURL
    case invalidResponse
    case httpError(statusCode: Int, data: Data?)
    case decodingError(Error)
    case networkError(Error)
    case unauthorized
    case serverError(String)
    case limitExceeded(LimitExceededInfo)

    struct LimitExceededInfo: Equatable {
        let code: String
        let limitType: String
        let used: Int
        let limit: Int
        let packageId: UInt?
        let message: String
    }

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .invalidResponse:
            return "Invalid response from server"
        case .httpError(let code, _):
            return "HTTP error: \(code)"
        case .decodingError(let error):
            return "Decoding failed: \(error.localizedDescription)"
        case .networkError(let error):
            return error.localizedDescription
        case .unauthorized:
            return "Unauthorized"
        case .serverError(let message):
            return message
        case .limitExceeded(let info):
            return info.message
        }
    }
}

import Foundation

// Backend wrapper: { "isSuccess", "message", "statusCode", "data" }
struct ServiceResponse<T: Decodable>: Decodable {
    let isSuccess: Bool
    let message: String
    let statusCode: Int
    let data: T?
}

// Auth payload from backend login/google and login/apple
struct AuthData: Decodable {
    let accessToken: String
    let refreshToken: String
    let expiresAt: String
    let tokenType: String
}

struct LoginResponse: Decodable {
    let token: String
    let refreshToken: String?
    let user: UserDTO?
}

struct UserDTO: Decodable, Equatable {
    let id: Int
    let email: String?
    let name: String?
    let username: String?
    let avatarUrl: String?
    let hasPassword: Bool?
}

/// Backend data: null yanıtları için placeholder
struct EmptyData: Decodable {}

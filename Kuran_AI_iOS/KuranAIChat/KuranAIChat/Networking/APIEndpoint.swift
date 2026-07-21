import Foundation

enum APIEndpoint: Sendable {
    // MARK: - Auth
    case login
    case register
    case loginGoogle
    case loginApple
    case logout
    case refreshToken

    // MARK: - Chat
    case sendMessage
    case transcribe
    case tts
    case chatSessions
    case chatMessages(sessionId: UInt)

    // MARK: - Subscription
    case subscriptionMe
    case subscriptionPackages
    case subscriptionPurchase

    // MARK: - Prayer
    case prayerTimes(lat: Double?, lng: Double?)

    // MARK: - Verse
    case verseOfDay(lang: String)

    // MARK: - Future
    case chatHistory
    case userProfile
    case me
    case updateUsername
    case updatePassword
    case uploadAvatar
    case deleteAvatar
    case deleteAccount

    var path: String {
        switch self {
        case .login: return "/api/user/auth/login"
        case .register: return "/api/user/auth/register"
        case .loginGoogle: return "/api/user/auth/login/google"
        case .loginApple: return "/api/user/auth/login/apple"
        case .logout: return "/api/user/auth/logout"
        case .refreshToken: return "/api/user/auth/refresh"
        case .sendMessage: return "/api/v1/chat/message"
        case .transcribe: return "/api/v1/chat/transcribe"
        case .tts: return "/api/v1/chat/tts"
        case .chatSessions: return "/api/v1/chat/sessions"
        case .chatMessages(let id): return "/api/v1/chat/sessions/\(id)/messages"
        case .subscriptionMe: return "/api/v1/subscription/me"
        case .subscriptionPackages: return "/api/v1/subscription/packages"
        case .subscriptionPurchase: return "/api/v1/subscription/purchase"
        case .prayerTimes: return "/api/v1/prayer-times"
        case .verseOfDay: return "/api/verse-of-day"
        case .chatHistory: return "/api/v1/chat/history"
        case .userProfile: return "/api/v1/user/profile"
        case .me: return "/api/user/auth/me"
        case .updateUsername: return "/api/user/profile/username"
        case .updatePassword: return "/api/user/profile/password"
        case .uploadAvatar: return "/api/user/profile/avatar"
        case .deleteAvatar: return "/api/user/profile/avatar"
        case .deleteAccount: return "/api/user/profile/account"
        }
    }

    var method: HTTPMethod {
        switch self {
        case .login, .register, .loginGoogle, .loginApple, .sendMessage, .logout, .uploadAvatar, .transcribe, .tts, .subscriptionPurchase: return .post
        case .updateUsername, .updatePassword: return .patch
        case .chatSessions, .chatMessages, .subscriptionMe, .subscriptionPackages, .refreshToken, .prayerTimes, .verseOfDay, .chatHistory, .userProfile, .me: return .get
        case .deleteAvatar, .deleteAccount: return .delete
        }
    }

    var needsAuth: Bool {
        switch self {
        case .login, .register, .loginGoogle, .loginApple, .prayerTimes, .verseOfDay: return false
        default: return true
        }
    }
}

enum HTTPMethod: String, Sendable {
    case get = "GET"
    case post = "POST"
    case put = "PUT"
    case patch = "PATCH"
    case delete = "DELETE"
}

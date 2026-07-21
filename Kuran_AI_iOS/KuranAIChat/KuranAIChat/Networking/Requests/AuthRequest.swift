import Foundation

// Backend expects "idToken" (camelCase); encoder uses snake_case so we force keys.
struct GoogleLoginRequest: Encodable {
    let idToken: String

    enum CodingKeys: String, CodingKey {
        case idToken = "idToken"
    }
}

struct AppleLoginRequest: Encodable {
    let identityToken: String

    enum CodingKeys: String, CodingKey {
        case identityToken = "identityToken"
    }
}

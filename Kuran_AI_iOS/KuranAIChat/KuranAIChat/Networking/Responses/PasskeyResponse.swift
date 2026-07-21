import Foundation

struct PasskeyChallengeData: Decodable {
    let transactionId: String
    let options: PasskeyPublicKeyOptions
}

struct PasskeyPublicKeyOptions: Decodable {
    let challenge: String
    let rpId: String?
    let user: PasskeyUserOptions?
}

struct PasskeyUserOptions: Decodable {
    let id: String
    let name: String
    let displayName: String
}

struct PasskeyStatusData: Decodable {
    let count: Int
    let enabled: Bool
}

struct PasskeyRegistrationCredential: Encodable {
    let id: String
    let rawId: String
    let type = "public-key"
    let authenticatorAttachment = "platform"
    let clientExtensionResults: [String: String] = [:]
    let response: RegistrationResponse

    struct RegistrationResponse: Encodable {
        let attestationObject: String
        let clientDataJSON: String
        let transports = ["internal"]
    }
}

struct PasskeyAssertionCredential: Encodable {
    let id: String
    let rawId: String
    let type = "public-key"
    let authenticatorAttachment = "platform"
    let clientExtensionResults: [String: String] = [:]
    let response: AssertionResponse

    struct AssertionResponse: Encodable {
        let authenticatorData: String
        let clientDataJSON: String
        let signature: String
        let userHandle: String?
    }
}

extension Data {
    init?(base64URLEncoded value: String) {
        var base64 = value.replacingOccurrences(of: "-", with: "+")
            .replacingOccurrences(of: "_", with: "/")
        base64.append(String(repeating: "=", count: (4 - base64.count % 4) % 4))
        self.init(base64Encoded: base64)
    }

    var base64URLEncodedString: String {
        base64EncodedString()
            .replacingOccurrences(of: "+", with: "-")
            .replacingOccurrences(of: "/", with: "_")
            .replacingOccurrences(of: "=", with: "")
    }
}

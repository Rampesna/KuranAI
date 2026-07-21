import Foundation

actor APIClient {
    static let shared = APIClient()
    private let session: URLSession
    private let baseURL: String
    private let timeout: TimeInterval
    private let decoder: JSONDecoder
    private let encoder: JSONEncoder

    private(set) var authToken: String?

    /// Aynı anda birden fazla 401 olunca tek bir refresh yapılmasını sağlar.
    private var isRefreshing = false
    private var refreshWaiters: [CheckedContinuation<Bool, Never>] = []

    func setAuthToken(_ token: String?) {
        authToken = token
    }

    private init(
        session: URLSession = .shared,
        baseURL: String = APIEnvironment.current.baseURL,
        timeout: TimeInterval = APIEnvironment.current.timeout,
        decoder: JSONDecoder = JSONDecoder(),
        encoder: JSONEncoder = JSONEncoder()
    ) {
        self.timeout = timeout
        self.session = session
        self.baseURL = baseURL
        self.decoder = decoder
        self.encoder = encoder
        self.decoder.keyDecodingStrategy = .convertFromSnakeCase
        self.encoder.keyEncodingStrategy = .convertToSnakeCase
    }

    // MARK: - Public Request Methods

    func request<T: Decodable>(
        _ endpoint: APIEndpoint,
        body: Encodable? = nil,
        queryItems: [URLQueryItem]? = nil
    ) async throws -> T {
        let data = try await performRequestWithRetry(endpoint, body: body, queryItems: queryItems)
        return try decoder.decode(T.self, from: data)
    }

    func request(_ endpoint: APIEndpoint, body: Encodable? = nil) async throws {
        _ = try await performRequestWithRetry(endpoint, body: body, queryItems: nil)
    }

    // MARK: - Request with Auto-Retry on 401

    private func performRequestWithRetry(
        _ endpoint: APIEndpoint,
        body: Encodable?,
        queryItems: [URLQueryItem]?
    ) async throws -> Data {
        do {
            return try await performRequest(endpoint, body: body, queryItems: queryItems)
        } catch let error as APIError {
            // Sadece auth gerektiren endpoint'lerde 401'i yakala
            guard endpoint.needsAuth, case .httpError(let code, _) = error, code == 401 else {
                throw error
            }

            // Refresh dene
            let refreshed = await coordinatedRefresh()
            guard refreshed else {
                throw APIError.unauthorized
            }

            // Yeni token ile tekrar dene
            return try await performRequest(endpoint, body: body, queryItems: queryItems)
        }
    }

    /// Aynı anda birden fazla 401 gelirse hepsini tek bir refresh ile çözer.
    private func coordinatedRefresh() async -> Bool {
        if isRefreshing {
            // Zaten refresh yapılıyor, bekle
            return await withCheckedContinuation { continuation in
                refreshWaiters.append(continuation)
            }
        }

        isRefreshing = true
        let success = await AuthService.shared.refreshTokens()
        isRefreshing = false

        // Bekleyen tüm isteklere sonucu bildir
        let waiters = refreshWaiters
        refreshWaiters.removeAll()
        for waiter in waiters {
            waiter.resume(returning: success)
        }

        return success
    }

    // MARK: - Core Request

    private func performRequest(
        _ endpoint: APIEndpoint,
        body: Encodable?,
        queryItems: [URLQueryItem]?
    ) async throws -> Data {
        guard var components = URLComponents(string: baseURL + endpoint.path) else {
            throw APIError.invalidURL
        }
        components.queryItems = queryItems

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = endpoint.method.rawValue
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.timeoutInterval = timeout

        if endpoint.needsAuth, let token = authToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body = body {
            request.httpBody = try encoder.encode(AnyEncodable(body))
        }

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }

        return data
    }

    // MARK: - Refresh Token Request (AuthService tarafından çağrılır)

    func performRefreshToken(refreshToken: String) async throws -> AuthData {
        guard let url = URL(string: baseURL + APIEndpoint.refreshToken.path) else {
            throw APIError.invalidURL
        }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.timeoutInterval = timeout

        // Backend camelCase "refreshToken" bekliyor — snake_case encoder kullanılmamalı
        let camelEncoder = JSONEncoder()
        camelEncoder.keyEncodingStrategy = .useDefaultKeys
        struct RefreshReq: Encodable { let refreshToken: String }
        request.httpBody = try camelEncoder.encode(RefreshReq(refreshToken: refreshToken))

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }
        let wrapper = try decoder.decode(ServiceResponse<AuthData>.self, from: data)
        guard wrapper.isSuccess, let authData = wrapper.data else {
            throw APIError.serverError(wrapper.message)
        }
        return authData
    }

    // MARK: - Auth Endpoints (camelCase body)

    func loginWithGoogle(idToken: String) async throws -> AuthData {
        let body = ["idToken": idToken]
        return try await performAuthRequest(path: APIEndpoint.loginGoogle.path, body: body)
    }

    func loginWithApple(identityToken: String) async throws -> AuthData {
        let body = ["identityToken": identityToken]
        return try await performAuthRequest(path: APIEndpoint.loginApple.path, body: body)
    }

    func loginWithEmail(credential: String, password: String) async throws -> AuthData {
        struct Req: Encodable { let credential, password: String }
        return try await performJSONAuthRequest(path: APIEndpoint.login.path, body: Req(credential: credential, password: password))
    }

    func register(email: String, password: String, name: String, username: String) async throws {
        struct Req: Encodable { let email, password, name, username: String }
        let body = Req(email: email, password: password, name: name, username: username)
        _ = try await performRegisterRequest(path: APIEndpoint.register.path, body: body)
    }

    // MARK: - Passkeys (WebAuthn)

    func beginPasskeyLogin() async throws -> PasskeyChallengeData {
        let wrapper: ServiceResponse<PasskeyChallengeData> = try await request(.passkeyLoginOptions)
        guard wrapper.isSuccess, let data = wrapper.data else { throw APIError.serverError(wrapper.message) }
        return data
    }

    func finishPasskeyLogin(transactionId: String, credential: PasskeyAssertionCredential) async throws -> AuthData {
        let data = try await performPasskeyRequest(.passkeyLoginFinish, transactionId: transactionId, body: credential)
        let wrapper = try decoder.decode(ServiceResponse<AuthData>.self, from: data)
        guard wrapper.isSuccess, let auth = wrapper.data else { throw APIError.serverError(wrapper.message) }
        return auth
    }

    func beginPasskeyRegistration() async throws -> PasskeyChallengeData {
        let wrapper: ServiceResponse<PasskeyChallengeData> = try await request(.passkeyRegistrationOptions)
        guard wrapper.isSuccess, let data = wrapper.data else { throw APIError.serverError(wrapper.message) }
        return data
    }

    func finishPasskeyRegistration(transactionId: String, credential: PasskeyRegistrationCredential) async throws {
        _ = try await performPasskeyRequest(.passkeyRegistrationFinish, transactionId: transactionId, body: credential)
    }

    func passkeyStatus() async throws -> PasskeyStatusData {
        let wrapper: ServiceResponse<PasskeyStatusData> = try await request(.passkeyStatus)
        guard wrapper.isSuccess, let data = wrapper.data else { throw APIError.serverError(wrapper.message) }
        return data
    }

    func deletePasskeys() async throws {
        try await request(.passkeyDelete)
    }

    private func performPasskeyRequest<T: Encodable>(
        _ endpoint: APIEndpoint,
        transactionId: String,
        body: T,
        retryAfterRefresh: Bool = true
    ) async throws -> Data {
        guard let url = URL(string: baseURL + endpoint.path) else { throw APIError.invalidURL }
        var request = URLRequest(url: url)
        request.httpMethod = endpoint.method.rawValue
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.setValue(transactionId, forHTTPHeaderField: "X-Passkey-Transaction")
        request.setValue("ios", forHTTPHeaderField: "X-Device-Type")
        request.timeoutInterval = timeout
        if endpoint.needsAuth, let token = authToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        let passkeyEncoder = JSONEncoder()
        passkeyEncoder.keyEncodingStrategy = .useDefaultKeys
        request.httpBody = try passkeyEncoder.encode(body)

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else { throw APIError.invalidResponse }
        if httpResponse.statusCode == 401, endpoint.needsAuth, retryAfterRefresh, await coordinatedRefresh() {
            return try await performPasskeyRequest(endpoint, transactionId: transactionId, body: body, retryAfterRefresh: false)
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            if let wrapper = try? decoder.decode(ServiceResponseMessage.self, from: data) {
                throw APIError.serverError(wrapper.message)
            }
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }
        return data
    }

    private func performRegisterRequest(path: String, body: some Encodable) async throws -> Data {
        guard let url = URL(string: baseURL + path) else { throw APIError.invalidURL }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.timeoutInterval = APIEnvironment.current.timeout
        request.httpBody = try encoder.encode(body)

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else { throw APIError.invalidResponse }
        guard (200...299).contains(httpResponse.statusCode) else {
            if let w = try? decoder.decode(ServiceResponseMessage.self, from: data) {
                throw APIError.serverError(w.message)
            }
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }
        return data
    }

    private func performJSONAuthRequest<T: Encodable>(path: String, body: T) async throws -> AuthData {
        guard let url = URL(string: baseURL + path) else { throw APIError.invalidURL }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.timeoutInterval = APIEnvironment.current.timeout
        request.httpBody = try encoder.encode(body)

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else { throw APIError.invalidResponse }

        let wrapper = try? decoder.decode(ServiceResponse<AuthData>.self, from: data)
        if let w = wrapper, !w.isSuccess {
            throw APIError.serverError(w.message)
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            if let w = try? decoder.decode(ServiceResponseMessage.self, from: data) {
                throw APIError.serverError(w.message)
            }
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }
        guard let authData = wrapper?.data else {
            throw APIError.decodingError(NSError(domain: "Auth", code: -1, userInfo: [NSLocalizedDescriptionKey: "Missing auth data"]))
        }
        return authData
    }

    private func performAuthRequest(path: String, body: [String: String]) async throws -> AuthData {
        guard let url = URL(string: baseURL + path) else { throw APIError.invalidURL }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.timeoutInterval = APIEnvironment.current.timeout
        let camelEncoder = JSONEncoder()
        camelEncoder.keyEncodingStrategy = .useDefaultKeys
        request.httpBody = try camelEncoder.encode(body)

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else { throw APIError.invalidResponse }

        let wrapper = try? decoder.decode(ServiceResponse<AuthData>.self, from: data)
        if let w = wrapper, !w.isSuccess {
            throw APIError.serverError(w.message)
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }
        guard let authData = wrapper?.data else {
            throw APIError.decodingError(NSError(domain: "Auth", code: -1, userInfo: [NSLocalizedDescriptionKey: "Missing auth data"]))
        }
        return authData
    }
}

private struct ServiceResponseMessage: Decodable {
    let isSuccess: Bool
    let message: String
}

private struct AnyEncodable: Encodable {
    private let encode: (Encoder) throws -> Void

    init<T: Encodable>(_ value: T) {
        encode = value.encode
    }

    func encode(to encoder: Encoder) throws {
        try encode(encoder)
    }
}

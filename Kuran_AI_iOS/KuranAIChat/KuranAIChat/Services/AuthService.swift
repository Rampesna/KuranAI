import Foundation
import UIKit
import Combine
import AuthenticationServices
#if canImport(GoogleSignIn)
import GoogleSignIn
#endif

/// Apple + Google sign-in; calls backend login/google or login/apple and stores token.
@MainActor
final class AuthService: NSObject, ObservableObject {
    static let shared = AuthService()

    /// Oturum geçersiz olduğunda (refresh de başarısız) gönderilir; ContentView dinler.
    static let sessionExpiredNotification = Notification.Name("KuranAI_SessionExpired")

    private let accessTokenKey = "kuran_ai_access_token"
    private let refreshTokenKey = "kuran_ai_refresh_token"

    private var appleContinuation: CheckedContinuation<String, Error>?

    override private init() {
        super.init()
    }

    /// Restore stored token and set on APIClient (call at app launch).
    func restoreSession() async {
        let token = UserDefaults.standard.string(forKey: accessTokenKey)
        await APIClient.shared.setAuthToken(token)
    }

    var isLoggedIn: Bool {
        UserDefaults.standard.string(forKey: accessTokenKey) != nil
    }

    var refreshToken: String? {
        UserDefaults.standard.string(forKey: refreshTokenKey)
    }

    func saveAuth(_ auth: AuthData) {
        UserDefaults.standard.set(auth.accessToken, forKey: accessTokenKey)
        UserDefaults.standard.set(auth.refreshToken, forKey: refreshTokenKey)
    }

    func clearSession() {
        UserDefaults.standard.removeObject(forKey: accessTokenKey)
        UserDefaults.standard.removeObject(forKey: refreshTokenKey)
        Task { await APIClient.shared.setAuthToken(nil) }
    }

    /// Backend'e logout isteği gönderir, sonra yerel oturumu temizler.
    func logout() async {
        do {
            try await APIClient.shared.request(APIEndpoint.logout)
        } catch { /* Ağ hatası olsa da yerelde çıkış yapılacak */ }
        clearSession()
    }

    // MARK: - Token Refresh

    /// Refresh token ile yeni access + refresh token alır.
    /// Başarısız olursa session temizlenir ve sessionExpired notification gönderilir.
    func refreshTokens() async -> Bool {
        guard let currentRefresh = UserDefaults.standard.string(forKey: refreshTokenKey) else {
            broadcastSessionExpired()
            return false
        }
        do {
            let auth = try await APIClient.shared.performRefreshToken(refreshToken: currentRefresh)
            saveAuth(auth)
            await APIClient.shared.setAuthToken(auth.accessToken)
            return true
        } catch {
            clearSession()
            broadcastSessionExpired()
            return false
        }
    }

    /// Oturum geçersiz bildirimini gönderir (ana ekran dinler).
    private func broadcastSessionExpired() {
        NotificationCenter.default.post(name: Self.sessionExpiredNotification, object: nil)
    }

    // MARK: - Email login

    func signInWithEmail(credential: String, password: String) async throws -> AuthData {
        let auth = try await APIClient.shared.loginWithEmail(credential: credential, password: password)
        saveAuth(auth)
        await APIClient.shared.setAuthToken(auth.accessToken)
        return auth
    }

    /// Kayıt sonrası otomatik login.
    func registerAndLogin(email: String, password: String, name: String, username: String) async throws -> AuthData {
        try await APIClient.shared.register(email: email, password: password, name: name, username: username)
        return try await signInWithEmail(credential: email, password: password)
    }

    // MARK: - Sign in with Apple

    func signInWithApple() async throws -> AuthData {
        let token = try await performAppleSignIn()
        let auth = try await APIClient.shared.loginWithApple(identityToken: token)
        saveAuth(auth)
        await APIClient.shared.setAuthToken(auth.accessToken)
        return auth
    }

    private func performAppleSignIn() async throws -> String {
        try await withCheckedThrowingContinuation { continuation in
            appleContinuation = continuation
            let provider = ASAuthorizationAppleIDProvider()
            let request = provider.createRequest()
            request.requestedScopes = [.fullName, .email]
            let controller = ASAuthorizationController(authorizationRequests: [request])
            controller.delegate = self
            controller.presentationContextProvider = self
            controller.performRequests()
        }
    }

    // MARK: - Sign in with Google

    func signInWithGoogle() async throws -> AuthData {
#if canImport(GoogleSignIn)
        let token = try await performGoogleSignIn()
        let auth = try await APIClient.shared.loginWithGoogle(idToken: token)
        saveAuth(auth)
        await APIClient.shared.setAuthToken(auth.accessToken)
        return auth
#else
        throw APIError.serverError("Google Sign-In not linked. Add package: https://github.com/google/GoogleSignIn-iOS")
#endif
    }

#if canImport(GoogleSignIn)
    private static func topViewControllerStatic(from base: UIViewController?) -> UIViewController? {
        guard let base = base else { return nil }
        if let nav = base as? UINavigationController, let visible = nav.visibleViewController {
            return topViewControllerStatic(from: visible)
        }
        if let presented = base.presentedViewController {
            return topViewControllerStatic(from: presented)
        }
        return base
    }

    private func performGoogleSignIn() async throws -> String {
        return try await withCheckedThrowingContinuation { continuation in
            DispatchQueue.main.async {
                guard let windowScene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
                      let rootVC = windowScene.windows.first(where: { $0.isKeyWindow })?.rootViewController ?? windowScene.windows.first?.rootViewController else {
                    continuation.resume(throwing: APIError.serverError("No window for Google Sign-In"))
                    return
                }
                let presentingVC = Self.topViewControllerStatic(from: rootVC) ?? rootVC
                GIDSignIn.sharedInstance.signIn(withPresenting: presentingVC) { result, error in
                    if let error = error {
                        continuation.resume(throwing: error)
                        return
                    }
                    guard let idToken = result?.user.idToken?.tokenString else {
                        continuation.resume(throwing: APIError.serverError("Google ID token missing"))
                        return
                    }
                    continuation.resume(returning: idToken)
                }
            }
        }
    }
#endif
}

// MARK: - Apple Sign-In delegate

extension AuthService: ASAuthorizationControllerDelegate, ASAuthorizationControllerPresentationContextProviding {
    nonisolated func authorizationController(
        controller: ASAuthorizationController,
        didCompleteWithAuthorization authorization: ASAuthorization
    ) {
        guard let credential = authorization.credential as? ASAuthorizationAppleIDCredential,
              let tokenData = credential.identityToken,
              let token = String(data: tokenData, encoding: .utf8) else {
            Task { @MainActor in
                appleContinuation?.resume(throwing: APIError.serverError("Apple identity token missing"))
                appleContinuation = nil
            }
            return
        }
        Task { @MainActor in
            appleContinuation?.resume(returning: token)
            appleContinuation = nil
        }
    }

    nonisolated func authorizationController(
        controller: ASAuthorizationController,
        didCompleteWithError error: Error
    ) {
        Task { @MainActor in
            appleContinuation?.resume(throwing: error)
            appleContinuation = nil
        }
    }

    nonisolated func presentationAnchor(for controller: ASAuthorizationController) -> ASPresentationAnchor {
        DispatchQueue.main.sync {
            MainActor.assumeIsolated {
                (UIApplication.shared.connectedScenes.first as? UIWindowScene)?.windows.first ?? UIWindow()
            }
        }
    }
}

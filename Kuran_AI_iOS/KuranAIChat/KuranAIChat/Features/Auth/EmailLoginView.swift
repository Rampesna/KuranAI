import SwiftUI

struct EmailLoginView: View {
    let language: Language
    let onSuccess: () -> Void
    let onDismiss: () -> Void
    let onShowRegister: () -> Void

    @State private var credential = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String?

    var body: some View {
        NavigationStack {
            ZStack {
                Theme.background.ignoresSafeArea()

                ScrollView(showsIndicators: false) {
                    VStack(spacing: 0) {
                        headerWithClose

                        Spacer(minLength: 24)

                        formSection

                        Spacer(minLength: 32)

                        registerPrompt
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 24)
                    .padding(.bottom, 40)
                }
                .scrollDismissesKeyboard(.interactively)
            }
            .navigationBarHidden(true)
        }
    }

    private var headerWithClose: some View {
        HStack(alignment: .top) {
            VStack(alignment: .leading, spacing: 6) {
                Text(language.localized("E-posta ile giriş", "Email login"))
                    .font(.system(size: 26, weight: .bold))
                    .foregroundColor(Theme.textPrimary)
                Text(language.localized("E-posta veya kullanıcı adı ile devam edin", "Continue with email or username"))
                    .font(.system(size: 14))
                    .foregroundColor(Theme.textMuted)
            }
            Spacer()
            Button {
                onDismiss()
            } label: {
                Image(systemName: "chevron.down")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(Theme.textMuted)
                    .frame(width: 36, height: 36)
                    .background(Theme.glassColor)
                    .clipShape(Circle())
            }
        }
        .padding(.top, 20)
        .padding(.bottom, 32)
    }

    private var formSection: some View {
        VStack(spacing: 20) {
            if let msg = errorMessage {
                HStack(spacing: 8) {
                    Image(systemName: "exclamationmark.triangle.fill")
                        .font(.system(size: 13))
                        .foregroundColor(.orange)
                    Text(msg)
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textPrimary.opacity(0.95))
                }
                .padding(14)
                .frame(maxWidth: .infinity, alignment: .leading)
                .background(Theme.cardBackground)
                .overlay(RoundedRectangle(cornerRadius: 10).stroke(Theme.border, lineWidth: 1))
                .cornerRadius(10)
            }

            AuthTextField(
                icon: "person.fill",
                placeholder: language.localized("E-posta veya kullanıcı adı", "Email or username"),
                text: $credential
            )
            .keyboardType(.emailAddress)
            .autocapitalization(.none)
            .autocorrectionDisabled()

            AuthTextField(
                icon: "lock.fill",
                placeholder: language.localized("Şifre", "Password"),
                text: $password,
                isSecure: true
            )

            Button {
                Task { await login() }
            } label: {
                HStack(spacing: 10) {
                    if isLoading {
                        ProgressView().tint(.black).scaleEffect(0.9)
                    } else {
                        Text(language.localized("Giriş Yap", "Sign in"))
                            .font(.system(size: 16, weight: .semibold))
                    }
                }
                .frame(maxWidth: .infinity)
                .padding(.vertical, 14)
                .background(Theme.neon)
                .foregroundColor(.black)
                .cornerRadius(10)
            }
            .disabled(isLoading || credential.isEmpty || password.isEmpty)
            .opacity(credential.isEmpty || password.isEmpty ? 0.6 : 1)
        }
    }

    private var registerPrompt: some View {
        VStack(spacing: 20) {
            HStack {
                Rectangle().fill(Theme.border).frame(height: 1)
                Text(language.localized("veya", "or"))
                    .font(.system(size: 13, weight: .medium))
                    .foregroundColor(Theme.textMuted)
                Rectangle().fill(Theme.border).frame(height: 1)
            }
            .padding(.top, 28)

            Button {
                onShowRegister()
            } label: {
                Text(language.localized("Hesabınız yok mu? Kayıt olun", "Don't have an account? Register"))
                    .font(.system(size: 15, weight: .medium))
                    .foregroundStyle(Theme.neon)
            }
        }
    }

    private func login() async {
        errorMessage = nil
        isLoading = true
        defer { isLoading = false }
        do {
            _ = try await AuthService.shared.signInWithEmail(credential: credential.trimmingCharacters(in: .whitespaces), password: password)
            onSuccess()
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

import SwiftUI

struct RegisterView: View {
    let language: Language
    let onSuccess: () -> Void
    let onDismiss: () -> Void

    @State private var email = ""
    @State private var password = ""
    @State private var name = ""
    @State private var username = ""
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
                Text(language.localized("Hesap oluştur", "Create account"))
                    .font(.system(size: 26, weight: .bold))
                    .foregroundColor(Theme.textPrimary)
                Text(language.localized("Bilgilerinizi girerek kayıt olun", "Register by entering your details"))
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
                icon: "envelope.fill",
                placeholder: language.localized("E-posta", "Email"),
                text: $email
            )
            .keyboardType(.emailAddress)
            .autocapitalization(.none)
            .autocorrectionDisabled()

            AuthTextField(
                icon: "lock.fill",
                placeholder: language.localized("Şifre (min. 6 karakter)", "Password (min. 6 characters)"),
                text: $password,
                isSecure: true
            )

            AuthTextField(
                icon: "person.fill",
                placeholder: language.localized("Ad Soyad", "Full name"),
                text: $name
            )
            .autocapitalization(.words)

            AuthTextField(
                icon: "at",
                placeholder: language.localized("Kullanıcı adı", "Username"),
                text: $username
            )
            .autocapitalization(.none)
            .autocorrectionDisabled()

            Button {
                Task { await register() }
            } label: {
                HStack(spacing: 10) {
                    if isLoading {
                        ProgressView().tint(.black).scaleEffect(0.9)
                    } else {
                        Text(language.localized("Kayıt Ol", "Register"))
                            .font(.system(size: 16, weight: .semibold))
                    }
                }
                .frame(maxWidth: .infinity)
                .padding(.vertical, 14)
                .background(Theme.neon)
                .foregroundColor(.black)
                .cornerRadius(10)
            }
            .disabled(isLoading || !isFormValid)
            .opacity(isFormValid ? 1 : 0.6)
        }
    }

    private var isFormValid: Bool {
        !email.trimmingCharacters(in: .whitespaces).isEmpty &&
        !password.isEmpty &&
        password.count >= 6 &&
        !name.trimmingCharacters(in: .whitespaces).isEmpty &&
        !username.trimmingCharacters(in: .whitespaces).isEmpty
    }

    private func register() async {
        errorMessage = nil
        isLoading = true
        defer { isLoading = false }
        do {
            _ = try await AuthService.shared.registerAndLogin(
                email: email.trimmingCharacters(in: .whitespaces),
                password: password,
                name: name.trimmingCharacters(in: .whitespaces),
                username: username.trimmingCharacters(in: .whitespaces)
            )
            onSuccess()
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

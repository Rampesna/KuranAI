import SwiftUI
import AuthenticationServices

private enum PasswordField { case current, new, confirm }

struct SecuritySettingsView: View {
    let language: Language
    let onBack: () -> Void
    var onAccountDeleted: (() -> Void)?

    @State private var hasPassword: Bool = true   // varsayılan true; /me'den güncellenir
    @State private var isCheckingPassword = true  // ilk /me yüklemesi
    @State private var currentPassword: String = ""
    @State private var newPassword: String = ""
    @State private var confirmPassword: String = ""
    @State private var isSaving = false
    @State private var errorMessage: String?
    @State private var successMessage: String?
    @State private var showDeleteConfirm = false
    @State private var isDeleting = false
    @State private var passkeyEnabled = false
    @State private var passkeyCount = 0
    @State private var isUpdatingPasskey = false
    @State private var showPasskeyDeleteConfirm = false
    @FocusState private var focusedField: PasswordField?

    private var roundedBg: some View {
        RoundedRectangle(cornerRadius: 10).fill(Theme.cardBackground)
    }

    /// "Kaydet" butonunun aktif olma koşulu
    private var canSave: Bool {
        if isSaving { return false }
        if newPassword.isEmpty || confirmPassword.isEmpty { return false }
        if newPassword != confirmPassword { return false }
        if newPassword.count < 6 { return false }
        // Mevcut şifresi varsa eski şifre de dolu olmalı
        if hasPassword && currentPassword.isEmpty { return false }
        return true
    }

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            if isCheckingPassword {
                ProgressView().tint(Theme.textPrimary)
            } else {
                ScrollView {
                    VStack(alignment: .leading, spacing: 20) {
                        VStack(alignment: .leading, spacing: 10) {
                            HStack {
                                Image(systemName: "person.badge.key.fill")
                                    .foregroundColor(Theme.neonPrimary)
                                VStack(alignment: .leading, spacing: 2) {
                                    Text(language.localized("Geçiş Anahtarı", "Passkey"))
                                        .font(.system(size: 16, weight: .semibold))
                                        .foregroundColor(Theme.textPrimary)
                                    Text(passkeyEnabled
                                         ? language.localized("Parolasız giriş etkin (\(passkeyCount))", "Passwordless sign-in is active (\(passkeyCount))")
                                         : language.localized("Face ID veya Touch ID ile güvenli giriş", "Secure sign-in with Face ID or Touch ID"))
                                        .font(.system(size: 12))
                                        .foregroundColor(Theme.textMuted)
                                }
                                Spacer()
                            }

                            Button {
                                Task { await createPasskey() }
                            } label: {
                                HStack {
                                    if isUpdatingPasskey {
                                        ProgressView().tint(.black).scaleEffect(0.85)
                                    } else {
                                        Text(passkeyEnabled
                                             ? language.localized("Yeni geçiş anahtarı ekle", "Add another passkey")
                                             : language.localized("Geçiş anahtarı oluştur", "Create a passkey"))
                                    }
                                }
                                .font(.system(size: 14, weight: .semibold))
                                .foregroundColor(.black)
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 11)
                                .background(Theme.neon)
                                .cornerRadius(10)
                            }
                            .disabled(isUpdatingPasskey)

                            if passkeyEnabled {
                                Button(role: .destructive) {
                                    showPasskeyDeleteConfirm = true
                                } label: {
                                    Text(language.localized("Geçiş anahtarlarını kaldır", "Remove passkeys"))
                                        .font(.system(size: 13, weight: .medium))
                                        .frame(maxWidth: .infinity)
                                }
                                .disabled(isUpdatingPasskey)
                            }
                        }
                        .padding(14)
                        .background(roundedBg)

                        // Açıklama metni: şifre durumuna göre
                        if hasPassword {
                            Text(language.localized(
                                "Şifrenizi güvenli tutun. En az 6 karakter kullanın.",
                                "Keep your password secure. Use at least 6 characters."
                            ))
                            .font(.system(size: 13))
                            .foregroundColor(Theme.textMuted)
                        } else {
                            // Social login ile gelen kullanıcıya bilgilendirme
                            HStack(spacing: 10) {
                                Image(systemName: "info.circle.fill")
                                    .foregroundColor(Theme.neonPrimary)
                                    .font(.system(size: 18))
                                Text(language.localized(
                                    "Google veya Apple ile giriş yaptığınız için henüz bir şifreniz yok. Aşağıdan yeni bir şifre belirleyebilirsiniz.",
                                    "You signed in with Google or Apple, so you don't have a password yet. You can set a new password below."
                                ))
                                .font(.system(size: 13))
                                .foregroundColor(Theme.textMuted)
                            }
                            .padding(14)
                            .background(
                                RoundedRectangle(cornerRadius: 10)
                                    .fill(Theme.neonPrimary.opacity(0.08))
                            )
                        }

                        // Mevcut şifre alanı — sadece şifresi varsa göster
                        if hasPassword {
                            inputRow(
                                label: language.localized("Mevcut şifre", "Current password"),
                                text: $currentPassword,
                                field: .current
                            )
                        }

                        inputRow(
                            label: language.localized("Yeni şifre (min. 6)", "New password (min. 6)"),
                            text: $newPassword,
                            field: .new
                        )
                        inputRow(
                            label: language.localized("Yeni şifre tekrar", "Confirm password"),
                            text: $confirmPassword,
                            field: .confirm
                        )

                        if let err = errorMessage { Text(err).font(.caption).foregroundColor(.red) }
                        if let ok = successMessage { Text(ok).font(.caption).foregroundColor(.green) }

                        Button {
                            Task { await savePassword() }
                        } label: {
                            HStack {
                                if isSaving {
                                    ProgressView().tint(.black).scaleEffect(0.9)
                                } else {
                                    Text(hasPassword
                                         ? language.localized("Şifreyi güncelle", "Update password")
                                         : language.localized("Şifre belirle", "Set password"))
                                        .font(.system(size: 14, weight: .semibold))
                                }
                            }
                            .foregroundColor(.black)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 10)
                            .background(Theme.neon)
                            .cornerRadius(10)
                        }
                        .disabled(!canSave)
                        .opacity(canSave ? 1 : 0.5)

                        // Hesap silme (Apple Guideline 5.1.1(v))
                        Divider().padding(.vertical, 16)

                        VStack(alignment: .leading, spacing: 8) {
                            Text(language.localized("Hesabı sil", "Delete account"))
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundColor(.red)

                            Text(language.localized(
                                "Hesabınız ve tüm verileriniz kalıcı olarak silinir. Bu işlem geri alınamaz.",
                                "Your account and all data will be permanently deleted. This action cannot be undone."
                            ))
                            .font(.system(size: 13))
                            .foregroundColor(Theme.textMuted)

                            Button {
                                showDeleteConfirm = true
                            } label: {
                                HStack {
                                    if isDeleting {
                                        ProgressView().tint(Theme.textPrimary).scaleEffect(0.9)
                                    } else {
                                        Image(systemName: "trash")
                                        Text(language.localized("Hesabımı kalıcı olarak sil", "Permanently delete my account"))
                                            .font(.system(size: 14, weight: .semibold))
                                    }
                                }
                                .foregroundColor(.white) // Kırmızı buton üzerinde beyaz kalmalı
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 12)
                                .background(Color.red.opacity(0.8))
                                .cornerRadius(10)
                            }
                            .disabled(isDeleting)
                        }
                    }
                    .padding(20)
                }
                .scrollDismissesKeyboard(.interactively)
            }
        }
        .navigationTitle(language.localized("Güvenlik ayarları", "Security settings"))
        .navigationBarTitleDisplayMode(.inline)
        .alert(
            language.localized("Hesabı sil", "Delete account"),
            isPresented: $showDeleteConfirm
        ) {
            Button(language.localized("İptal", "Cancel"), role: .cancel) {}
            Button(language.localized("Sil", "Delete"), role: .destructive) {
                Task { await performDeleteAccount() }
            }
        } message: {
            Text(language.localized(
                "Bu işlem geri alınamaz. Hesabınız ve tüm verileriniz kalıcı olarak silinecektir.",
                "This action cannot be undone. Your account and all data will be permanently deleted."
            ))
        }
        .confirmationDialog(
            language.localized("Geçiş anahtarlarını kaldır", "Remove passkeys"),
            isPresented: $showPasskeyDeleteConfirm,
            titleVisibility: .visible
        ) {
            Button(language.localized("Tümünü kaldır", "Remove all"), role: .destructive) {
                Task { await deletePasskeys() }
            }
            Button(language.localized("İptal", "Cancel"), role: .cancel) {}
        } message: {
            Text(language.localized(
                "Bu hesaba ait tüm geçiş anahtarları kaldırılacak. Apple Kimliği, Google veya e-posta ile giriş yapmaya devam edebilirsiniz.",
                "All passkeys for this account will be removed. You can continue signing in with Apple, Google, or email."
            ))
        }
        .task {
            await checkHasPassword()
            await checkPasskeyStatus()
        }
    }

    // MARK: - API

    /// /me'den hasPassword bilgisini çek
    private func checkHasPassword() async {
        isCheckingPassword = true
        defer { isCheckingPassword = false }
        do {
            let user = try await ProfileService.shared.fetchMe()
            hasPassword = user.hasPassword ?? true
        } catch {
            hasPassword = true // hata durumunda güvenli varsayılan
        }
    }

    private func checkPasskeyStatus() async {
        do {
            let status = try await APIClient.shared.passkeyStatus()
            passkeyEnabled = status.enabled
            passkeyCount = status.count
        } catch {
            passkeyEnabled = false
            passkeyCount = 0
        }
    }

    private func createPasskey() async {
        isUpdatingPasskey = true
        errorMessage = nil
        successMessage = nil
        defer { isUpdatingPasskey = false }
        do {
            try await AuthService.shared.registerPasskey()
            await checkPasskeyStatus()
            successMessage = language.localized("Geçiş anahtarınız hazır.", "Your passkey is ready.")
        } catch let error as ASAuthorizationError where error.code == .canceled {
            // Sistem sayfası kullanıcı tarafından kapatıldı.
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func deletePasskeys() async {
        isUpdatingPasskey = true
        errorMessage = nil
        successMessage = nil
        defer { isUpdatingPasskey = false }
        do {
            try await APIClient.shared.deletePasskeys()
            passkeyEnabled = false
            passkeyCount = 0
            successMessage = language.localized("Geçiş anahtarları kaldırıldı.", "Passkeys removed.")
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func performDeleteAccount() async {
        isDeleting = true
        errorMessage = nil
        defer { isDeleting = false }
        do {
            try await ProfileService.shared.deleteAccount()
            await AuthService.shared.logout()
            onAccountDeleted?()
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func inputRow(label: String, text: Binding<String>, field: PasswordField) -> some View {
        VStack(alignment: .leading, spacing: 6) {
            Text(label)
                .font(.system(size: 13, weight: .medium))
                .foregroundColor(Theme.textMuted)

            SecureField("", text: text)
                .textFieldStyle(.plain)
                .foregroundColor(Theme.textPrimary)
                .font(.system(size: 16))
                .submitLabel(.done)
                .focused($focusedField, equals: field)
                .padding(.horizontal, 16)
                .padding(.vertical, 14)
                .frame(minHeight: 48)
                .background(roundedBg)
                .contentShape(Rectangle())
                .onTapGesture { focusedField = field }
        }
    }

    private func savePassword() async {
        guard newPassword == confirmPassword, newPassword.count >= 6 else {
            errorMessage = language.localized("Yeni şifreler eşleşmiyor veya 6 karakterden kısa", "Passwords don't match or too short")
            return
        }
        isSaving = true
        errorMessage = nil
        successMessage = nil
        defer { isSaving = false }
        do {
            // Social login kullanıcısı ise currentPassword boş gönderilir, backend buna izin verir
            try await ProfileService.shared.updatePassword(
                currentPassword: hasPassword ? currentPassword : "",
                newPassword: newPassword
            )
            successMessage = hasPassword
                ? language.localized("Şifre güncellendi", "Password updated")
                : language.localized("Şifre başarıyla belirlendi", "Password set successfully")
            currentPassword = ""
            newPassword = ""
            confirmPassword = ""
            // Artık kullanıcının şifresi var
            hasPassword = true
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

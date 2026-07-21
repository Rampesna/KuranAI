import SwiftUI

struct ProfileMenuView: View {
    let language: Language
    @Binding var path: NavigationPath
    let onLogout: () -> Void
    let onDismiss: () -> Void

    @State private var user: UserDTO?
    @State private var isLoading = true

    private var roundedBg: some View {
        RoundedRectangle(cornerRadius: 10).fill(Theme.cardBackground)
    }

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            if isLoading {
                ProgressView().tint(Theme.textPrimary)
            } else if let user = user {
                ScrollView {
                    VStack(spacing: 0) {
                        headerSection(user: user)

                        VStack(spacing: 1) {
                            menuRow(
                                icon: "person.crop.circle",
                                title: language.localized("Profil ayarları", "Profile settings"),
                                subtitle: language.localized("Fotoğraf ve kullanıcı adı", "Photo and username")
                            ) { path.append(Route.profileSettings) }

                            menuRow(
                                icon: "lock.shield",
                                title: language.localized("Güvenlik ayarları", "Security settings"),
                                subtitle: language.localized("Şifre değiştir", "Change password")
                            ) { path.append(Route.securitySettings) }

                            menuRow(
                                icon: "waveform",
                                title: language.localized("AI ayarları", "AI settings"),
                                subtitle: language.localized("Sesli okuma sesi", "Voice reading")
                            ) { path.append(Route.aiSettings) }

                            menuRow(
                                icon: "paintbrush.fill",
                                title: language.localized("Görünüm", "Appearance"),
                                subtitle: language.localized("Tema ve renk modu", "Theme and color mode")
                            ) { path.append(Route.appearanceSettings) }

                            menuRow(
                                icon: "crown.fill",
                                title: language.localized("Abonelik ayarları", "Subscription settings"),
                                subtitle: language.localized("Paket ve kullanım", "Package and usage")
                            ) { path.append(Route.subscriptionSettings) }
                        }
                        .padding(.horizontal, 20)
                        .padding(.top, 8)

                        Button { onLogout() } label: {
                            Text(language.localized("Çıkış Yap", "Log out"))
                                .font(.system(size: 14, weight: .medium))
                                .foregroundColor(Theme.textMuted)
                        }
                        .padding(.top, 32)
                        .padding(.bottom, 40)
                    }
                }
            }
        }
        .navigationBarTitleDisplayMode(.inline)
        .navigationTitle(language.localized("Profil", "Profile"))
        .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
                Button {
                    onDismiss()
                } label: {
                    Image(systemName: "chevron.left")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundColor(Theme.textPrimary)
                }
            }
        }
        .task { await loadUser() }
    }

    private func headerSection(user: UserDTO) -> some View {
        HStack(spacing: 14) {
            avatarView(user: user)
            VStack(alignment: .leading, spacing: 2) {
                Text(user.username ?? user.email ?? "-")
                    .font(.system(size: 17, weight: .semibold))
                    .foregroundColor(Theme.textPrimary)
                Text(user.email ?? "")
                    .font(.system(size: 13))
                    .foregroundColor(Theme.textMuted)
            }
            Spacer()
        }
        .padding(16)
        .background(roundedBg)
        .cornerRadius(12)
        .padding(.horizontal, 20)
        .padding(.top, 12)
        .padding(.bottom, 4)
    }

    private func avatarView(user: UserDTO) -> some View {
        Group {
            if let urlStr = user.avatarUrl, !urlStr.isEmpty, let url = fullAvatarURL(urlStr) {
                AsyncImage(url: url) { phase in
                    switch phase {
                    case .success(let img): img.resizable().scaledToFill()
                    default: placeholderAvatar
                    }
                }
            } else {
                placeholderAvatar
            }
        }
        .frame(width: 52, height: 52)
        .clipShape(Circle())
        .overlay(Circle().stroke(Theme.neonPrimary.opacity(0.6), lineWidth: 1))
    }

    private var placeholderAvatar: some View {
        ZStack {
            Circle().fill(Theme.glassColor)
            Image(systemName: "person.fill").font(.system(size: 22)).foregroundColor(Theme.textMuted)
        }
    }

    private func menuRow(icon: String, title: String, subtitle: String, action: @escaping () -> Void) -> some View {
        Button(action: action) {
            HStack(spacing: 14) {
                Image(systemName: icon)
                    .font(.system(size: 18))
                    .foregroundColor(Theme.neonPrimary)
                    .frame(width: 28, alignment: .center)

                VStack(alignment: .leading, spacing: 2) {
                    Text(title)
                        .font(.system(size: 15, weight: .medium))
                        .foregroundColor(Theme.textPrimary)
                    Text(subtitle)
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textMuted)
                }

                Spacer()
                Image(systemName: "chevron.right")
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundColor(Theme.textMuted.opacity(0.7))
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 12)
            .background(roundedBg)
        }
        .buttonStyle(.plain)
    }

    private func fullAvatarURL(_ path: String) -> URL? {
        if path.hasPrefix("http") { return URL(string: path) }
        let base = APIEnvironment.current.baseURL
        return URL(string: base.hasSuffix("/") ? String(base.dropLast()) + path : base + path)
    }

    private func loadUser() async {
        isLoading = true
        defer { isLoading = false }
        do {
            user = try await ProfileService.shared.fetchMe()
        } catch { }
    }
}

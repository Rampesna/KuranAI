import SwiftUI
import PhotosUI

/// Identifiable wrapper — fullScreenCover(item:) için gerekli
private struct CropItem: Identifiable {
    let id = UUID()
    let image: UIImage
}

struct ProfileSettingsView: View {
    let language: Language
    let user: UserDTO?
    let onUserUpdated: (UserDTO) -> Void
    let onBack: () -> Void

    @State private var currentUser: UserDTO?
    @State private var username: String = ""
    @State private var isSavingUsername = false
    @State private var isUploadingAvatar = false
    @State private var errorMessage: String?
    @State private var successMessage: String?
    @State private var selectedItem: PhotosPickerItem?
    @State private var avatarImage: UIImage?
    @State private var showAvatarPicker = false
    @State private var pendingCropImage: UIImage?    // sheet kapandıktan sonra crop'a geçilecek
    @State private var cropItem: CropItem?           // fullScreenCover(item:) tetikleyicisi
    @FocusState private var isUsernameFocused: Bool
    @FocusState private var isPreloadFocused: Bool

    private var roundedBg: some View {
        RoundedRectangle(cornerRadius: 10).fill(Theme.cardBackground)
    }

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            // Görünmez TextField - klavye ön-yüklemesi için
            TextField("", text: .constant(""))
                .frame(width: 1, height: 1)
                .opacity(0)
                .focused($isPreloadFocused)

            ScrollView {
                VStack(alignment: .leading, spacing: 24) {
                    avatarSection
                    usernameSection
                }
                .padding(20)
            }
            .scrollDismissesKeyboard(.interactively)
        }
        .navigationTitle(language.localized("Profil ayarları", "Profile settings"))
        .navigationBarTitleDisplayMode(.inline)
        .onAppear {
            currentUser = user
            username = user?.username ?? ""
            // Klavye gecikmesini azaltmak için arka planda klavyeyi önceden yükle (görünmez)
            Task { @MainActor in
                try? await Task.sleep(for: .milliseconds(400))
                guard !isUsernameFocused else { return }
                isPreloadFocused = true
                try? await Task.sleep(for: .milliseconds(25))
                isPreloadFocused = false
            }
        }
        .onChange(of: user) { _, u in
            currentUser = u
            username = u?.username ?? ""
        }
        .task {
            if user == nil {
                currentUser = try? await ProfileService.shared.fetchMe()
                username = currentUser?.username ?? ""
            }
        }
        .onChange(of: selectedItem) { _, newValue in
            Task { await loadRawImage(from: newValue) }
        }
        .fullScreenCover(item: $cropItem) { item in
            ImageCropView(
                image: item.image,
                language: language,
                onCropped: { croppedImage in
                    cropItem = nil
                    Task { await uploadCroppedImage(croppedImage) }
                },
                onCancel: {
                    cropItem = nil
                }
            )
        }
    }

    private var avatarSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(language.localized("Profil fotoğrafı", "Profile photo"))
                .font(.system(size: 13, weight: .medium))
                .foregroundColor(Theme.textMuted)

            Button { showAvatarPicker = true } label: {
                HStack(spacing: 16) {
                    ZStack {
                        Group {
                            if let img = avatarImage {
                                Image(uiImage: img).resizable().scaledToFill()
                            } else if let u = currentUser, let urlStr = u.avatarUrl, !urlStr.isEmpty, let url = fullAvatarURL(urlStr) {
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
                        .frame(width: 56, height: 56)
                        .clipShape(Circle())
                        .overlay(Circle().stroke(Theme.neonPrimary, lineWidth: 1.5))
                        .overlay(alignment: .bottomTrailing) {
                            if !isUploadingAvatar {
                                Image(systemName: "camera.fill")
                                    .font(.system(size: 10))
                                    .foregroundColor(.black)
                                    .padding(4)
                                    .background(Circle().fill(Theme.neonPrimary))
                                    .offset(x: 2, y: 2)
                            } else {
                                ProgressView().tint(Theme.textPrimary).scaleEffect(0.7)
                            }
                        }
                    }

                    VStack(alignment: .leading, spacing: 2) {
                        Text(language.localized("Fotoğraf ekle veya değiştir", "Add or change photo"))
                            .font(.system(size: 15, weight: .medium))
                            .foregroundColor(Theme.textPrimary)
                        Text(language.localized("Dokunarak seçebilir veya kaldırabilirsin", "Tap to select or remove"))
                            .font(.system(size: 13))
                            .foregroundColor(Theme.textMuted)
                    }

                    Spacer()
                    Image(systemName: "chevron.right")
                        .font(.system(size: 13))
                        .foregroundColor(Theme.textMuted)
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 14)
                .background(roundedBg)
            }
            .buttonStyle(.plain)
            .disabled(isUploadingAvatar)

            if let err = errorMessage { Text(err).font(.caption).foregroundColor(.red) }
            if let ok = successMessage { Text(ok).font(.caption).foregroundColor(.green) }
        }
        .sheet(isPresented: $showAvatarPicker, onDismiss: {
            // Sheet tamamen kapandıktan sonra crop ekranını aç
            if let img = pendingCropImage {
                pendingCropImage = nil
                cropItem = CropItem(image: img)
            }
        }) {
            AvatarPickerSheet(
                language: language,
                selectedItem: $selectedItem,
                hasAvatar: (currentUser?.avatarUrl ?? "").isEmpty == false,
                onRemove: { Task { await removeAvatar() } },
                onDismiss: { showAvatarPicker = false }
            )
        }
    }

    private var placeholderAvatar: some View {
        ZStack {
            Circle().fill(Theme.glassColor)
            Image(systemName: "person.fill").font(.system(size: 24)).foregroundColor(Theme.textMuted)
        }
    }

    private var usernameSection: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(language.localized("Kullanıcı adı", "Username"))
                .font(.system(size: 13, weight: .medium))
                .foregroundColor(Theme.textMuted)

            HStack(spacing: 12) {
                TextField("", text: $username)
                    .textFieldStyle(.plain)
                    .foregroundColor(Theme.textPrimary)
                    .font(.system(size: 16))
                    .autocapitalization(.none)
                    .autocorrectionDisabled()
                    .submitLabel(.done)
                    .focused($isUsernameFocused)
                    .padding(.horizontal, 16)
                    .padding(.vertical, 14)
                    .frame(minHeight: 48)
                    .background(roundedBg)
                    .contentShape(Rectangle())
                    .onTapGesture { isUsernameFocused = true }

                if isSavingUsername {
                    ProgressView().tint(Theme.textPrimary).scaleEffect(0.7)
                } else {
                    Button(language.localized("Kaydet", "Save")) { Task { await saveUsername() } }
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundStyle(Theme.neon)
                        .disabled(username.isEmpty || username == currentUser?.username)
                }
            }
        }
    }

    private func fullAvatarURL(_ path: String) -> URL? {
        if path.hasPrefix("http") { return URL(string: path) }
        let base = APIEnvironment.current.baseURL
        return URL(string: base.hasSuffix("/") ? String(base.dropLast()) + path : base + path)
    }

    /// Galeriden seçilen fotoğrafı yükle ve sheet kapandıktan sonra crop aç
    private func loadRawImage(from item: PhotosPickerItem?) async {
        guard let item = item,
              let data = try? await item.loadTransferable(type: Data.self),
              let img = UIImage(data: data) else { return }
        // Fotoğrafı geçici tut, sheet kapanınca crop açılacak
        pendingCropImage = img
        showAvatarPicker = false
    }

    /// Kırpılmış fotoğrafı sunucuya yükle
    private func uploadCroppedImage(_ img: UIImage) async {
        avatarImage = img
        isUploadingAvatar = true
        errorMessage = nil
        successMessage = nil
        defer { isUploadingAvatar = false }
        do {
            let updated = try await ProfileService.shared.uploadAvatar(image: img)
            avatarImage = nil
            currentUser = updated
            onUserUpdated(updated)
            successMessage = language.localized("Fotoğraf güncellendi", "Photo updated")
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func removeAvatar() async {
        isUploadingAvatar = true
        errorMessage = nil
        successMessage = nil
        defer { isUploadingAvatar = false }
        do {
            if let updated = try await ProfileService.shared.removeAvatar() {
                currentUser = updated
                onUserUpdated(updated)
                successMessage = language.localized("Fotoğraf kaldırıldı", "Photo removed")
                showAvatarPicker = false
            }
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func saveUsername() async {
        guard !username.trimmingCharacters(in: .whitespaces).isEmpty else { return }
        isSavingUsername = true
        errorMessage = nil
        successMessage = nil
        defer { isSavingUsername = false }
        do {
            let updated = try await ProfileService.shared.updateUsername(username.trimmingCharacters(in: .whitespaces))
            currentUser = updated
            onUserUpdated(updated)
            successMessage = language.localized("Kaydedildi", "Saved")
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

// MARK: - Avatar Picker Sheet

private struct AvatarPickerSheet: View {
    let language: Language
    @Binding var selectedItem: PhotosPickerItem?
    let hasAvatar: Bool
    let onRemove: () async -> Void
    let onDismiss: () -> Void

    @State private var isRemoving = false

    var body: some View {
        NavigationStack {
            Theme.background.ignoresSafeArea()
            VStack(spacing: 20) {
                PhotosPicker(selection: $selectedItem, matching: .images, photoLibrary: .shared()) {
                    Label(language.localized("Fotoğraf Seç", "Select Photo"), systemImage: "photo.on.rectangle.angled")
                        .font(.subheadline.weight(.medium))
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 12)
                        .background(Theme.neon)
                        .foregroundColor(.black)
                        .cornerRadius(10)
                }
                .padding(.horizontal, 24)

                if hasAvatar {
                    Button {
                        Task {
                            isRemoving = true
                            await onRemove()
                            isRemoving = false
                        }
                    } label: {
                        HStack {
                            if isRemoving {
                                ProgressView().tint(Theme.textPrimary).scaleEffect(0.8)
                            } else {
                                Image(systemName: "trash")
                                Text(language.localized("Fotoğrafı Kaldır", "Remove Photo"))
                                    .font(.subheadline.weight(.medium))
                            }
                        }
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 12)
                        .foregroundColor(.red.opacity(0.9))
                        .background(Color.red.opacity(0.15))
                        .cornerRadius(10)
                    }
                    .disabled(isRemoving)
                    .padding(.horizontal, 24)
                }

                Spacer()
            }
            .padding(.top, 32)
            .navigationTitle(language.localized("Profil Fotoğrafı", "Profile Photo"))
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button(language.localized("İptal", "Cancel")) { onDismiss() }
                        .foregroundColor(Theme.neonPrimary)
                }
            }
        }
    }
}

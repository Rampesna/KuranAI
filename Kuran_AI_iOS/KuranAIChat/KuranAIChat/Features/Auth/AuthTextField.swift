import SwiftUI

/// Profil ayarları ve diğer ekranlarla aynı tema: cardBackground, 48pt yükseklik, 10pt radius.
struct AuthTextField: View {
    let icon: String
    let placeholder: String
    @Binding var text: String
    var isSecure: Bool = false

    private var roundedBg: some View {
        RoundedRectangle(cornerRadius: 10).fill(Theme.cardBackground)
    }

    var body: some View {
        HStack(spacing: 14) {
            Image(systemName: icon)
                .font(.system(size: 16))
                .foregroundColor(Theme.neonPrimary)
                .frame(width: 24, alignment: .center)

            if isSecure {
                SecureField(placeholder, text: $text)
                    .textFieldStyle(.plain)
                    .font(.system(size: 16))
                    .foregroundColor(Theme.textPrimary)
                    .frame(maxWidth: .infinity)
                    .submitLabel(.next)
            } else {
                TextField(placeholder, text: $text)
                    .textFieldStyle(.plain)
                    .font(.system(size: 16))
                    .foregroundColor(Theme.textPrimary)
                    .frame(maxWidth: .infinity)
                    .submitLabel(.next)
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 14)
        .frame(minHeight: 48)
        .background(roundedBg)
        .overlay(
            RoundedRectangle(cornerRadius: 10)
                .stroke(Theme.border, lineWidth: 1)
        )
    }
}

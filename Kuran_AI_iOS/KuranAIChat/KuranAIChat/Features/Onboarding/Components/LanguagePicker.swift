import SwiftUI

struct LanguagePicker: View {
    @Binding var selectedLanguage: Language
    @State private var showLanguageSheet = false

    var body: some View {
        HStack {
            Spacer()
            Button {
                showLanguageSheet = true
            } label: {
                HStack(spacing: 6) {
                    Image(systemName: "globe")
                    Text(selectedLanguage.displayName)
                }
                .font(.system(size: 13, weight: .semibold))
                .padding(.horizontal, 14)
                .padding(.vertical, 10)
                .background(Theme.glassColor)
                .cornerRadius(16)
                .foregroundColor(Theme.textPrimary)
            }
            .sheet(isPresented: $showLanguageSheet) {
                LanguageSheet(language: selectedLanguage, selectedLanguage: $selectedLanguage, isPresented: $showLanguageSheet)
            }
        }
    }
}

private struct LanguageSheet: View {
    let language: Language
    @Binding var selectedLanguage: Language
    @Binding var isPresented: Bool

    var body: some View {
        NavigationStack {
            List(Language.allCases, id: \.id) { lang in
                Button {
                    selectedLanguage = lang
                    isPresented = false
                } label: {
                    HStack {
                        Text(lang.displayName)
                            .foregroundColor(.primary)
                        if lang == selectedLanguage {
                            Spacer()
                            Image(systemName: "checkmark.circle.fill")
                                .foregroundColor(Theme.neonPrimary)
                        }
                    }
                }
            }
            .navigationTitle(language.localized("Dil", "Language"))
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button(language.localized("Kapat", "Close")) {
                        isPresented = false
                    }
                }
            }
        }
    }
}

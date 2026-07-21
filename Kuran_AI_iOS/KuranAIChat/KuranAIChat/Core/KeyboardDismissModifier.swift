import SwiftUI

extension View {
    /// Klavye açıkken ekrana tıklanınca klavye kapanır.
    func dismissKeyboardOnTap() -> some View {
        modifier(KeyboardDismissModifier())
    }
}

private struct KeyboardDismissModifier: ViewModifier {
    @State private var keyboardVisible = false

    func body(content: Content) -> some View {
        content
            .overlay {
                if keyboardVisible {
                    Color.black.opacity(0.001)
                        .contentShape(Rectangle())
                        .onTapGesture { hideKeyboard() }
                        .ignoresSafeArea()
                }
            }
            .onReceive(NotificationCenter.default.publisher(for: UIResponder.keyboardWillShowNotification)) { _ in
                keyboardVisible = true
            }
            .onReceive(NotificationCenter.default.publisher(for: UIResponder.keyboardWillHideNotification)) { _ in
                keyboardVisible = false
            }
    }

    private func hideKeyboard() {
        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
    }
}

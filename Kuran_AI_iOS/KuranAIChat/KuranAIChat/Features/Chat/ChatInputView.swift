import SwiftUI
import UIKit

/// UITextView tabanlı input - gerçek içerik boyutuna göre yükseklik.
struct ChatInputView: UIViewRepresentable {
    @Binding var text: String
    let placeholder: String
    let onHeightChange: (CGFloat) -> Void

    func makeCoordinator() -> Coordinator {
        Coordinator(self)
    }

    func makeUIView(context: Context) -> ChatInputWrapper {
        let wrapper = ChatInputWrapper()
        wrapper.setup(placeholder: placeholder, coordinator: context.coordinator)
        context.coordinator.wrapper = wrapper
        return wrapper
    }

    func updateUIView(_ wrapper: ChatInputWrapper, context: Context) {
        if wrapper.textView.text != text {
            wrapper.textView.text = text
            wrapper.placeholderLabel.isHidden = !text.isEmpty
        }
        wrapper.placeholderLabel.text = placeholder
        context.coordinator.onHeightChange = onHeightChange
    }

    class Coordinator: NSObject, UITextViewDelegate {
        var parent: ChatInputView
        weak var wrapper: ChatInputWrapper?
        var onHeightChange: ((CGFloat) -> Void)?

        init(_ parent: ChatInputView) {
            self.parent = parent
            super.init()
        }

        func textViewDidChange(_ textView: UITextView) {
            let newText = textView.text ?? ""
            if parent.text != newText {
                parent.text = newText
            }
            wrapper?.placeholderLabel.isHidden = !newText.isEmpty
            guard !newText.isEmpty else { return }
            wrapper?.reportHeightAfterLayout(coordinator: self)
        }
    }
}

final class ChatInputWrapper: UIView {
    let textView = UITextView()
    let placeholderLabel = UILabel()

    private static let minHeight: CGFloat = 44
    private static let maxHeight: CGFloat = 120
    private static let verticalPadding: CGFloat = 24

    private var lastReportedHeight: CGFloat = 44

    func setup(placeholder: String, coordinator: ChatInputView.Coordinator) {
        backgroundColor = .clear

        textView.translatesAutoresizingMaskIntoConstraints = false
        textView.delegate = coordinator
        textView.font = .systemFont(ofSize: 16)
        textView.textColor = UIColor { trait in
            trait.userInterfaceStyle == .dark
                ? .white
                : UIColor(red: 0.12, green: 0.10, blue: 0.06, alpha: 1)
        }
        textView.backgroundColor = .clear
        textView.textContainerInset = UIEdgeInsets(top: 12, left: 20, bottom: 12, right: 20)
        textView.textContainer.lineFragmentPadding = 0
        textView.isScrollEnabled = true
        textView.showsVerticalScrollIndicator = false
        textView.returnKeyType = .default
        textView.keyboardType = .default
        textView.autocapitalizationType = .sentences
        textView.autocorrectionType = .default
        addSubview(textView)

        placeholderLabel.translatesAutoresizingMaskIntoConstraints = false
        placeholderLabel.text = placeholder
        placeholderLabel.font = .systemFont(ofSize: 16)
        placeholderLabel.textColor = UIColor { trait in
            trait.userInterfaceStyle == .dark
                ? UIColor.white.withAlphaComponent(0.5)
                : UIColor(red: 0.42, green: 0.37, blue: 0.28, alpha: 0.7)
        }
        placeholderLabel.isUserInteractionEnabled = false
        addSubview(placeholderLabel)

        NSLayoutConstraint.activate([
            textView.topAnchor.constraint(equalTo: topAnchor),
            textView.leadingAnchor.constraint(equalTo: leadingAnchor),
            textView.trailingAnchor.constraint(equalTo: trailingAnchor),
            textView.bottomAnchor.constraint(equalTo: bottomAnchor),
            placeholderLabel.topAnchor.constraint(equalTo: topAnchor, constant: 12),
            placeholderLabel.leadingAnchor.constraint(equalTo: leadingAnchor, constant: 20),
            placeholderLabel.trailingAnchor.constraint(equalTo: trailingAnchor, constant: -20)
        ])
    }

    override func layoutSubviews() {
        super.layoutSubviews()
        if !(textView.text ?? "").isEmpty, let coordinator = textView.delegate as? ChatInputView.Coordinator {
            reportHeight(coordinator: coordinator)
        }
    }

    private var heightReportWorkItem: DispatchWorkItem?

    func reportHeightAfterLayout(coordinator: ChatInputView.Coordinator) {
        heightReportWorkItem?.cancel()
        let item = DispatchWorkItem { [weak self, weak coordinator] in
            guard let self = self, let coordinator = coordinator else { return }
            self.reportHeight(coordinator: coordinator)
        }
        heightReportWorkItem = item
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.05, execute: item)
    }

    private func notifyHeight(_ height: CGFloat, coordinator: ChatInputView.Coordinator) {
        guard abs(height - lastReportedHeight) > 0.5 else { return }
        lastReportedHeight = height
        coordinator.onHeightChange?(height)
    }

    func reportHeight(coordinator: ChatInputView.Coordinator) {
        let text = textView.text ?? ""
        if text.isEmpty {
            lastReportedHeight = Self.minHeight
            coordinator.onHeightChange?(Self.minHeight)
            return
        }
        let newlineCount = text.filter { $0 == "\n" }.count
        let lineHeight: CGFloat = 22
        let lineCount: Int
        if newlineCount > 0 {
            lineCount = newlineCount + 1
        } else {
            let layoutManager = textView.layoutManager
            let textContainer = textView.textContainer
            let originalSize = textContainer.size
            textContainer.size = CGSize(width: originalSize.width, height: 10000)
            layoutManager.ensureLayout(for: textContainer)
            var count = 0
            let range = layoutManager.glyphRange(for: textContainer)
            if range.length > 0 {
                layoutManager.enumerateLineFragments(forGlyphRange: range) { _, _, _, _, _ in count += 1 }
            }
            textContainer.size = originalSize
            layoutManager.ensureLayout(for: textContainer)
            lineCount = max(1, count)
        }
        let totalHeight = min(Self.maxHeight, Self.minHeight + CGFloat(lineCount - 1) * lineHeight)
        textView.isScrollEnabled = lineCount > 5
        notifyHeight(totalHeight, coordinator: coordinator)
    }
}

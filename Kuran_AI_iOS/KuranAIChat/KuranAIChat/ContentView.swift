import SwiftUI
import StoreKit
import AuthenticationServices

// MARK: - Root

struct AppRootView: View {
    @State private var isAuthenticated = false
    @State private var showLogin = false
    @State private var showProfile = false
    @State private var path = NavigationPath()
    @State private var profilePath = NavigationPath()
    @AppStorage("kuranai_language") private var languageCode = Language.turkish.rawValue
    @State private var selectedSource: SourceFilter = .kuran
    @State private var pendingInitialQuery: String?
    @State private var skipNavigateFromRestore = false
    @ObservedObject private var themeManager = ThemeManager.shared

    private var language: Language {
        Language(rawValue: languageCode) ?? .turkish
    }

    private var languageBinding: Binding<Language> {
        Binding(
            get: { language },
            set: { languageCode = $0.rawValue }
        )
    }

    var body: some View {
        ZStack {
            // Ana NavigationStack (chat vb.)
            NavigationStack(path: $path) {
                OnboardingView(
                    language: languageBinding,
                    selectedSource: $selectedSource,
                    isAuthenticated: isAuthenticated,
                    onProfileTap: { openProfile() },
                    onPrimaryAction: { handleAction(initialQuery: $0) }
                )
                .background(Theme.background.ignoresSafeArea())
                .navigationDestination(for: Route.self) { route in
                    switch route {
                    case .chat(let initialQuery, let sessionId):
                        ChatView(
                            language: language,
                            selectedSource: $selectedSource,
                            isAuthenticated: $isAuthenticated,
                            showLogin: $showLogin,
                            initialQuery: initialQuery,
                            sessionId: sessionId,
                            onBack: { path.removeLast() },
                            onSelectSession: { id in
                                path.removeLast()
                                path.append(Route.chat(initialQuery: nil, sessionId: id))
                            },
                            onNewChat: {
                                path.removeLast()
                                path.append(Route.chat(initialQuery: nil, sessionId: nil))
                            }
                        )
                    default:
                        EmptyView()
                    }
                }
            }

            // Profil overlay – soldan sağa kayarak gelir
            if showProfile {
                Theme.dimOverlay
                    .ignoresSafeArea()
                    .onTapGesture { closeProfile() }
                    .transition(.opacity)

                NavigationStack(path: $profilePath) {
                    ProfileMenuView(
                        language: language,
                        path: $profilePath,
                        onLogout: {
                            closeProfile()
                            Task { await AuthService.shared.logout() }
                            isAuthenticated = false
                            path = NavigationPath()
                        },
                        onDismiss: { closeProfile() }
                    )
                    .navigationDestination(for: Route.self) { route in
                        profileSubDestination(route: route, language: language, path: $profilePath, onAccountDeleted: {
                            closeProfile()
                            isAuthenticated = false
                            path = NavigationPath()
                        })
                    }
                }
                .transition(.move(edge: .leading))
                .gesture(
                    DragGesture(minimumDistance: 50, coordinateSpace: .local)
                        .onEnded { value in
                            let shouldClose = language == .arabic
                                ? value.translation.width > 50
                                : value.translation.width < -50
                            if shouldClose {
                                closeProfile()
                            }
                        }
                )
            }
        }
        .animation(.easeInOut(duration: 0.3), value: showProfile)
        .sheet(isPresented: $showLogin) {
            LoginGateView(language: language, isAuthenticated: $isAuthenticated)
        }
        .onChange(of: isAuthenticated) { _, newValue in
            if newValue {
                showLogin = false
                guard !skipNavigateFromRestore else {
                    skipNavigateFromRestore = false
                    return
                }
                if let query = pendingInitialQuery, path.isEmpty {
                    pendingInitialQuery = nil
                    path.append(Route.chat(initialQuery: query, sessionId: nil))
                } else {
                    pendingInitialQuery = nil
                }
            }
        }
        .onAppear {
            Task {
                await AuthService.shared.restoreSession()
                if AuthService.shared.isLoggedIn {
                    skipNavigateFromRestore = true
                    isAuthenticated = true
                }
            }
        }
        .onReceive(NotificationCenter.default.publisher(for: AuthService.sessionExpiredNotification)) { _ in
            // Refresh token da geçersiz oldu — oturumu kapat, login ekranını göster
            isAuthenticated = false
            showProfile = false
            path = NavigationPath()
            profilePath = NavigationPath()
            showLogin = true
        }
        .preferredColorScheme(themeManager.colorScheme)
        .environment(\.layoutDirection, language.layoutDirection)
    }

    func handleAction(initialQuery: String? = nil) {
        path.append(Route.chat(initialQuery: initialQuery, sessionId: nil))
    }

    private func openProfile() {
        if isAuthenticated {
            profilePath = NavigationPath()
            showProfile = true
        } else {
            showLogin = true
        }
    }

    private func closeProfile() {
        showProfile = false
        profilePath = NavigationPath()
    }
}

@ViewBuilder
private func profileSubDestination(route: Route, language: Language, path: Binding<NavigationPath>, onAccountDeleted: @escaping () -> Void) -> some View {
    switch route {
    case .profileSettings:
        ProfileSettingsView(
            language: language,
            user: nil,
            onUserUpdated: { _ in },
            onBack: { path.wrappedValue.removeLast() }
        )
    case .securitySettings:
        SecuritySettingsView(
            language: language,
            onBack: { path.wrappedValue.removeLast() },
            onAccountDeleted: onAccountDeleted
        )
    case .aiSettings:
        AISettingsView(language: language, onBack: { path.wrappedValue.removeLast() })
    case .appearanceSettings:
        AppearanceSettingsView(language: language, onBack: { path.wrappedValue.removeLast() })
    case .subscriptionSettings:
        SubscriptionSettingsView(language: language, onBack: { path.wrappedValue.removeLast() })
    default:
        EmptyView()
    }
}

// MARK: - Chat

private struct DisplayMessage: Identifiable {
    let id: UUID
    let role: String
    var content: String
    var isStreaming: Bool
    init(role: String, content: String, isStreaming: Bool, id: UUID = UUID()) {
        self.id = id
        self.role = role
        self.content = content
        self.isStreaming = isStreaming
    }
}

private enum PendingAIConsentAction {
    case sendMessage
    case startVoiceRecording
}

struct ChatView: View {
    @Environment(\.requestReview) private var requestReview
    let language: Language
    @Binding var selectedSource: SourceFilter
    @Binding var isAuthenticated: Bool
    @Binding var showLogin: Bool
    let initialQuery: String?
    let sessionId: UInt?
    let onBack: () -> Void
    let onSelectSession: (UInt) -> Void
    let onNewChat: () -> Void
    @State private var input = ""
    @State private var messages: [DisplayMessage] = []
    @State private var currentSessionId: UInt?
    @State private var isRecording = false
    @State private var isTranscribing = false
    @State private var isSending = false
    @State private var transcribeError: String?
    @State private var sendError: String?
    @FocusState private var isInputFocused: Bool
    @State private var showSidebar = false
    @State private var sessions: [ChatSessionDTO] = []
    @State private var isLoadingMessages = false
    @State private var scrollTrigger: Int = 0
    @State private var userScrolledAway = false
    @ObservedObject private var speechService = SpeechService.shared
    @State private var showUpgradeSheet = false
    @State private var upgradeSheetLimitInfo: APIError.LimitExceededInfo?
    @State private var showAIConsentAlert = false
    @State private var pendingAIConsentAction: PendingAIConsentAction?

    /// Kullanıcı AI veri paylaşımını kabul etti mi (Apple Guideline 5.1.1(i))
    private static let aiConsentKey = "kuranai_ai_data_consent"
    private var hasAIConsent: Bool {
        UserDefaults.standard.bool(forKey: Self.aiConsentKey)
    }

    var body: some View {
        ZStack(alignment: .leading) {
            Theme.background.ignoresSafeArea()

            VStack(spacing: 0) {
                chatHeader
                sourceBar

                ScrollViewReader { proxy in
                    ScrollView {
                        LazyVStack(spacing: 14) {
                            if messages.isEmpty && !isLoadingMessages {
                                GlassBubble(text: language.localized("Selam! Sesli soru sorarak Kur'an, tefsir ve hadis kaynaklarından anında cevap alabilirsin.", "Hi! Use your voice to pull answers from Quran, tafsir, and hadith."), isUser: false, speechService: speechService)
                            }
                            ForEach(messages) { msg in
                                GlassBubble(
                                    text: msg.content,
                                    isUser: msg.role == "user",
                                    showSpeakButton: msg.role == "assistant" && !msg.isStreaming && !msg.content.isEmpty,
                                    speakButtonLabel: language.localized("Sesli oku", "Read aloud"),
                                    stopButtonLabel: language.localized("Durdur", "Stop"),
                                    onSpeakTap: msg.role == "assistant" ? { speechService.toggleSpeak(msg.content, language: language.speechLanguage) } : nil,
                                    speechService: speechService
                                )
                                .id(msg.id)
                            }
                        }
                        .padding(.leading, 20)
                        .padding(.trailing, 8)
                    }
                    .scrollDismissesKeyboard(.interactively)
                    .contentShape(Rectangle())
                    .onTapGesture {
                        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
                    }
                    .simultaneousGesture(
                        DragGesture().onChanged { _ in userScrolledAway = true }
                    )
                    .onChange(of: scrollTrigger) { _, _ in
                        guard !userScrolledAway, let lastId = messages.last?.id else { return }
                        withAnimation(.easeOut(duration: 0.25)) { proxy.scrollTo(lastId, anchor: .bottom) }
                    }
                }

                VStack(spacing: 4) {
                    if let err = transcribeError ?? sendError {
                        Text(err)
                            .font(.caption)
                            .foregroundColor(.red)
                            .padding(.horizontal, 4)
                    }
                    InputStack(
                        language: language,
                        input: $input,
                        isRecording: $isRecording,
                        isTranscribing: isTranscribing || isSending,
                        onSend: handleSend,
                        onMicTap: handleMicTap
                    )
                }
            }
            .padding(.vertical, 12)
        }
        .overlay(alignment: .leading) {
            if showSidebar {
                Theme.dimOverlay
                    .ignoresSafeArea()
                    .onTapGesture { showSidebar = false }
            }
        }
        .sheet(isPresented: $showUpgradeSheet) {
            SubscriptionUpgradeView(
                language: language,
                limitInfo: upgradeSheetLimitInfo,
                onDismiss: { showUpgradeSheet = false },
                onPurchaseSuccess: { showUpgradeSheet = false }
            )
        }
        .alert(
            language.localized("AI Veri Paylaşımı", "AI Data Sharing"),
            isPresented: $showAIConsentAlert
        ) {
            Button(language.localized("İzin Ver ve Devam Et", "Allow and Continue")) {
                UserDefaults.standard.set(true, forKey: Self.aiConsentKey)
                let action = pendingAIConsentAction
                pendingAIConsentAction = nil
                switch action {
                case .sendMessage:
                    handleSend()
                case .startVoiceRecording:
                    Task { await handleMicTap() }
                case .none:
                    break
                }
            }
            Button(language.localized("Şimdi Değil", "Not Now"), role: .cancel) {
                pendingAIConsentAction = nil
            }
            Button(language.localized("Gizlilik Politikasını Aç", "Open Privacy Policy")) {
                pendingAIConsentAction = nil
                if let url = URL(string: "https://kuranai.talhacan.com/privacy-policy") {
                    UIApplication.shared.open(url)
                }
            }
        } message: {
            Text(language.localized(
                "Yazdığınız mesaj metni ve sesli soru kullandığınızda ses kaydınız, yanıt üretmek veya sesi metne çevirmek için OpenAI, Inc. ile paylaşılır. Adınız, e-posta adresiniz ve hesap bilgileriniz OpenAI'a gönderilmez. İzin vermezseniz AI sohbeti ve sesli soru özelliği kullanılmaz; diğer özellikleri kullanmaya devam edebilirsiniz.",
                "Your message text and, when you use voice questions, your audio recording are shared with OpenAI, Inc. to generate a response or transcribe the audio. Your name, email address, and account details are not sent to OpenAI. If you decline, AI chat and voice questions remain unavailable while other features continue to work."
            ))
        }
        .onChange(of: speechService.limitExceededInfo) { _, newVal in
            if let info = newVal {
                upgradeSheetLimitInfo = info
                showUpgradeSheet = true
                speechService.limitExceededInfo = nil
            }
        }
        .onDisappear {
            if isRecording {
                AudioRecorderService.shared.cancelRecording()
                isRecording = false
            }
        }
        .overlay(alignment: .bottomTrailing) {
            if userScrolledAway, !messages.isEmpty {
                Button {
                    userScrolledAway = false
                    scrollTrigger += 1
                } label: {
                    Image(systemName: "chevron.down.circle.fill")
                        .font(.system(size: 36))
                        .foregroundStyle(Theme.neon)
                        .symbolRenderingMode(.hierarchical)
                }
                .padding(.trailing, 20)
                .padding(.bottom, 100)
            }
        }
        .overlay(alignment: .leading) {
            ChatSidebar(
                language: language,
                isPresented: $showSidebar,
                sessions: sessions,
                onSelectSession: { onSelectSession($0); showSidebar = false },
                onNewChat: { onNewChat(); showSidebar = false }
            )
            .offset(
                x: showSidebar
                    ? 0
                    : (language == .arabic ? 1 : -1) * (UIScreen.main.bounds.width * 0.8 + 20)
            )
            .animation(.easeInOut(duration: 0.2), value: showSidebar)
        }
        .onAppear {
            currentSessionId = sessionId
            if let q = initialQuery, !q.isEmpty { input = q }
            Task {
                _ = try? await AudioRecorderService.shared.requestPermissionAndPrepare()
                await loadSessions()
                await loadMessagesIfNeeded()
            }
        }
        .onChange(of: sessionId) { _, newVal in
            currentSessionId = newVal
            Task { await loadMessagesIfNeeded() }
        }
        .navigationBarHidden(true)
    }

    private var chatHeader: some View {
        HStack {
            Button {
                onBack()
            } label: {
                Image(systemName: "chevron.left")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(Theme.textPrimary)
                    .frame(width: 40, height: 40)
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(language.localized("Kuran AI", "Quran AI"))
                    .font(.system(size: 20, weight: .bold))
                    .foregroundStyle(Theme.neon)
                Text(language.localized("Ses odaklı asistan", "Voice-first assistant"))
                    .font(.system(size: 13))
                    .foregroundColor(Theme.textMuted)
            }
            .frame(maxWidth: .infinity, alignment: .leading)

            Button {
                showSidebar = true
            } label: {
                Image(systemName: "list.bullet")
                    .font(.system(size: 20, weight: .medium))
                    .foregroundColor(Theme.textPrimary)
                    .frame(width: 40, height: 40)
            }
        }
        .padding(.horizontal, 12)
    }

    private var sourceBar: some View {
        HStack(spacing: 8) {
            ForEach(SourceFilter.allCases) { filter in
                Button {
                    selectedSource = filter
                } label: {
                    Text(filter.title(for: language))
                        .font(.system(size: 13, weight: .semibold))
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 10)
                        .background(
                            RoundedRectangle(cornerRadius: 14)
                                .fill(selectedSource == filter ? Theme.neon : Theme.pillGradient)
                        )
                        .foregroundColor(selectedSource == filter ? .black : Theme.textPrimary)
                }
            }
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 10)
    }

    private struct GlassBubble: View {
        let text: String
        let isUser: Bool
        var showSpeakButton: Bool = false
        var speakButtonLabel: String = "Sesli oku"
        var stopButtonLabel: String = "Durdur"
        var onSpeakTap: (() -> Void)?
        @ObservedObject var speechService: SpeechService

        var body: some View {
            VStack(alignment: isUser ? .trailing : .leading, spacing: 0) {
                if isUser {
                    HStack(alignment: .bottom, spacing: 0) {
                        Spacer(minLength: 0)
                        MarkdownText(text: text)
                            .font(.system(size: 15))
                            .foregroundColor(.black)
                            .padding(16)
                            .background(Theme.neon.opacity(0.92))
                            .cornerRadius(20)
                            .frame(maxWidth: 280, alignment: .trailing)
                    }
                    .frame(maxWidth: .infinity, alignment: .trailing)
                } else {
                    HStack(alignment: .top, spacing: 12) {
                        Circle()
                            .fill(Theme.neonPrimary.opacity(0.5))
                            .frame(width: 6, height: 6)
                            .padding(.top, 10)
                        MarkdownText(text: text)
                            .font(.system(size: 15))
                            .foregroundColor(Theme.textPrimary)
                            .padding(.vertical, 12)
                        Spacer()
                    }
                    .frame(maxWidth: .infinity, alignment: .leading)

                    if showSpeakButton {
                        Button(action: { onSpeakTap?() }) {
                            HStack(spacing: 6) {
                                if speechService.isLoading {
                                    ProgressView()
                                        .scaleEffect(0.8)
                                        .tint(Theme.neonPrimary)
                                } else {
                                    Image(systemName: speechService.isSpeaking ? "stop.circle.fill" : "speaker.wave.2.fill")
                                        .font(.system(size: 16))
                                }
                                Text(speechService.isSpeaking ? stopButtonLabel : speakButtonLabel)
                                    .font(.caption)
                            }
                            .foregroundColor(Theme.neonPrimary.opacity(0.9))
                        }
                        .buttonStyle(.plain)
                        .disabled(speechService.isLoading)
                        .padding(.top, 4)
                        .padding(.leading, 18)
                    }
                }
            }
        }
    }

    /// Markdown string'i paragraflar, boşluklar ve markdown (**, ##, -) ile düzgün render eder.
    private struct MarkdownText: View {
        let text: String

        var body: some View {
            let blocks = paragraphs
            if blocks.isEmpty {
                Text(text)
                    .font(.system(size: 15))
                    .lineSpacing(6)
            } else {
                VStack(alignment: .leading, spacing: 12) {
                    ForEach(Array(blocks.enumerated()), id: \.offset) { _, block in
                        if let attr = try? AttributedString(
                            markdown: portableBlock(block),
                            options: .init(interpretedSyntax: .inlineOnlyPreservingWhitespace)
                        ) {
                            Text(attr)
                                .font(.system(size: 15))
                                .lineSpacing(6)
                                .fixedSize(horizontal: false, vertical: true)
                        } else {
                            Text(block)
                                .font(.system(size: 15))
                                .lineSpacing(6)
                                .fixedSize(horizontal: false, vertical: true)
                        }
                    }
                }
            }
        }

        private var paragraphs: [String] {
            text.components(separatedBy: "\n\n")
                .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
                .filter { !$0.isEmpty }
        }

        private func portableBlock(_ block: String) -> String {
            block.components(separatedBy: .newlines)
                .compactMap { rawLine -> String? in
                    var line = rawLine.trimmingCharacters(in: .whitespaces)
                    while line.hasPrefix("#") {
                        line.removeFirst()
                        line = line.trimmingCharacters(in: .whitespaces)
                    }
                    if line.hasPrefix("- ") || line.hasPrefix("* ") || line.hasPrefix("+ ") {
                        line = "• " + String(line.dropFirst(2))
                    }
                    if line.hasPrefix("|"), line.hasSuffix("|") {
                        let cells = line
                            .dropFirst()
                            .dropLast()
                            .split(separator: "|")
                            .map { $0.trimmingCharacters(in: .whitespaces) }
                        let isSeparator = !cells.isEmpty && cells.allSatisfy {
                            let value = $0.replacingOccurrences(of: ":", with: "")
                            return value.count >= 3 && value.allSatisfy { $0 == "-" }
                        }
                        if isSeparator { return nil }
                        if cells.count > 1 {
                            line = cells[0] + ": " + cells.dropFirst().joined(separator: " — ")
                        }
                    }
                    return line
                }
                .joined(separator: "\n")
        }
    }

    /// Kayıt sırasında stop ikonu etrafında blur’un artıp azaldığı nabız.
    private struct RecordingWaveView: View {
        private static let iconSize: CGFloat = 36
        private static let period: Double = 1.2

        var body: some View {
            TimelineView(.animation(minimumInterval: 0.04)) { context in
                let t = context.date.timeIntervalSinceReferenceDate
                let phase = sin(2 * .pi * t / Self.period)
                let blurRadius: CGFloat = 2 + 4 * CGFloat(0.5 + 0.5 * phase)
                ZStack {
                    Circle()
                        .fill(Color.red.opacity(0.35))
                        .blur(radius: blurRadius)
                    Image(systemName: "stop.circle.fill")
                        .font(.system(size: 30))
                        .foregroundStyle(Color.red.opacity(0.9))
                        .symbolRenderingMode(.hierarchical)
                }
                .frame(width: Self.iconSize, height: Self.iconSize)
            }
        }
    }

    private struct InputStack: View {
        let language: Language
        @Binding var input: String
        @Binding var isRecording: Bool
        let isTranscribing: Bool
        let onSend: () -> Void
        let onMicTap: () async -> Void

        private static let iconSize: CGFloat = 36
        private static let minHeight: CGFloat = 44

        @State private var inputHeight: CGFloat = 44

        private var hasText: Bool {
            !input.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty
        }

        /// Boş input → sadece mic. Metin varken → sadece gönder. Kayıt varken → sadece stop + animasyon.
        private var showOnlyMic: Bool { !hasText && !isRecording }
        private var showOnlySend: Bool { hasText && !isRecording && !isTranscribing }

        @ViewBuilder
        private var micButtonContent: some View {
            if isTranscribing {
                ProgressView()
                    .tint(Theme.neonPrimary)
                    .scaleEffect(0.8)
            } else if isRecording {
                RecordingWaveView()
            
            } else {
                Image(systemName: "mic.circle.fill")
                    .font(.system(size: 30))
                    .foregroundStyle(Theme.neonPrimary)
                    .symbolRenderingMode(.hierarchical)
            }
        }

        var body: some View {
            HStack(alignment: .bottom, spacing: 10) {
                ChatInputView(
                    text: $input,
                    placeholder: language.localized("Sorunu yaz veya mikrofona bas...", "Type your question or tap the mic..."),
                    onHeightChange: { h in inputHeight = h }
                )
                .frame(height: inputHeight)
                .onChange(of: input) { _, newVal in
                    if newVal.isEmpty { inputHeight = Self.minHeight }
                }

                HStack(spacing: 4) {
                    if showOnlyMic || isRecording || isTranscribing {
                        Button {
                            Task { await onMicTap() }
                        } label: {
                            micButtonContent
                                .frame(width: Self.iconSize, height: Self.iconSize)
                                .contentShape(Rectangle())
                        }
                        .buttonStyle(.plain)
                        .disabled(isTranscribing)
                    }
                    if showOnlySend {
                        Button(action: onSend) {
                            Image(systemName: "arrow.up.circle.fill")
                                .font(.system(size: 30))
                                .foregroundStyle(Theme.neonPrimary)
                                .symbolRenderingMode(.hierarchical)
                                .frame(width: Self.iconSize, height: Self.iconSize)
                                .contentShape(Rectangle())
                        }
                        .buttonStyle(.plain)
                    }
                }
                .padding(.trailing, 8)
                .padding(.bottom, 4)
            }
            .background(
                RoundedRectangle(cornerRadius: 24)
                    .fill(Theme.glassColor)
                    .overlay(
                        RoundedRectangle(cornerRadius: 24)
                            .stroke(Theme.border.opacity(0.5), lineWidth: 1)
                    )
            )
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
            .padding(.bottom, 8)
        }
    }

    private func handleSend() {
        guard ensureAuthenticated() else { return }
        let trimmed = input.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else { return }

        // Apple Guideline 5.1.1(i): İlk kullanımda AI veri paylaşımı izni iste
        if !hasAIConsent {
            pendingAIConsentAction = .sendMessage
            showAIConsentAlert = true
            return
        }

        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
        input = ""
        sendError = nil
        isSending = true
        messages.append(DisplayMessage(role: "user", content: trimmed, isStreaming: false))
        messages.append(DisplayMessage(role: "assistant", content: "", isStreaming: true))
        userScrolledAway = false
        scrollTrigger += 1
        let userContent = trimmed
        Task {
            defer { isSending = false }
            do {
                let sid = try await ChatService.shared.sendMessageStream(
                    sessionId: currentSessionId,
                    sourceType: selectedSource,
                    content: userContent,
                    onChunk: { chunk in
                        DispatchQueue.main.async {
                            guard let idx = messages.lastIndex(where: { $0.role == "assistant" && $0.isStreaming }) else { return }
                            var copy = messages
                            var m = copy[idx]
                            m.content += chunk
                            copy[idx] = m
                            messages = copy
                            scrollTrigger += 1
                        }
                    },
                    onCompleted: { content in
                        DispatchQueue.main.async {
                            guard let idx = messages.lastIndex(where: { $0.role == "assistant" && $0.isStreaming }) else { return }
                            var copy = messages
                            var m = copy[idx]
                            m.content = content
                            copy[idx] = m
                            messages = copy
                            scrollTrigger += 1
                        }
                    }
                )
                currentSessionId = sid
                if let idx = messages.lastIndex(where: { $0.role == "assistant" && $0.isStreaming }) {
                    var copy = messages
                    var m = copy[idx]
                    m.isStreaming = false
                    copy[idx] = m
                    messages = copy
                }
                await loadSessions()
                if InAppReviewCoordinator.shared.recordSuccessfulAnswerAndShouldRequest() {
                    requestReview()
                }
            } catch let e as APIError {
                if case .limitExceeded(let info) = e {
                    if messages.count >= 2 {
                        messages.removeLast(2)
                    }
                    upgradeSheetLimitInfo = info
                    showUpgradeSheet = true
                } else {
                    sendError = e.localizedDescription
                }
                if let idx = messages.lastIndex(where: { $0.role == "assistant" && $0.isStreaming }) {
                    var copy = messages
                    var m = copy[idx]
                    m.isStreaming = false
                    copy[idx] = m
                    messages = copy
                }
            } catch {
                sendError = error.localizedDescription
                if let idx = messages.lastIndex(where: { $0.role == "assistant" && $0.isStreaming }) {
                    var copy = messages
                    var m = copy[idx]
                    m.isStreaming = false
                    copy[idx] = m
                    messages = copy
                }
            }
        }
    }

    private func handleMicTap() async {
        guard ensureAuthenticated() else { return }
        if !hasAIConsent {
            pendingAIConsentAction = .startVoiceRecording
            showAIConsentAlert = true
            return
        }
        if isRecording {
            await finishRecording()
        } else {
            transcribeError = nil
            do {
                try await AudioRecorderService.shared.requestPermissionAndPrepare()
                try AudioRecorderService.shared.startRecording { reason in
                    Task { @MainActor in
                        await finishRecording(reason: reason)
                    }
                }
                isRecording = true
            } catch {
                isRecording = false
                transcribeError = error.localizedDescription
            }
        }
    }

    private func finishRecording(reason: AudioRecordingStopReason? = nil) async {
        guard isRecording || AudioRecorderService.shared.isRecording else { return }
        isRecording = false
        transcribeError = nil

        if reason == .noSpeech {
            AudioRecorderService.shared.cancelRecording()
            transcribeError = language.localized(
                "Konuşma algılanamadı. Mikrofona biraz daha yakın konuşup tekrar deneyin.",
                "No speech was detected. Move a little closer to the microphone and try again.",
                "لم يتم اكتشاف كلام. اقترب قليلاً من الميكروفون وحاول مرة أخرى."
            )
            return
        }

        isTranscribing = true
        defer { isTranscribing = false }
        do {
            let data = try AudioRecorderService.shared.stopRecording()
            let text = try await TranscribeService.shared.transcribe(audioData: data)
            let trimmed = text.trimmingCharacters(in: .whitespacesAndNewlines)
            if trimmed.isEmpty {
                transcribeError = language.localized(
                    "Ses anlaşılamadı. Lütfen tekrar deneyin.",
                    "The audio could not be understood. Please try again.",
                    "تعذر فهم الصوت. يرجى المحاولة مرة أخرى."
                )
            } else {
                input = input.isEmpty ? trimmed : (input + " " + trimmed)
            }
        } catch {
            transcribeError = error.localizedDescription
        }
    }

    private func ensureAuthenticated() -> Bool {
        if !isAuthenticated {
            showLogin = true
            return false
        }
        return true
    }

    private func loadSessions() async {
        do {
            let (items, _) = try await ChatService.shared.fetchSessions(pageIndex: 0, pageSize: 20)
            await MainActor.run { sessions = items }
        } catch { }
    }

    private func loadMessagesIfNeeded() async {
        guard let sid = sessionId else {
            await MainActor.run { messages = [] }
            return
        }
        await MainActor.run { isLoadingMessages = true }
        defer { Task { @MainActor in isLoadingMessages = false } }
        do {
            let (items, _) = try await ChatService.shared.fetchMessages(sessionId: sid, pageIndex: 0, pageSize: 100)
            await MainActor.run {
                messages = items.map { DisplayMessage(role: $0.role, content: $0.content, isStreaming: false) }
                userScrolledAway = false
                scrollTrigger += 1
            }
        } catch { }
    }
}

// MARK: - Chat Sidebar

private struct ChatSidebar: View {
    let language: Language
    @Binding var isPresented: Bool
    let sessions: [ChatSessionDTO]
    let onSelectSession: (UInt) -> Void
    let onNewChat: () -> Void

    private static let sidebarBg = Theme.sidebarBackground
    private static let sidebarBorder = Theme.border

    private func formatDate(_ s: String) -> String {
        let fmt = ISO8601DateFormatter()
        fmt.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let d = fmt.date(from: s) ?? ISO8601DateFormatter().date(from: s) {
            let cal = Calendar.current
            if cal.isDateInToday(d) { return language.localized("Bugün", "Today") }
            if cal.isDateInYesterday(d) { return language.localized("Dün", "Yesterday") }
            let df = DateFormatter()
            df.dateStyle = .short
            df.locale = Locale(identifier: language.localeIdentifier)
            return df.string(from: d)
        }
        return s
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            HStack {
                Text(language.localized("Sohbet Geçmişi", "Chat History"))
                    .font(.system(size: 18, weight: .bold))
                    .foregroundColor(Theme.textPrimary)
                Spacer()
                Button {
                    isPresented = false
                } label: {
                    Image(systemName: "xmark.circle.fill")
                        .font(.system(size: 24))
                        .foregroundStyle(Theme.textMuted)
                }
            }
            .padding(20)

            Button {
                onNewChat()
            } label: {
                HStack {
                    Image(systemName: "plus.circle.fill")
                        .foregroundStyle(Theme.neonPrimary)
                    Text(language.localized("Yeni Sohbet", "New Chat"))
                        .font(.system(size: 15, weight: .medium))
                        .foregroundColor(Theme.textPrimary)
                }
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.vertical, 12)
                .padding(.horizontal, 16)
            }
            .buttonStyle(.plain)

            Divider().background(Self.sidebarBorder).padding(.horizontal, 16)

            ScrollView {
                VStack(spacing: 0) {
                    ForEach(sessions) { session in
                        Button {
                            onSelectSession(session.id)
                        } label: {
                            HStack {
                                VStack(alignment: .leading, spacing: 4) {
                                    Text((session.firstMessage ?? "").isEmpty ? language.localized("Sohbet", "Chat") : String((session.firstMessage ?? "").prefix(60)))
                                        .font(.system(size: 15, weight: .medium))
                                        .foregroundColor(Theme.textPrimary)
                                        .lineLimit(2)
                                    Text(formatDate(session.updatedAt))
                                        .font(.system(size: 13))
                                        .foregroundColor(Theme.textMuted)
                                }
                                Spacer()
                                Image(systemName: "chevron.right")
                                    .font(.system(size: 13, weight: .medium))
                                    .foregroundColor(Theme.textMuted)
                            }
                            .padding(.vertical, 14)
                            .padding(.horizontal, 16)
                        }
                        .buttonStyle(.plain)
                        Divider().background(Self.sidebarBorder).padding(.leading, 16)
                    }
                }
            }

            Spacer()
        }
        .frame(width: UIScreen.main.bounds.width * 0.8)
        .frame(maxHeight: .infinity)
        .background(Self.sidebarBg)
        .overlay(
            Rectangle()
                .fill(Self.sidebarBorder)
                .frame(width: 1)
                .frame(maxHeight: .infinity),
            alignment: .trailing
        )
        .shadow(color: .black.opacity(0.6), radius: 24, x: 4, y: 0)
    }
}

// MARK: - Login Gate

struct LoginGateView: View {
    let language: Language
    @Binding var isAuthenticated: Bool
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showEmailLogin = false
    @State private var showRegister = false

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            VStack(spacing: 18) {
                Text(language.localized("Devam etmek için giriş yap", "Authenticate to continue"))
                    .font(.system(size: 22, weight: .bold))
                    .foregroundColor(Theme.textPrimary)
                Text(language.localized("Apple veya Google ile tek dokunuşta giriş yaparak sohbeti aktif et.", "Sign in with Apple or Google to enable chatting."))
                    .font(.system(size: 14))
                    .foregroundColor(Theme.textMuted)
                    .multilineTextAlignment(.center)

                if let msg = errorMessage {
                    Text(msg)
                        .font(.caption)
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)
                }

                VStack(spacing: 12) {
                    LoginButton(language: language, labelKey: ("Apple ile Giriş", "Sign in with Apple"), icon: "applelogo", disabled: isLoading) {
                        await signInWithApple()
                    }
                    LoginButton(language: language, labelKey: ("Google ile Giriş", "Sign in with Google"), icon: "globe", disabled: isLoading) {
                        await signInWithGoogle()
                    }
                    LoginButton(language: language, labelKey: ("Geçiş Anahtarı ile Giriş", "Sign in with a Passkey"), icon: "person.badge.key.fill", disabled: isLoading) {
                        await signInWithPasskey()
                    }
                    LoginButton(language: language, labelKey: ("E-posta ile Devam Et", "Continue with email"), icon: "envelope.fill", filled: true, disabled: isLoading) {
                        showEmailLogin = true
                    }
                }
            }
            .padding(28)
            .overlay {
                if isLoading {
                    Theme.dimOverlay
                        .ignoresSafeArea()
                    ProgressView()
                        .tint(Theme.textPrimary)
                        .scaleEffect(1.2)
                }
            }
        }
        .sheet(isPresented: $showEmailLogin) {
            EmailLoginView(
                language: language,
                onSuccess: {
                    showEmailLogin = false
                    isAuthenticated = true
                },
                onDismiss: { showEmailLogin = false },
                onShowRegister: {
                    showEmailLogin = false
                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.35) {
                        showRegister = true
                    }
                }
            )
        }
        .sheet(isPresented: $showRegister) {
            RegisterView(
                language: language,
                onSuccess: {
                    showRegister = false
                    isAuthenticated = true
                },
                onDismiss: { showRegister = false }
            )
        }
    }

    private func signInWithApple() async {
        errorMessage = nil
        isLoading = true
        defer { isLoading = false }
        do {
            // Kısa bekleme: loading overlay'in render olması için (OAuth sheet öncesi kasma azaltır)
            try? await Task.sleep(for: .milliseconds(120))
            _ = try await AuthService.shared.signInWithApple()
            isAuthenticated = true
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func signInWithGoogle() async {
        errorMessage = nil
        isLoading = true
        defer { isLoading = false }
        do {
            // Kısa bekleme: loading overlay'in render olması için (OAuth sheet öncesi kasma azaltır)
            try? await Task.sleep(for: .milliseconds(120))
            _ = try await AuthService.shared.signInWithGoogle()
            isAuthenticated = true
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private func signInWithPasskey() async {
        errorMessage = nil
        isLoading = true
        defer { isLoading = false }
        do {
            _ = try await AuthService.shared.signInWithPasskey()
            isAuthenticated = true
        } catch let error as ASAuthorizationError where error.code == .canceled {
            // Kullanıcı sistem geçiş anahtarı sayfasını kapattı; hata göstermeyelim.
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    private struct LoginButton: View {
        let language: Language
        let labelKey: (tr: String, en: String)
        let icon: String
        var filled = false
        var disabled = false
        let action: () async -> Void

        var body: some View {
            Button {
                Task { await action() }
            } label: {
                HStack {
                    Image(systemName: icon)
                    Text(language.localized(labelKey.tr, labelKey.en)).fontWeight(.semibold)
                    Spacer()
                }
                .padding()
                .background(
                    RoundedRectangle(cornerRadius: 18)
                        .fill(filled ? Theme.neon : Theme.whiteCardGradient)
                )
                .foregroundColor(.black)
                .opacity(disabled ? 0.7 : 1)
            }
            .disabled(disabled)
        }
    }
}

// MARK: - ContentView

struct ContentView: View {
    var body: some View {
        AppRootView()
    }
}

#Preview {
    ContentView()
}

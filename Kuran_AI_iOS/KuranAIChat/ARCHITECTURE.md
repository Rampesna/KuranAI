# Kuran AI – iOS Architecture

## Folder Structure

```
KuranAIChat/
├── Core/                    # Shared models & utilities
│   ├── Theme.swift          # Design tokens, colors, gradients
│   ├── Language.swift       # Localization
│   ├── Route.swift          # Navigation routes
│   └── SourceFilter.swift   # Meal / Tafsir / Hadith filter
│
├── Networking/              # API layer (axios-style)
│   ├── APIClient.swift      # Base HTTP client (async/await)
│   ├── APIEnvironment.swift # Base URL, env config
│   ├── APIEndpoint.swift    # Endpoint definitions
│   ├── APIError.swift       # Error types
│   ├── Requests/            # Request DTOs
│   │   └── ChatRequest.swift
│   └── Responses/           # Response DTOs
│       ├── AuthResponse.swift
│       └── ChatResponse.swift
│
├── Services/                # Business logic
│   └── ChatService.swift    # Chat API calls
│
├── Features/                # Feature modules
│   └── Onboarding/
│       ├── OnboardingView.swift
│       └── Components/
│           ├── HomeCarousel.swift
│           ├── LanguagePicker.swift
│           └── SourceSelector.swift
│
├── ContentView.swift        # App root, Chat, Auth views
├── KuranAIChatApp.swift
└── Assets.xcassets/
```

## API Usage

```swift
// Set token after login
await APIClient.shared.setAuthToken("jwt-token")

// Send chat message
let response: ChatMessageResponse = try await ChatService.shared.sendMessage(
    source: .tefsir,
    message: "Fatiha tefsiri"
)
```

## Auth (Apple & Google)

- **Sign in with Apple**: Xcode’da **Signing & Capabilities** → **+ Capability** → **Sign in with Apple** ekle. Backend `POST /api/user/auth/login/apple` ile `identityToken` alıp kullanıcıyı login/kayıt eder.
- **Google Sign-In**: Xcode → **File → Add Package Dependencies…** → URL: `https://github.com/google/GoogleSignIn-iOS` → **Add Package** → **GoogleSignIn** ürününü seç. Google Cloud Console’da OAuth 2.0 Client ID (iOS) oluştur, `GoogleService-Info.plist` indirip projeye ekle ve URL scheme’i (reversed client ID) Info’ya ekle. Backend `POST /api/user/auth/login/google` ile `idToken` alıp kullanıcıyı login/kayıt eder.
- **Token**: Giriş sonrası `AuthService.shared.saveAuth(_:)` ile access/refresh token `UserDefaults`’a yazılır; `APIClient.shared.setAuthToken` ile header’a konur. Uygulama açılışında `AuthService.shared.restoreSession()` ile oturum geri yüklenir.

## Environment

- **Development**: Varsayılan `http://localhost:9300` (simulator için).
- **Fiziksel iPhone**: `localhost` telefona işaret eder, backend’e ulaşamaz. Mac IP kullan:
  1. Mac’in IP’sini bul: Terminal’de `ipconfig getifaddr en0` (Wi‑Fi genelde en0).
  2. Xcode → **Product → Scheme → Edit Scheme…** → **Run** → **Arguments** sekmesi → **Environment Variables** → **+**
  3. Name: `KURAN_API_BASE_URL`, Value: `http://MAC_IP:9300` (örn. `http://192.168.1.5:9300`)
  4. OK → Uygulamayı tekrar cihazda çalıştır.
- **Staging / Production**: `APIEnvironment.swift` içinde güncelle.

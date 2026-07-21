import Foundation

enum Route: Hashable {
    case chat(initialQuery: String?, sessionId: UInt?)
    case profile
    case profileSettings
    case securitySettings
    case aiSettings
    case appearanceSettings
    case subscriptionSettings
}

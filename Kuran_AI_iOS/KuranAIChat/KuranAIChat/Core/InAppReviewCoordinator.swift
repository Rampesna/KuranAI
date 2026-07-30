import Foundation

/// Keeps rating prompts tied to meaningful, successful use of the app.
///
/// StoreKit applies its own display quota. These local rules further prevent
/// repetitive prompts and ensure that a new user has time to evaluate KuranAI.
final class InAppReviewCoordinator {
    static let shared = InAppReviewCoordinator()

    private enum Key {
        static let firstValueAt = "review.firstValueAt"
        static let successfulAnswers = "review.successfulAnswers"
        static let lastPromptAt = "review.lastPromptAt"
        static let lastPromptVersion = "review.lastPromptVersion"
    }

    private let defaults: UserDefaults
    private let requiredSuccessfulAnswers = 3
    private let minimumExperience: TimeInterval = 2 * 24 * 60 * 60
    private let promptCooldown: TimeInterval = 120 * 24 * 60 * 60

    private init(defaults: UserDefaults = .standard) {
        self.defaults = defaults
    }

    func recordSuccessfulAnswerAndShouldRequest(now: Date = Date()) -> Bool {
        let firstValueAt: Date
        if let stored = defaults.object(forKey: Key.firstValueAt) as? Date {
            firstValueAt = stored
        } else {
            firstValueAt = now
            defaults.set(now, forKey: Key.firstValueAt)
        }

        let successfulAnswers = defaults.integer(forKey: Key.successfulAnswers) + 1
        defaults.set(successfulAnswers, forKey: Key.successfulAnswers)

        let version = Bundle.main.object(
            forInfoDictionaryKey: "CFBundleShortVersionString"
        ) as? String ?? "unknown"
        let lastPromptVersion = defaults.string(forKey: Key.lastPromptVersion)
        let lastPromptAt = defaults.object(forKey: Key.lastPromptAt) as? Date

        let enoughExperience = successfulAnswers >= requiredSuccessfulAnswers &&
            now.timeIntervalSince(firstValueAt) >= minimumExperience
        let outsideCooldown = lastPromptAt.map {
            now.timeIntervalSince($0) >= promptCooldown
        } ?? true
        let isNewVersion = lastPromptVersion != version

        guard enoughExperience, outsideCooldown, isNewVersion else { return false }

        defaults.set(now, forKey: Key.lastPromptAt)
        defaults.set(version, forKey: Key.lastPromptVersion)
        return true
    }
}

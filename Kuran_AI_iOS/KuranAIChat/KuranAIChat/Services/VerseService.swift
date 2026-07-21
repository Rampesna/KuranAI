import Foundation

final class VerseService {
    static let shared = VerseService()

    private init() {}

    func getVerseOfDay(lang: String) async throws -> DailyVerseDTO {
        let queryItems = [URLQueryItem(name: "lang", value: lang)]
        let response: VerseAPIResponse = try await APIClient.shared.request(
            .verseOfDay(lang: lang),
            queryItems: queryItems
        )
        guard response.isSuccess, let data = response.data else {
            throw APIError.serverError(response.message ?? "Verse not found")
        }
        return data
    }
}

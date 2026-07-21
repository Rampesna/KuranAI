import Foundation

final class PrayerTimesService {
    static let shared = PrayerTimesService()

    private init() {}

    func getPrayerTimes(lat: Double?, lng: Double?) async throws -> PrayerTimesData {
        var queryItems: [URLQueryItem]?
        if let lat = lat, let lng = lng {
            queryItems = [
                URLQueryItem(name: "lat", value: String(lat)),
                URLQueryItem(name: "lng", value: String(lng))
            ]
        }
        let response: PrayerTimesAPIResponse = try await APIClient.shared.request(
            .prayerTimes(lat: lat, lng: lng),
            queryItems: queryItems
        )
        guard response.isSuccess, let data = response.data else {
            throw APIError.serverError(response.message ?? "Invalid response")
        }
        return data
    }
}

import Foundation

struct PrayerTimesAPIResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: PrayerTimesData?
}

struct PrayerTimesData: Decodable {
    let place: PrayerPlace?
    let times: PrayerTimesMap?
    let date: String?
}

struct PrayerPlace: Decodable {
    let country: String?
    let name: String?
    let stateName: String?
}

struct PrayerTimesMap: Decodable {
    let imsak: String?
    let gunes: String?
    let ogle: String?
    let ikindi: String?
    let aksam: String?
    let yatsi: String?
}

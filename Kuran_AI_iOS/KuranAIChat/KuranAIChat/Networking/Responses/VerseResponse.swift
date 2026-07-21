import Foundation

struct VerseAPIResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: DailyVerseDTO?
}

struct DailyVerseDTO: Decodable {
    let id: UInt
    let surahNumber: Int
    let ayahNumber: Int
    let arabicText: String
    let displayDate: String?
    let isActive: Bool?
    let translations: [DailyVerseTranslationDTO]?
}

struct DailyVerseTranslationDTO: Decodable {
    let id: UInt?
    let languageCode: String
    let translationText: String
    let sourceName: String
}

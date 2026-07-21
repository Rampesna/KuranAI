import Foundation

struct ChatMessageResponse: Decodable {
    let id: String
    let content: String
    let sources: [ChatSource]?
    let createdAt: String
}

struct ChatSource: Decodable {
    let type: String
    let reference: String
    let excerpt: String?
}

struct APIResponseWrapper<T: Decodable>: Decodable {
    let data: T
    let success: Bool
    let message: String?
}

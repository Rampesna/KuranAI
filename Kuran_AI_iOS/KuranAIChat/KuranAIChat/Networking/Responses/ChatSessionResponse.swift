import Foundation

struct ChatSessionsAPIResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: ChatSessionsData?
}

struct ChatSessionsData: Decodable {
    let items: [ChatSessionDTO]
    let total: Int
}

struct ChatSessionDTO: Decodable, Identifiable {
    let id: UInt
    let authType: String
    let authId: UInt
    let firstMessage: String?
    let createdAt: String
    let updatedAt: String
}

struct ChatMessagesAPIResponse: Decodable {
    let isSuccess: Bool
    let message: String?
    let statusCode: Int
    let data: ChatMessagesData?
}

struct ChatMessagesData: Decodable {
    let items: [ChatMessageDTO]
    let total: Int
}

struct ChatMessageDTO: Decodable, Identifiable {
    let id: UInt
    let chatSessionId: UInt
    let role: String
    let content: String
    let sourceType: String
    let createdAt: String
}

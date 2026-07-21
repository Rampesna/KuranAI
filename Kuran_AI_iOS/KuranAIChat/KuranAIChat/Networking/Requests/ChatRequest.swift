import Foundation

struct SendMessageRequest: Encodable {
    let sessionId: UInt?
    let sourceType: String
    let content: String

    init(sessionId: UInt? = nil, sourceType: SourceFilter, content: String) {
        self.sessionId = sessionId
        self.sourceType = sourceType.rawValue == "kuran" ? "meal" : sourceType.rawValue
        self.content = content
    }
}

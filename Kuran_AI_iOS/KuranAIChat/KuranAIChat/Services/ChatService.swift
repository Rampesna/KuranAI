import Foundation

final class ChatService {
    static let shared = ChatService()

    private init() {}

    func fetchSessions(pageIndex: Int = 0, pageSize: Int = 20) async throws -> (items: [ChatSessionDTO], total: Int) {
        let queryItems = [
            URLQueryItem(name: "pageIndex", value: String(pageIndex)),
            URLQueryItem(name: "pageSize", value: String(pageSize))
        ]
        let response: ChatSessionsAPIResponse = try await APIClient.shared.request(
            .chatSessions,
            queryItems: queryItems
        )
        guard response.isSuccess, let data = response.data else {
            throw APIError.serverError(response.message ?? "Invalid response")
        }
        return (data.items, data.total)
    }

    func fetchMessages(sessionId: UInt, pageIndex: Int = 0, pageSize: Int = 20) async throws -> (items: [ChatMessageDTO], total: Int) {
        let queryItems = [
            URLQueryItem(name: "pageIndex", value: String(pageIndex)),
            URLQueryItem(name: "pageSize", value: String(pageSize))
        ]
        let response: ChatMessagesAPIResponse = try await APIClient.shared.request(
            .chatMessages(sessionId: sessionId),
            queryItems: queryItems
        )
        guard response.isSuccess, let data = response.data else {
            throw APIError.serverError(response.message ?? "Invalid response")
        }
        return (data.items, data.total)
    }

    func sendMessageStream(
        sessionId: UInt?,
        sourceType: SourceFilter,
        content: String,
        onChunk: @escaping (String) -> Void,
        onCompleted: @escaping (String) -> Void
    ) async throws -> UInt {
        let token = await APIClient.shared.authToken
        guard let token else { throw APIError.serverError("Not authenticated") }

        let path = APIEnvironment.current.baseURL + "/api/v1/chat/message"
        guard let url = URL(string: path) else { throw APIError.invalidURL }

        var body: [String: Any] = ["content": content, "sourceType": sourceType.rawValue == "kuran" ? "meal" : sourceType.rawValue]
        if let sid = sessionId { body["sessionId"] = sid }
        let bodyData = try JSONSerialization.data(withJSONObject: body)

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("text/event-stream", forHTTPHeaderField: "Accept")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.httpBody = bodyData
        request.timeoutInterval = 120
        request.cachePolicy = .reloadIgnoringLocalCacheData

        var config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 120
        config.timeoutIntervalForResource = 120
        config.urlCache = nil
        config.requestCachePolicy = .reloadIgnoringLocalCacheData
        let session = URLSession(configuration: config)

        let (bytes, response) = try await session.bytes(for: request)
        guard let http = response as? HTTPURLResponse else { throw APIError.invalidResponse }
        guard (200...299).contains(http.statusCode) else {
            var errData = Data()
            for try await chunk in bytes { errData.append(chunk) }
            if http.statusCode == 402, let info = parseLimitExceeded(errData) {
                throw APIError.limitExceeded(info)
            }
            throw APIError.httpError(statusCode: http.statusCode, data: errData)
        }

        var lastSessionId: UInt?
        for try await line in bytes.lines {
            guard line.hasPrefix("data: ") else { continue }
            let payload = String(line.dropFirst(6))
            if payload == "[DONE]" { break }
            guard let data = payload.data(using: .utf8),
                  let event = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
                  let type = event["type"] as? String else { continue }

            if type == "started", let sid = event["sessionId"] as? Int {
                lastSessionId = UInt(sid)
            } else if type == "chunk", let content = event["content"] as? String {
                onChunk(content)
            } else if type == "completed", let content = event["content"] as? String {
                onCompleted(content)
            } else if type == "done", let sid = event["sessionId"] as? Int {
                lastSessionId = UInt(sid)
            } else if type == "error", let msg = event["message"] as? String {
                throw APIError.serverError(msg)
            }
        }
        return lastSessionId ?? sessionId ?? 0
    }

    /// TTS API'den ses verisi (MP3) alır.
    func fetchTTSAudio(text: String, voice: String = "ash") async throws -> Data {
        let token = await APIClient.shared.authToken
        guard let token else { throw APIError.serverError("Not authenticated") }

        let path = APIEnvironment.current.baseURL + APIEndpoint.tts.path
        guard let url = URL(string: path) else { throw APIError.invalidURL }

        let body: [String: Any] = ["text": text, "voice": voice]
        let bodyData = try JSONSerialization.data(withJSONObject: body)

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("audio/mpeg, */*", forHTTPHeaderField: "Accept")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.httpBody = bodyData
        request.timeoutInterval = 60

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse else { throw APIError.invalidResponse }
        guard (200...299).contains(http.statusCode) else {
            if http.statusCode == 402, let info = parseLimitExceeded(data) {
                throw APIError.limitExceeded(info)
            }
            throw APIError.httpError(statusCode: http.statusCode, data: data)
        }
        return data
    }
}

private func parseLimitExceeded(_ data: Data) -> APIError.LimitExceededInfo? {
    guard let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
          let d = json["data"] as? [String: Any],
          let code = d["code"] as? String,
          let limitType = d["limitType"] as? String,
          let used = d["used"] as? Int,
          let limit = d["limit"] as? Int else { return nil }
    let message = (json["message"] as? String) ?? "Limit aşıldı"
    var pkgId: UInt?
    if let pid = d["packageId"] as? Int, pid > 0 { pkgId = UInt(pid) }
    return APIError.LimitExceededInfo(code: code, limitType: limitType, used: used, limit: limit, packageId: pkgId, message: message)
}

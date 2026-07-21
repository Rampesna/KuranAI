import Foundation

@MainActor
final class TranscribeService {
    static let shared = TranscribeService()

    private init() {}

    func transcribe(audioData: Data, filename: String = "recording.m4a") async throws -> String {
        let baseURL = APIEnvironment.current.baseURL
        let path = APIEndpoint.transcribe.path
        guard let url = URL(string: baseURL + path) else {
            throw APIError.invalidURL
        }

        let boundary = UUID().uuidString
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        let token = UserDefaults.standard.string(forKey: "kuran_ai_access_token") ?? ""
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
        request.timeoutInterval = 60

        var body = Data()
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"audio\"; filename=\"\(filename)\"\r\n".data(using: .utf8)!)
        body.append("Content-Type: audio/mp4\r\n\r\n".data(using: .utf8)!)
        body.append(audioData)
        body.append("\r\n--\(boundary)--\r\n".data(using: .utf8)!)
        request.httpBody = body

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            if let wrapper = try? JSONDecoder().decode(ServiceResponseMessage.self, from: data) {
                throw APIError.serverError(wrapper.message)
            }
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let wrapper = try decoder.decode(ServiceResponse<TranscribeData>.self, from: data)
        guard wrapper.isSuccess, let payload = wrapper.data else {
            throw APIError.serverError(wrapper.message)
        }
        return payload.text ?? ""
    }
}

private struct ServiceResponseMessage: Decodable {
    let isSuccess: Bool
    let message: String
}

private struct TranscribeData: Decodable {
    let text: String?
}

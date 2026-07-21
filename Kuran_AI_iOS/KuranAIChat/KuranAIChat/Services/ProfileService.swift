import Foundation
import UIKit

@MainActor
final class ProfileService {
    static let shared = ProfileService()

    private init() {}

    func fetchMe() async throws -> UserDTO {
        let wrapper: ServiceResponse<UserDTO> = try await APIClient.shared.request(APIEndpoint.me)
        guard wrapper.isSuccess, let user = wrapper.data else {
            throw APIError.serverError(wrapper.message)
        }
        return user
    }

    func updateUsername(_ username: String) async throws -> UserDTO {
        struct Req: Encodable { let username: String }
        let wrapper: ServiceResponse<UserDTO> = try await APIClient.shared.request(
            APIEndpoint.updateUsername,
            body: Req(username: username)
        )
        guard wrapper.isSuccess, let user = wrapper.data else {
            throw APIError.serverError(wrapper.message)
        }
        return user
    }

    func updatePassword(currentPassword: String, newPassword: String) async throws {
        struct Req: Encodable {
            let currentPassword: String
            let newPassword: String
        }
        let wrapper: ServiceResponse<EmptyData> = try await APIClient.shared.request(
            APIEndpoint.updatePassword,
            body: Req(currentPassword: currentPassword, newPassword: newPassword)
        )
        guard wrapper.isSuccess else {
            throw APIError.serverError(wrapper.message)
        }
    }

    func uploadAvatar(image: UIImage) async throws -> UserDTO {
        guard let jpeg = image.jpegData(compressionQuality: 0.8) else {
            throw APIError.serverError("Failed to convert image")
        }
        let baseURL = APIEnvironment.current.baseURL
        guard let url = URL(string: baseURL + APIEndpoint.uploadAvatar.path) else {
            throw APIError.invalidURL
        }
        let boundary = UUID().uuidString
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        let token = UserDefaults.standard.string(forKey: "kuran_ai_access_token") ?? ""
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
        request.timeoutInterval = APIEnvironment.current.timeout

        var body = Data()
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"avatar\"; filename=\"avatar.jpg\"\r\n".data(using: .utf8)!)
        body.append("Content-Type: image/jpeg\r\n\r\n".data(using: .utf8)!)
        body.append(jpeg)
        body.append("\r\n--\(boundary)--\r\n".data(using: .utf8)!)
        request.httpBody = body

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            if let wrapper = try? JSONDecoder().decode(ServiceResponseMessage.self, from: data) {
                throw APIError.serverError(wrapper.message)
            }
            throw APIError.httpError(statusCode: (response as? HTTPURLResponse)?.statusCode ?? -1, data: data)
        }
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let wrapper = try decoder.decode(ServiceResponse<AvatarUploadData>.self, from: data)
        guard wrapper.isSuccess, let payload = wrapper.data else {
            throw APIError.serverError(wrapper.message)
        }
        return payload.user
    }

    func removeAvatar() async throws -> UserDTO? {
        let wrapper: ServiceResponse<UserDTO> = try await APIClient.shared.request(APIEndpoint.deleteAvatar)
        guard wrapper.isSuccess else {
            throw APIError.serverError(wrapper.message)
        }
        return wrapper.data
    }

    func deleteAccount() async throws {
        let wrapper: ServiceResponse<EmptyData> = try await APIClient.shared.request(APIEndpoint.deleteAccount)
        guard wrapper.isSuccess else {
            throw APIError.serverError(wrapper.message)
        }
    }
}

private struct ServiceResponseMessage: Decodable {
    let isSuccess: Bool
    let message: String
}
private struct AvatarUploadData: Decodable {
    let avatarUrl: String
    let user: UserDTO
}

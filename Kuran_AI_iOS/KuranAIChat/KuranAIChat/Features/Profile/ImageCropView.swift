import SwiftUI

/// Daire şeklinde profil fotoğrafı kırpma view'ı.
/// Kullanıcı pinch-to-zoom ve drag ile fotoğrafı konumlandırır,
/// "Kaydet" ile daire içindeki kısmı kırpılmış UIImage olarak döner.
struct ImageCropView: View {
    let image: UIImage
    let language: Language
    let onCropped: (UIImage) -> Void
    let onCancel: () -> Void

    // Görünür daire çapı
    private let circleDiameter: CGFloat = 300

    @State private var scale: CGFloat = 1.0
    @State private var lastScale: CGFloat = 1.0
    @State private var offset: CGSize = .zero
    @State private var lastOffset: CGSize = .zero

    var body: some View {
        ZStack {
            Color.black.ignoresSafeArea()

            GeometryReader { geo in
                let center = CGPoint(x: geo.size.width / 2, y: geo.size.height / 2)

                ZStack {
                    // Fotoğraf
                    Image(uiImage: image)
                        .resizable()
                        .scaledToFill()
                        .frame(width: circleDiameter, height: circleDiameter)
                        .scaleEffect(scale)
                        .offset(offset)
                        .gesture(dragGesture)
                        .gesture(magnificationGesture)

                    // Karartma maskesi (daire dışı)
                    CropOverlay(circleCenter: center, circleDiameter: circleDiameter)
                        .fill(Color.black.opacity(0.6), style: FillStyle(eoFill: true))
                        .allowsHitTesting(false)

                    // Daire çerçevesi
                    Circle()
                        .stroke(Color.white.opacity(0.7), lineWidth: 1.5)
                        .frame(width: circleDiameter, height: circleDiameter)
                        .position(center)
                        .allowsHitTesting(false)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            }

            // Üst bar
            VStack {
                HStack {
                    Button {
                        onCancel()
                    } label: {
                        Text(language.localized("İptal", "Cancel"))
                            .font(.system(size: 16, weight: .medium))
                            .foregroundColor(.white)
                    }

                    Spacer()

                    Text(language.localized("Fotoğrafı Kırp", "Crop Photo"))
                        .font(.system(size: 17, weight: .semibold))
                        .foregroundColor(.white)

                    Spacer()

                    Button {
                        let cropped = cropImage()
                        onCropped(cropped)
                    } label: {
                        Text(language.localized("Kaydet", "Save"))
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(Theme.neonPrimary)
                    }
                }
                .padding(.horizontal, 20)
                .padding(.top, 12)

                Spacer()

                // Alt ipucu
                Text(language.localized(
                    "Kaydırmak ve yakınlaştırmak için parmağınızı kullanın",
                    "Use your fingers to pan and zoom"
                ))
                .font(.system(size: 13))
                .foregroundColor(.white.opacity(0.6))
                .padding(.bottom, 40)
            }
        }
    }

    // MARK: - Gestures

    private var dragGesture: some Gesture {
        DragGesture()
            .onChanged { value in
                offset = CGSize(
                    width: lastOffset.width + value.translation.width,
                    height: lastOffset.height + value.translation.height
                )
            }
            .onEnded { _ in
                lastOffset = offset
            }
    }

    private var magnificationGesture: some Gesture {
        MagnificationGesture()
            .onChanged { value in
                let newScale = lastScale * value
                scale = min(max(newScale, 1.0), 5.0)
            }
            .onEnded { _ in
                lastScale = scale
            }
    }

    // MARK: - Crop Logic

    private func cropImage() -> UIImage {
        let imageSize = image.size
        let displaySize = circleDiameter
        let imgAspect = imageSize.width / imageSize.height

        // scaledToFill ile görüntülenen boyut
        let displayW: CGFloat
        let displayH: CGFloat
        if imgAspect > 1 {
            displayH = displaySize
            displayW = displaySize * imgAspect
        } else {
            displayW = displaySize
            displayH = displaySize / imgAspect
        }

        // Toplam ölçeklenmiş boyut
        let scaledW = displayW * scale
        let scaledH = displayH * scale

        // Piksel başına oran
        let pxPerPtX = imageSize.width / scaledW
        let pxPerPtY = imageSize.height / scaledH

        // Daire merkezinin, ölçeklenmiş görüntü üzerindeki konumu
        let cropCenterX = (scaledW / 2 - offset.width) * pxPerPtX
        let cropCenterY = (scaledH / 2 - offset.height) * pxPerPtY
        let cropSize = displaySize * pxPerPtX

        let cropRect = CGRect(
            x: cropCenterX - cropSize / 2,
            y: cropCenterY - cropSize / 2,
            width: cropSize,
            height: cropSize
        ).intersection(CGRect(origin: .zero, size: imageSize))

        // Kırp (kare olarak — daire maskeleme client tarafında yapılıyor)
        guard let cgImage = image.cgImage?.cropping(to: cropRect) else { return image }
        return UIImage(cgImage: cgImage, scale: image.scale, orientation: image.imageOrientation)
    }
}

// MARK: - Crop Overlay Shape

/// Ekranın tamamını kaplayıp, ortadaki daire bölümünü açık bırakan şekil.
private struct CropOverlay: Shape {
    let circleCenter: CGPoint
    let circleDiameter: CGFloat

    func path(in rect: CGRect) -> Path {
        var path = Path(rect)
        let circleRect = CGRect(
            x: circleCenter.x - circleDiameter / 2,
            y: circleCenter.y - circleDiameter / 2,
            width: circleDiameter,
            height: circleDiameter
        )
        path.addEllipse(in: circleRect)
        return path
    }
}

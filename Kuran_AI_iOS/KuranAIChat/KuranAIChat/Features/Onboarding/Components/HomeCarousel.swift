import SwiftUI
import Combine
import CoreLocation
import UIKit

struct HomeCarousel: View {
    let language: Language
    let onPrimaryAction: (String?) -> Void

    private static let cardHeight: CGFloat = 200

    var body: some View {
        TabView {
            VerseOfDayCard(language: language)
                .frame(height: Self.cardHeight)
            PrayerTimesCard(language: language)
                .frame(height: Self.cardHeight)
            QuickQuestionsCard(language: language, onTap: onPrimaryAction)
                .frame(height: Self.cardHeight)
        }
        .tabViewStyle(.page(indexDisplayMode: .automatic))
        .frame(height: Self.cardHeight + 24)
    }
}

private struct VerseOfDayCard: View {
    let language: Language
    @StateObject private var vm = VerseOfDayViewModel()

    // Fallback ayetler (API'ye bağlanılamazsa)
    private static let fallbackVerses: [(surah: Int, ayah: Int, arabic: String, textTr: String, textEn: String, nameTr: String, nameEn: String)] = [
        (1, 1, "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ", "Bismillâhirrahmânirrahîm", "In the name of Allah, the Entirely Merciful, the Especially Merciful.", "Fatiha", "Al-Fatiha"),
        (2, 255, "اللَّهُ لَا إِلَٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ...", "Allah... O'nun kürsüsü gökleri ve yeri kaplar.", "Allah... His Kursi extends over the heavens and the earth.", "Bakara", "Al-Baqarah"),
        (55, 26, "كُلُّ مَنْ عَلَيْهَا فَانٍ", "Yeryüzünde bulunan her şey fânidir.", "What is on earth will perish.", "Rahman", "Ar-Rahman"),
        (57, 4, "هُوَ الَّذِي خَلَقَ السَّمَاوَاتِ وَالْأَرْضَ فِي سِتَّةِ أَيَّامٍ", "O, gökleri ve yeri altı günde yarattı.", "It is He who created the heavens and earth in six days.", "Hadid", "Al-Hadid"),
        (112, 1, "قُلْ هُوَ اللَّهُ أَحَدٌ", "De ki: O Allah birdir.", "Say: He is Allah, the One.", "İhlas", "Al-Ikhlas"),
    ]

    private var fallbackVerse: (arabic: String, translation: String, ref: String) {
        let dayOfYear = Calendar.current.ordinality(of: .day, in: .year, for: Date()) ?? 1
        let v = Self.fallbackVerses[dayOfYear % Self.fallbackVerses.count]
        return (
            arabic: v.arabic,
            translation: language.localized(v.textTr, v.textEn),
            ref: language.localized("\(v.nameTr) \(v.ayah)", "\(v.nameEn) \(v.ayah)")
        )
    }

    private var resolvedVerse: (arabic: String, translation: String, ref: String) {
        if let verse = vm.verse, let t = verse.translations?.first {
            return (verse.arabicText, t.translationText, "\(t.sourceName) \(verse.ayahNumber)")
        }
        return fallbackVerse
    }

    var body: some View {
        let v = resolvedVerse

        CarouselCard {
            VStack(alignment: .leading, spacing: 10) {
                Text(language.localized("Günün Ayeti", "Verse of the Day"))
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundColor(Theme.neonPrimary)
                Text(v.arabic)
                    .font(.system(size: 16, weight: .medium))
                    .foregroundColor(Theme.textPrimary)
                    .lineLimit(2)
                Text("\"\(v.translation)\"")
                    .font(.system(size: 13, weight: .regular))
                    .italic()
                    .foregroundColor(Theme.textPrimary.opacity(0.9))
                    .lineLimit(2)
                Spacer(minLength: 4)
                HStack {
                    Spacer()
                    Text(v.ref)
                        .font(.system(size: 13, weight: .bold))
                        .italic()
                        .foregroundColor(Theme.neonPrimary)
                }
            }
        }
        .onAppear {
            vm.fetchIfNeeded(lang: language == .turkish ? "tr" : "en")
        }
        .onChange(of: language) { _, newLang in
            vm.fetch(lang: newLang == .turkish ? "tr" : "en")
        }
    }
}

private final class VerseOfDayViewModel: ObservableObject {
    @Published var verse: DailyVerseDTO?
    private var hasFetched = false

    func fetchIfNeeded(lang: String) {
        guard !hasFetched else { return }
        fetch(lang: lang)
    }

    func fetch(lang: String) {
        hasFetched = true
        Task { @MainActor in
            do {
                let result = try await VerseService.shared.getVerseOfDay(lang: lang)
                self.verse = result
            } catch {
                // API'den çekilemezse fallback kullanılır (verse nil kalır)
            }
        }
    }
}

private struct PrayerTimesCard: View {
    let language: Language
    @StateObject private var prayerVM = PrayerTimesViewModel()

    var body: some View {
        CarouselCard {
            VStack(alignment: .leading, spacing: 10) {
                HStack {
                    Text(language.localized("Namaz Vakitleri", "Prayer Times"))
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundColor(Theme.neonPrimary)
                    Spacer()
                    if !prayerVM.cityName.isEmpty {
                        HStack(spacing: 4) {
                            Image(systemName: "mappin.circle.fill")
                                .font(.system(size: 13))
                                .foregroundColor(Theme.neonPrimary.opacity(0.7))
                            Text(prayerVM.cityName)
                                .font(.system(size: 13, weight: .medium))
                                .foregroundColor(Theme.textPrimary.opacity(0.85))
                        }
                    }
                }
                if let display = prayerVM.displayRows(language: language) {
                    VStack(alignment: .leading, spacing: 6) {
                        ForEach(display) { row in
                            PrayerRow(name: row.name, time: row.time, isNext: row.isNext)
                        }
                    }
                } else {
                    VStack(alignment: .leading, spacing: 6) {
                        ForEach(0..<6, id: \.self) { _ in
                            SkeletonRow()
                        }
                    }
                }
            }
        }
        .onAppear { prayerVM.fetchIfNeeded() }
        .onReceive(NotificationCenter.default.publisher(for: UIApplication.willEnterForegroundNotification)) { _ in
            prayerVM.fetchIfNeeded()
        }
    }
}

private struct SkeletonRow: View {
    @State private var opacity: Double = 0.4
    var body: some View {
        HStack {
            RoundedRectangle(cornerRadius: 4)
                .fill(Theme.textMuted)
                .frame(width: 50, height: 10)
            Spacer()
            RoundedRectangle(cornerRadius: 4)
                .fill(Theme.textMuted.opacity(0.6))
                .frame(width: 36, height: 10)
        }
        .opacity(opacity)
        .onAppear {
            withAnimation(.easeInOut(duration: 0.8).repeatForever(autoreverses: true)) {
                opacity = 0.7
            }
        }
    }
}

private struct PrayerRow: View {
    let name: String
    let time: String
    var isNext: Bool = false
    var body: some View {
        HStack {
            Text(name).font(.system(size: 13)).foregroundColor(Theme.textMuted)
            Spacer()
            Text(time)
                .font(.system(size: 13, weight: isNext ? .semibold : .medium))
                .foregroundColor(isNext ? Theme.neonPrimary : Theme.textPrimary)
        }
    }
}

private struct QuickQuestionsCard: View {
    let language: Language
    let onTap: (String?) -> Void

    private let questions: [(tr: String, en: String)] = [
        ("Fatiha suresi tefsirini anlat", "Explain the tafsir of Surah Fatiha"),
        ("Ayetel Kürsi ne anlatıyor?", "What does Ayat al-Kursi describe?"),
        ("Namaz ile ilgili hadisler", "Hadiths about prayer"),
    ]

    var body: some View {
        CarouselCard {
            VStack(alignment: .leading, spacing: 12) {
                Text(language.localized("Hızlı Sorular", "Quick Questions"))
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundColor(Theme.neonPrimary)
                VStack(spacing: 8) {
                    ForEach(Array(questions.enumerated()), id: \.offset) { _, q in
                        Button {
                            onTap(language.localized(q.tr, q.en))
                        } label: {
                            HStack {
                                Text(language.localized(q.tr, q.en))
                                    .font(.system(size: 13))
                                    .foregroundColor(Theme.textPrimary)
                                    .lineLimit(1)
                                    .multilineTextAlignment(.leading)
                                Spacer()
                                Image(systemName: "chevron.right")
                                    .font(.system(size: 13))
                                    .foregroundColor(Theme.textMuted)
                            }
                            .padding(.vertical, 8)
                            .padding(.horizontal, 10)
                            .background(Theme.pillColor.opacity(0.5))
                            .cornerRadius(10)
                        }
                        .buttonStyle(.plain)
                    }
                }
            }
        }
    }
}

private struct CarouselCard<Content: View>: View {
    @ViewBuilder let content: Content
    var body: some View {
        content
            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topLeading)
            .padding(16)
            .background(
                RoundedRectangle(cornerRadius: 20)
                    .fill(Theme.glassGradient)
                    .overlay(RoundedRectangle(cornerRadius: 20).stroke(Theme.border, lineWidth: 1))
            )
    }
}

// MARK: - Prayer Times ViewModel

private struct PrayerRowDisplay: Identifiable {
    let id = UUID()
    let name: String
    let time: String
    let isNext: Bool
}

private final class PrayerTimesViewModel: NSObject, ObservableObject {
    @Published var prayerData: PrayerTimesData?
    @Published var isLoading = false
    @Published var error: String?
    @Published var nowTick = Date()
    @Published var cityName: String = ""

    private let locationManager = CLLocationManager()
    private let geocoder = CLGeocoder()
    private var lastFetchedLocation: (lat: Double, lng: Double)?
    private var hasRequestedPermission = false
    private var hasFetchedData = false
    private var timer: Timer?

    override init() {
        super.init()
        startTimer()
    }

    deinit {
        timer?.invalidate()
    }

    private func startTimer() {
        timer?.invalidate()
        timer = Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { [weak self] _ in
            self?.nowTick = Date()
        }
    }

    func displayRows(language: Language) -> [PrayerRowDisplay]? {
        guard let times = prayerData?.times else { return nil }
        let order: [(key: String, nameTr: String, nameEn: String)] = [
            ("imsak", "İmsak", "Fajr"),
            ("gunes", "Güneş", "Sunrise"),
            ("ogle", "Öğle", "Dhuhr"),
            ("ikindi", "İkindi", "Asr"),
            ("aksam", "Akşam", "Maghrib"),
            ("yatsi", "Yatsı", "Isha"),
        ]
        let calendar = Calendar.current
        let components = calendar.dateComponents([.year, .month, .day], from: Date())
        let today = calendar.date(from: components) ?? Date()
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        formatter.timeZone = TimeZone.current
        let now = nowTick
        var nextIndex = -1
        var minFuture: Date?
        var rows: [PrayerRowDisplay] = []
        for (index, item) in order.enumerated() {
            let timeStr: String = {
                switch item.key {
                case "imsak": return times.imsak
                case "gunes": return times.gunes
                case "ogle": return times.ogle
                case "ikindi": return times.ikindi
                case "aksam": return times.aksam
                case "yatsi": return times.yatsi
                default: return nil
                }
            }() ?? "--:--"
            let name = language.localized(item.nameTr, item.nameEn)
            guard let prayerDate = formatter.date(from: String(timeStr.prefix(5))) else {
                rows.append(PrayerRowDisplay(name: name, time: timeStr, isNext: false))
                continue
            }
            var comps = calendar.dateComponents([.year, .month, .day], from: today)
            comps.hour = calendar.component(.hour, from: prayerDate)
            comps.minute = calendar.component(.minute, from: prayerDate)
            let fullPrayerDate = calendar.date(from: comps) ?? prayerDate
            let isFuture = fullPrayerDate > now
            if isFuture && (minFuture == nil || fullPrayerDate < minFuture!) {
                minFuture = fullPrayerDate
                nextIndex = index
            }
            rows.append(PrayerRowDisplay(name: name, time: timeStr, isNext: false))
        }
        if nextIndex >= 0, nextIndex < rows.count {
            let interval = max(0, (minFuture ?? now).timeIntervalSince(now))
            let hours = Int(interval) / 3600
            let minutes = (Int(interval) % 3600) / 60
            let seconds = Int(interval) % 60
            let countdownTr: String
            let countdownEn: String
            if hours > 0 {
                countdownTr = "\(hours) Saat \(minutes) dk \(seconds) sn"
                countdownEn = "\(hours)h \(minutes)m \(seconds)s"
            } else if minutes > 0 {
                countdownTr = "\(minutes) dk \(seconds) sn"
                countdownEn = "\(minutes)m \(seconds)s"
            } else {
                countdownTr = "\(seconds) sn"
                countdownEn = "\(seconds)s"
            }
            let countdown = language.localized(countdownTr, countdownEn)
            let timeStr = rows[nextIndex].time
            rows[nextIndex] = PrayerRowDisplay(
                name: rows[nextIndex].name,
                time: "\(countdown) (\(timeStr))",
                isNext: true
            )
        }
        return rows
    }

    func fetchIfNeeded() {
        locationManager.delegate = self
        let status = locationManager.authorizationStatus

        switch status {
        case .authorizedWhenInUse, .authorizedAlways:
            // İzin var → konum iste
            locationManager.requestLocation()
        case .denied, .restricted:
            // İzin reddedilmiş → sessizce İstanbul göster
            if !hasFetchedData {
                cityName = "İstanbul"
                Task { @MainActor in await fetchFromBackend(lat: nil, lng: nil) }
            }
        case .notDetermined:
            // Henüz karar verilmemiş → izin iste, cevabı delegate'te bekle
            if !hasRequestedPermission {
                hasRequestedPermission = true
                locationManager.requestWhenInUseAuthorization()
            }
            // Kullanıcı cevap verene kadar skeleton göster, fallback'e düşme
        @unknown default:
            if !hasFetchedData {
                cityName = "İstanbul"
                Task { @MainActor in await fetchFromBackend(lat: nil, lng: nil) }
            }
        }
    }

    private func reverseGeocode(location: CLLocation) {
        geocoder.reverseGeocodeLocation(location) { [weak self] placemarks, _ in
            guard let self = self else { return }
            Task { @MainActor in
                if let city = placemarks?.first?.locality, !city.isEmpty {
                    self.cityName = city
                } else if let admin = placemarks?.first?.administrativeArea, !admin.isEmpty {
                    self.cityName = admin
                }
            }
        }
    }

    private func fetchFromBackend(lat: Double?, lng: Double?) async {
        isLoading = true
        defer { isLoading = false }
        do {
            let data = try await PrayerTimesService.shared.getPrayerTimes(lat: lat, lng: lng)
            self.prayerData = data
            self.error = nil
            self.hasFetchedData = true
        } catch {
            self.error = error.localizedDescription
        }
    }
}

extension PrayerTimesViewModel: CLLocationManagerDelegate {
    // Kullanıcı izin verdi/reddetti → bu delegate çağrılır
    nonisolated func locationManagerDidChangeAuthorization(_ manager: CLLocationManager) {
        Task { @MainActor in
            let status = manager.authorizationStatus
            switch status {
            case .authorizedWhenInUse, .authorizedAlways:
                // Kullanıcı izin verdi → hemen konum iste
                manager.requestLocation()
            case .denied, .restricted:
                // Kullanıcı reddetti → sessizce İstanbul göster
                if !hasFetchedData {
                    cityName = "İstanbul"
                    await fetchFromBackend(lat: nil, lng: nil)
                }
            case .notDetermined:
                break // Hâlâ karar vermedi, beklemeye devam
            @unknown default:
                if !hasFetchedData {
                    cityName = "İstanbul"
                    await fetchFromBackend(lat: nil, lng: nil)
                }
            }
        }
    }

    nonisolated func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        guard let loc = locations.last else { return }
        Task { @MainActor in
            lastFetchedLocation = (loc.coordinate.latitude, loc.coordinate.longitude)
            reverseGeocode(location: loc)
            await fetchFromBackend(lat: loc.coordinate.latitude, lng: loc.coordinate.longitude)
        }
    }

    nonisolated func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
        Task { @MainActor in
            // Konum alınamadı ama sessizce İstanbul'a düş
            if !hasFetchedData {
                cityName = "İstanbul"
                await fetchFromBackend(lat: nil, lng: nil)
            }
        }
    }
}

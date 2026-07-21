package prayer

// PrayerTimesRequest namaz vakitleri için query parametreleri.
// lat/lng verilmezse İstanbul kullanılır.
type PrayerTimesRequest struct {
	Lat *float64 `form:"lat"`
	Lng *float64 `form:"lng"`
}

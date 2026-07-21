package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rampesna/kuran-ai-gateway/core"
	prayerRequest "github.com/rampesna/kuran-ai-gateway/requests/prayer"
)

const (
	aladhanBaseURL = "https://api.aladhan.com/v1"
	istanbulLat    = 41.0082
	istanbulLng    = 28.9784
	turkeyMethod   = 13 // Diyanet İşleri Başkanlığı, Turkey
)

// AladhanTimingsResponse Aladhan API timings response
type AladhanTimingsResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings struct {
			Imsak   string `json:"Imsak"`
			Fajr    string `json:"Fajr"`
			Sunrise string `json:"Sunrise"`
			Dhuhr   string `json:"Dhuhr"`
			Asr     string `json:"Asr"`
			Maghrib string `json:"Maghrib"`
			Isha    string `json:"Isha"`
		} `json:"timings"`
		Date struct {
			Readable string `json:"readable"`
			Gregorian struct {
				Date string `json:"date"` // DD-MM-YYYY
			} `json:"gregorian"`
		} `json:"date"`
		Meta struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Timezone  string  `json:"timezone"`
			Method    struct {
				Name string `json:"name"`
			} `json:"method"`
		} `json:"meta"`
	} `json:"data"`
}

// PrayerTimesService namaz vakitleri proxy servisi (Aladhan API)
type PrayerTimesService struct {
	client  *http.Client
	baseURL string
}

// NewPrayerTimesService PrayerTimesService oluşturur
func NewPrayerTimesService() *PrayerTimesService {
	return &PrayerTimesService{
		client:  &http.Client{Timeout: 15 * time.Second},
		baseURL: aladhanBaseURL,
	}
}

// GetPrayerTimes konuma göre namaz vakitlerini getirir; konum yoksa İstanbul kullanılır
func (s *PrayerTimesService) GetPrayerTimes(ctx context.Context, req *prayerRequest.PrayerTimesRequest) *core.ServiceResponse {
	lat := istanbulLat
	lng := istanbulLng
	if req.Lat != nil && req.Lng != nil {
		lat = *req.Lat
		lng = *req.Lng
	}

	now := time.Now()
	dateStr := now.Format("02-01-2006") // DD-MM-YYYY for Aladhan

	url := fmt.Sprintf("%s/timings/%s?latitude=%.6f&longitude=%.6f&method=%d",
		s.baseURL, dateStr, lat, lng, turkeyMethod)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return core.NewServiceResponse(false, "request build failed", 500, nil)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return core.NewServiceResponse(false, "upstream request failed", 502, nil)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.NewServiceResponse(false, "upstream error", 502, nil)
	}

	var aladhanResp AladhanTimingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&aladhanResp); err != nil {
		return core.NewServiceResponse(false, "upstream decode failed", 502, nil)
	}

	if aladhanResp.Code != 200 || aladhanResp.Data.Timings.Imsak == "" {
		return core.NewServiceResponse(false, "invalid times data", 502, nil)
	}

	t := &aladhanResp.Data.Timings
	// iOS ve mevcut sözleşme: place (opsiyonel), times (imsak, gunes, ogle, ikindi, aksam, yatsi), date (YYYY-MM-DD)
	dateISO := now.Format("2006-01-02")
	place := map[string]interface{}{
		"name":       aladhanResp.Data.Meta.Timezone,
		"country":    "Turkey",
		"stateName":  nil,
	}
	if aladhanResp.Data.Meta.Latitude != 0 {
		place["latitude"] = aladhanResp.Data.Meta.Latitude
		place["longitude"] = aladhanResp.Data.Meta.Longitude
	}

	data := map[string]interface{}{
		"place": place,
		"times": map[string]string{
			"imsak":   t.Imsak,
			"gunes":   t.Sunrise,
			"ogle":    t.Dhuhr,
			"ikindi":  t.Asr,
			"aksam":   t.Maghrib,
			"yatsi":   t.Isha,
		},
		"date": dateISO,
	}
	return core.NewServiceResponse(true, "ok", 200, data)
}

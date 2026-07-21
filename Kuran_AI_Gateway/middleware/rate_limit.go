package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
)

type rateLimitEntry struct {
	count    int
	resetAt  time.Time
}

// RateLimiter belirli bir sürede IP başına maksimum istek sayısını sınırlar.
// maxRequests: pencere başına izin verilen istek sayısı
// window: zaman penceresi süresi
func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	entries := make(map[string]*rateLimitEntry)

	// Arka planda eski kayıtları temizle (memory leak önleme)
	go func() {
		for {
			time.Sleep(window * 2)
			mu.Lock()
			now := time.Now()
			for ip, entry := range entries {
				if now.After(entry.resetAt) {
					delete(entries, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		entry, exists := entries[ip]
		now := time.Now()

		if !exists || now.After(entry.resetAt) {
			entries[ip] = &rateLimitEntry{count: 1, resetAt: now.Add(window)}
			mu.Unlock()
			c.Next()
			return
		}

		entry.count++
		if entry.count > maxRequests {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, core.NewServiceResponse(false, "Too many requests. Please try again later.", 429, nil))
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}

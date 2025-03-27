package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientData struct {
    lastRequest time.Time
    count       int
}

var rateLimitStore = make(map[string]*clientData)
var mu sync.Mutex

func RateLimiter(maxPerMinute int) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        mu.Lock()
        data, exists := rateLimitStore[ip]
        now := time.Now()
        if !exists || now.Sub(data.lastRequest) > time.Minute {
            rateLimitStore[ip] = &clientData{lastRequest: now, count: 1}
        } else {
            data.count++
            data.lastRequest = now
            if data.count > maxPerMinute {
                mu.Unlock()
                c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
                return
            }
        }
        mu.Unlock()
        c.Next()
    }
}
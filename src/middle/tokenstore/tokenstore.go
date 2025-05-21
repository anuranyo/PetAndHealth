package tokenstore

import (
	"sync"
	"time"
)

var (
	blacklistedTokens = make(map[string]time.Time)
	mutex             = &sync.Mutex{}
)

func AddToBlacklist(token string, expiration time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	blacklistedTokens[token] = expiration
}

func IsBlacklisted(token string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := blacklistedTokens[token]
	return exists
}

func CleanupExpiredTokens() {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()
	for token, expiration := range blacklistedTokens {
		if now.After(expiration) {
			delete(blacklistedTokens, token)
		}
	}
}

func GetTokenCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return len(blacklistedTokens)
}

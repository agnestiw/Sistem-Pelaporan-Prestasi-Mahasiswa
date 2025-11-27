package memory

import (
	"sync"
)

var tokenBlacklist sync.Map

func AddToBlacklist(token string) {
	tokenBlacklist.Store(token, true)
}

func IsBlacklisted(token string) bool {
	_, found := tokenBlacklist.Load(token)
	return found
}
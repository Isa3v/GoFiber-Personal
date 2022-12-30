package handlers

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type HandlerApi struct {
	cache *cache.Cache // Инициализированный кеш
}

// Инстанс для api
func New() *HandlerApi {
	return &HandlerApi{
		cache: cache.New(30*time.Minute, 2*time.Hour),
	}
}

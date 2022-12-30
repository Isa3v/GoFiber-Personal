package controller

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type ControllerApi struct {
	cache *cache.Cache // Инициализированный кеш
}

// Инстанс для api
func New() *ControllerApi {
	return &ControllerApi{
		cache: cache.New(30*time.Minute, 2*time.Hour),
	}
}

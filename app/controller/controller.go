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
		cache: cache.New(5*time.Minute, 1*time.Hour),
	}
}

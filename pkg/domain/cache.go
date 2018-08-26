package domain

import (
	"time"
)

type CacheService interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte) error
	Invalidate(key string) error
	SetWithExpiry(key string, content []byte, d time.Duration) error
}

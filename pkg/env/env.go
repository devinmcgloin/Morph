package env

import "os"

// Getenv is the same as os.Getenv but allows for secondary key.
func Getenv(key string, opt string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return opt
}

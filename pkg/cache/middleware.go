package cache

import (
	"net/http"

	"net/http/httptest"

	"time"

	"log"

	"github.com/fokal/fokal/pkg/handler"
)

func Handler(state *handler.State, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		b, err := Get(state.RD, url)
		if err != nil {
			c := httptest.NewRecorder()
			next.ServeHTTP(c, r)

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}

			content := c.Body.Bytes()
			w.WriteHeader(c.Code)
			w.Write(content)

			if c.Code == http.StatusOK {
				log.Printf("Setting Handler URL: %s\n", url)
				Set(state.RD, url, content)
				ExpireAt(state.RD, url, time.Minute*15)
			}
			return
		} else {
			log.Printf("Retrieving Handler URL: %s\n", url)
			w.Write(b)
			w.WriteHeader(http.StatusOK)
			return
		}
	})
}

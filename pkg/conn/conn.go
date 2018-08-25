package conn

import (
	"log"
	"time"

	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/vision/v1"
	"googlemaps.github.io/maps"
)

func DialPostgres(postgresURL string) (db *sqlx.DB) {
	var err error

	db, err = sqlx.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func DialRedis(server string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   15,
		Wait:        true,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(server)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if _, err := c.Do("PING"); err != nil {
				c.Close()
				log.Println(err)
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func DialGoogleServices(apiKey string) (*vision.Service, *maps.Client, error) {
	var err error

	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	visionService, err := vision.New(client)
	if err != nil {
		log.Fatal(err)
	}

	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	return visionService, mapsClient, nil
}

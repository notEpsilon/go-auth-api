package sessions

import (
	"go-auth/constants"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

var RSS *session.Store

func MustInitRedisStore() {
	storage := redis.New(redis.Config{
		URL: os.Getenv("REDIS_URL"),
	})

	RSS = session.New(session.Config{
		CookieHTTPOnly: true,
		Expiration:     7 * 24 * time.Hour,
		CookieSameSite: "Lax",
		Storage:        storage,
		CookieSecure:   os.Getenv("APP_ENV") == string(constants.Production),
	})
}

package cfg

import (
	"os"
)

var (
	PayPal struct {
		ClientId  string
		ApiSecret string
	}

	PostgreSQL struct {
		Database string
		User     string
		Password string
	}
)

func init() {
	PayPal.ClientId = getenv("PAYPAL_CLIENT_ID", "")
	PayPal.ApiSecret = getenv("PAYPAL_API_SECRET", "")

	PostgreSQL.Database = getenv("DATABASE_SERVICE_DB", "app")
	PostgreSQL.User = getenv("DATABASE_SERVICE_USER", "root")
	PostgreSQL.Password = getenv("DATABASE_SERVICE_PASSWORD", "")
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

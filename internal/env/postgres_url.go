package env

import (
	"net/url"
	"os"
)

// PostgresURL returns a postgres DSN. Uses DB_ADDR if set; otherwise builds from
// POSTGRES_USER, POSTGRES_PASSWORD, DB_HOST, DB_PORT, POSTGRES_DB (password URL-encoded).
func PostgresURL() string {
	if addr, ok := os.LookupEnv("DB_ADDR"); ok && addr != "" {
		return addr
	}

	user := GetString("POSTGRES_USER", "app")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := GetString("DB_HOST", "localhost")
	port := GetString("DB_PORT", "5432")
	db := GetString("POSTGRES_DB", "currency_exchange")
	ssl := GetString("DB_SSLMODE", "disable")

	u := &url.URL{
		Scheme: "postgres",
		Host:   host + ":" + port,
		Path:   "/" + db,
	}
	if user != "" || pass != "" {
		u.User = url.UserPassword(user, pass)
	}
	q := u.Query()
	q.Set("sslmode", ssl)
	u.RawQuery = q.Encode()

	return u.String()
}

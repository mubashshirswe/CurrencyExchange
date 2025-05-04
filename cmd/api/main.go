package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/mubashshir3767/currencyExchange/internal/db"
	"github.com/mubashshir3767/currencyExchange/internal/env"
	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/store/cache"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	// Load configuration from environment variables
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://qweqweqeqe12312314:qweqweqeqe12312314@localhost:5432/currency_exchange?sslmode=disable"),
			maxOpenConns: env.GetInt("MAX_OPEN_CONNS", 50),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 50),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisConfig: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", true),
		},
		env: env.GetString("ENV", "PROD"),
	}

	rdb := cache.NewRedisClient(cfg.redisConfig.addr, cfg.redisConfig.pw, cfg.redisConfig.db)
	log.Println("redis cache connection established")

	// Initialize database connection
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Fatalf("failed to establish a database connection: %v", err)
	}

	defer db.Close()
	log.Println("DATABASE HAS BEEN SUCCESSFULLY ESTABLISHED")

	// Initialize store
	store := store.NewStorage(db)
	cacheStore := cache.NewRedisStorage(rdb)

	// Initialize the application
	app := application{
		config:     cfg,
		store:      store,
		cacheStore: cacheStore,
	}

	// Mount routes and run the server
	mux := app.mount()
	log.Fatal(app.run(mux))
}

func GetIdFromContext(r *http.Request) int64 {
	param := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

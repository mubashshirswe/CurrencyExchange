package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mubashshir3767/currencyExchange/internal/db"
	"github.com/mubashshir3767/currencyExchange/internal/env"
	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr        string
	db          dbConfig
	redisConfig redisConfig
	env         string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			app.writeResponse(w, http.StatusOK, "ok")
		})
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("server has been started on %v env %v", app.config.addr, app.config.env)
	return srv.ListenAndServe()
}

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

	// Initialize the application
	app := application{
		config: cfg,
		store:  store,
	}

	// Mount routes and run the server
	mux := app.mount()
	log.Fatal(app.run(mux))
}

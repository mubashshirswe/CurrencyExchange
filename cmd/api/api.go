package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/store/cache"
)

type application struct {
	config     config
	store      store.Storage
	cacheStore cache.Storage
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

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			app.writeResponse(w, http.StatusOK, "ok")
		})

		r.Post("/register", app.CreateUserHandler)
		r.Post("/login", app.CreateUserHandler)

		r.With(app.JWTUserMiddleware()).Route("/", func(r chi.Router) {

			r.Route("/users", func(r chi.Router) {
				r.Get("/all", app.GetAllUserHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", app.UpdateUserHandler)
					r.Delete("/", app.DeleteUserHandler)
				})
			})

			r.Route("/balances", func(r chi.Router) {
				r.Post("/", app.CreateBalanceHandler)
				r.Get("/all", app.GetAllBalanceHandler)
				r.Get("/user/{id}", app.GetBalanceByUserIdHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetBalanceByIdHandler)
					r.Put("/", app.UpdateBalanceHandler)
				})
			})

			r.Route("/currencies", func(r chi.Router) {
				r.Post("/", app.CreateCurrencyHandler)
				r.Get("/all", app.GetAllCurrencyHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", app.UpdateCurrencyHandler)
					r.Delete("/", app.DeleteCurrencyHandler)
				})
			})

			r.Route("/cities", func(r chi.Router) {
				r.Post("/", app.CreateCityHandler)
				r.Get("/all", app.GetAllCityHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", app.UpdateCityHandler)
					r.Delete("/", app.DeleteCityHandler)
				})
			})

			r.Route("/transactions", func(r chi.Router) {
				r.Post("/", app.CreateTransactionHandler)
				r.Get("/all", app.GetAllTransactionHandler)
				r.Get("/all/date", app.GetAllTransactionByDateHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetTransactionByIdHandler)
					r.Put("/", app.UpdateTransactionHandler)
				})
			})

			r.Route("/companies", func(r chi.Router) {
				r.Post("/", app.CreateCompanyHandler)
				r.Get("/all", app.GetAllCompanyHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetCompanyByIdHandler)
					r.Put("/", app.UpdateCompanyHandler)
					r.Delete("/", app.DeleteCompanyHandler)
				})
			})
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

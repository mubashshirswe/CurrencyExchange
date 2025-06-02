package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mubashshir3767/currencyExchange/internal/service"
	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/store/cache"
)

type application struct {
	config     config
	store      store.Storage
	service    service.Service
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
	// r.Use(chi.MiddlewareFunc(http.StripPrefix))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		// r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		// 	app.writeResponse(w, http.StatusOK, chi.URLParam(r, "id"))
		// })

		r.Post("/users/register", app.CreateUserHandler)
		r.Post("/users/login", app.LoginUserHandler)

		r.With(app.JWTUserMiddleware()).Route("/user", func(r chi.Router) {

			r.Get("/all", app.GetAllUserHandler)
			r.Put("/{id}", app.UpdateUserHandler)
			r.Delete("/{id}", app.DeleteUserHandler)

			r.Route("/balances", func(r chi.Router) {
				r.Post("/", app.CreateBalanceHandler)
				r.Get("/all", app.GetAllBalanceHandler)
				r.Get("/user/{id}", app.GetBalanceByUserIdHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetBalanceByIdHandler)
					r.Put("/", app.UpdateBalanceHandler)
					r.Delete("/", app.DeleteBalanceHandler)
				})
			})

			r.Route("/balances-records", func(r chi.Router) {
				r.Post("/", app.CreateBalanceRecordHandler)
				r.Get("/balance/{id}", app.GetBalanceRecordsByBalanceIdHandler)
				r.Get("/user/{id}", app.GetBalanceRecordsByUserIdHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", app.UpdateBalanceRecordHandler)
					r.Delete("/", app.DeleteBalanceRecordHandler)
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

			r.Route("/debtors", func(r chi.Router) {
				r.Post("/", app.CreateDebtorsHandler)
				r.Get("/user/{id}", app.GetDebtorsByUserIdHandler)
				r.Get("/receive/{id}", app.ReceivedDebtHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetDebtorsByIdHandler)
					r.Put("/", app.UpdateDebtorsHandler)
					r.Delete("/", app.DeleteDebtorsHandler)
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
				r.Post("/create", app.CreateTransactionHandler)
				r.Get("/complete/{serial_no}", app.CompleteTransactionHandler)
				r.Get("/all/balance/{id}", app.GetAllTransactionByBalanceIdHandler)
				r.Get("/all/user/{id}", app.GetAllTransactionByUserIdHandler)
				r.Get("/all/receiver/{id}", app.GetAllTransactionByReceiverIdHandler)
				r.Get("/all/date", app.GetAllTransactionByDateHandler)
				r.Get("/all/active/{status}", app.GetAllActiveTransactionsHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetTransactionByIdHandler)
					r.Put("/", app.UpdateTransactionHandler)
					r.Delete("/", app.DeleteTransactionHandler)
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

func getIDFromContext(r *http.Request) int64 {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 60)
	if err != nil {
		return 0
	}

	return ID
}

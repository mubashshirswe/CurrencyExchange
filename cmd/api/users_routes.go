package main

import "github.com/go-chi/chi/v5"

func (app *application) UserRoutes(r chi.Router) {
	r.With(app.JWTUserMiddleware()).Route("/user", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Get("/all", app.GetAllUserHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.Put("/", app.UpdateUserHandler)
				r.Delete("/", app.DeleteUserHandler)
			})
		})

		// r.Route("/employees", func(r chi.Router) {
		// 	r.Get("/register", app.CreateEmployeeHandler)
		// 	r.Get("/all", app.GetAllEmployeeHandler)
		// 	r.Route("/{id}", func(r chi.Router) {
		// 		r.Put("/", app.UpdateEmployeeHandler)
		// 		r.Delete("/", app.DeleteEmployeeHandler)
		// 	})
		// })

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
			r.Get("/all/balance/{id}", app.GetAllTransactionHandler)
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
}

package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s err: %s", r.Method, r.URL.Path, err)
	writeError(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s err: %s", r.Method, r.URL.Path, err)
	writeError(w, http.StatusBadRequest, err.Error())
}

// func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
// 	log.Printf("conflict error: %s path: %s err: %s", r.Method, r.URL.Path, err)
// 	writeError(w, http.StatusConflict, err.Error())
// }

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized error: %s path: %s err: %s", r.Method, r.URL.Path, err)
	writeError(w, http.StatusUnauthorized, "Unauthorized method used")
}

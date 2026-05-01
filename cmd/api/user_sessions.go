package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type upsertUserSessionPayload struct {
	DeviceID     string  `json:"device_id" validate:"required"`
	FCMToken     string  `json:"fcm_token" validate:"required"`
	RefreshToken *string `json:"refresh_token"`
	Platform     *string `json:"platform"`
	AppVersion   *string `json:"app_version"`
	UserAgent    *string `json:"user_agent"`
}

func (app *application) UpsertUserSessionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(int64)

	var payload upsertUserSessionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	row := &store.UserSession{
		UserID:       userID,
		DeviceID:     payload.DeviceID,
		FCMToken:     payload.FCMToken,
		RefreshToken: payload.RefreshToken,
		Platform:     payload.Platform,
		AppVersion:   payload.AppVersion,
		UserAgent:    payload.UserAgent,
	}

	if err := app.store.UserSessions.Upsert(r.Context(), row); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, row); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) ListUserSessionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(int64)

	sessions, err := app.store.UserSessions.ListByUserID(r.Context(), userID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, sessions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type patchUserSessionPayload struct {
	FCMToken     string  `json:"fcm_token" validate:"required"`
	RefreshToken *string `json:"refresh_token"`
}

func (app *application) UpdateUserSessionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(int64)
	id := getIDFromContext(r)

	var payload patchUserSessionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.UserSessions.UpdateFCM(r.Context(), id, userID, payload.FCMToken, payload.RefreshToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	updated, err := app.store.UserSessions.GetByIDForUser(r.Context(), id, userID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, updated); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteUserSessionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(int64)
	id := getIDFromContext(r)

	if err := app.store.UserSessions.Delete(r.Context(), id, userID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, id); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

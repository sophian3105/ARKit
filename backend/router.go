package main

import (
	util "aria/backend/utility"
	"encoding/json"

	"net/http"
)

// MainHandler provides the main starting handler
// for any http request.
func MainHandler() http.Handler {
	mux := http.NewServeMux()

	defaultBase := util.NewRouter(mux, "/")
	defaultBase.Use(util.LoggerMiddleware())

	defaultBase.Handle("/ping", getPing, http.MethodGet, http.MethodPost)

	// Any routes using this middleware must be fully authorized
	authBase := defaultBase.Branch("/api")
	authBase.Use(util.AuthMiddleware())
	authBase.Handle("/auth", getAuth, http.MethodGet)
	return mux
}

// Simple get method to ping the backend
func getPing(w http.ResponseWriter, r *http.Request, d *util.MiddlewareData) {
	if r.Method != http.MethodGet {
		return
	}

	w.WriteHeader(http.StatusOK)
	ginContent := map[string]string{
		"message": "pong",
	}
	json.NewEncoder(w).Encode(ginContent)
}

// Simple method to test auth
func getAuth(w http.ResponseWriter, r *http.Request, d *util.MiddlewareData) {
	w.WriteHeader(http.StatusOK)
}

package main

import (
	"aria/backend/routes"
	util "aria/backend/utility"
	"net/http"
)

// MainHandler provides the main starting handler
// for any http request.
func MainHandler() http.Handler {
	mux := http.NewServeMux()

	defaultBase := util.NewRouter(mux, "/")
	defaultBase.Use(util.LoggerMiddleware)
	defaultBase.Use(util.DatabaseMiddleware)

	defaultBase.Handle("/ping", getPing, http.MethodGet, http.MethodPost)
	defaultBase.Handle("/test", routes.PostUser, http.MethodPost)

	// Any routes using this middleware must be fully authorized
	authBase := defaultBase.Branch("/api")
	authBase.Use(util.AuthMiddleware)
	authBase.Handle("/auth", getAuth, http.MethodGet)
	authBase.Handle("/user", routes.PostUser, http.MethodPost)
	authBase.Handle("/user", routes.GetUser, http.MethodGet)
	return mux
}

// Simple get method to ping the backend
func getPing(ctx *util.Context) {
	ginContent := map[string]string{
		"message": "pong",
	}
	ctx.Json(http.StatusOK, ginContent)
}

// Simple method to test auth
func getAuth(ctx *util.Context) {
	ctx.WriteHeader(http.StatusOK)
}

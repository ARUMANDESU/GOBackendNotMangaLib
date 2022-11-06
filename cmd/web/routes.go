package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	router.GET("/", app.home)
	router.POST("/manga/create", app.AuthMiddleware(app.createManga))
	router.GET("/manga/:id", app.getManga)
	router.POST("/signup", app.signUp)
	router.POST("/signin", app.signIN)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.MiddleCORS)

	return standard.Then(router)
}

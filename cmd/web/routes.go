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
			header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
			header.Set("Access-Control-Allow-Credentials", "true")
		}
	})
	fileServer := http.FileServer(http.Dir("./public/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", app.ImageMiddleware(fileServer)))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/manga/create", app.AuthMiddleware(app.createManga))
	router.HandlerFunc(http.MethodPost, "/manga/add-chapter/:id", app.AuthMiddleware(app.addChapter))
	router.HandlerFunc(http.MethodGet, "/manga/:id", app.getManga)
	router.HandlerFunc(http.MethodGet, "/manga/:id/:v/:ch,", app.getChapter)
	router.HandlerFunc(http.MethodPost, "/signup", app.signUp)
	router.HandlerFunc(http.MethodPost, "/signin", app.signIN)
	router.HandlerFunc(http.MethodPost, "/logout", app.Logout)
	router.HandlerFunc(http.MethodGet, "/user/:id", app.isOwner(app.GetUser))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

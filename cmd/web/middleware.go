package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"notmangalib.com/internal/models"
	"strconv"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		accessToken, err := r.Cookie("AccessToken")
		log.Print(accessToken)
		if err != nil {
			w.WriteHeader(303)
			return
		}
		_, e := app.VerifyToken(accessToken.Value)
		if e != nil {
			if errors.Is(e.Err, models.ExpiredToken) {
				userId, _ := strconv.Atoi(fmt.Sprint(e.Payload["Id"]))

				fmt.Println(userId)
				token, err := app.RefreshAccessToken(e.Payload)
				if err != nil {
					http.Error(w, "internal server error", 500)
				}
				newCookie := &http.Cookie{
					Name:     "AccessToken",
					Value:    token,
					HttpOnly: true,
					MaxAge:   2592000,
				}

				http.SetCookie(w, newCookie)
			} else {
				http.Redirect(w, r, "/signUp", 303)
				return
			}
		}
		next(w, r, ps)
	}
}

func (app *application) AuthAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {

	})
}

func (app *application) MiddleCORS(next http.Handler) http.Handler {
	return func(w http.ResponseWriter,
		r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r, ps)
	}
}

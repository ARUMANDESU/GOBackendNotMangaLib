package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"notmangalib.com/internal/models"
	"strconv"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
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

func (app *application) AuthMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		next(w, r)
	})
}

func (app *application) AuthAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (app *application) ImageMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jng")
		next.ServeHTTP(w, r)
	})
}

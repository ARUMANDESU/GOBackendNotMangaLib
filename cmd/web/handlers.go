package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"notmangalib.com/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	user := models.User{1, "Arman", "", "user", ""}
	manga, err := app.manga.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	resp := make(map[string]any)
	resp["user"] = user
	resp["manga"] = manga
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (app *application) createManga(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := app.manga.Insert("OnePunchMan", "Something...", "ME", "Manga")
	if err != nil {
		app.serverError(w, err)
	}

	w.Write([]byte("Create new manga..."))

}

func (app *application) getManga(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	manga, err := app.manga.Get(id)
	if err != nil {
		app.serverError(w, err)
	}

	resp := make(map[string]any)
	resp["manga"] = manga
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

	var newUser = &models.SignModel{}
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

	user, accesstoken, err := app.SignUpService(newUser.Name, newUser.Email, newUser.Password)
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
	AccessTokenCookie := &http.Cookie{
		Name:     "AccessToken",
		Value:    accesstoken,
		MaxAge:   170000,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, AccessTokenCookie)

	resp := make(map[string]any)
	resp["user"] = user
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (app *application) signIN(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

	var user = &models.SignModel{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
	log.Println(user.Email, user.Password)

	result, accesstoken, err := app.SignINService(user.Email, user.Password)
	if err != nil {
		app.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}

	AccessTokenCookie := &http.Cookie{
		Name:     "AccessToken",
		Value:    accesstoken,
		MaxAge:   170000,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, AccessTokenCookie)

	resp := make(map[string]any)
	resp["user"] = result
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

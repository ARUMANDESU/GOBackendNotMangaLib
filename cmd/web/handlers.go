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

	w.Header().Set("Content-Type", "application/json")
	user := models.User{1, "Arman"}
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
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]any)
	resp["manga"] = manga
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

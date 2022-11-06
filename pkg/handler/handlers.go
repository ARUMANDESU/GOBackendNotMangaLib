package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"notmangalib.com/pkg/models"
	"strconv"
)

type Handler struct {
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	user := models.User{1, "Arman", "zh.arumandes@gmail.com"}
	manga, err := app.Manga.Latest()
	if err != nil {
		app.ServerError(w, err)
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

func (h *Handler) createManga(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) getManga(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {

}

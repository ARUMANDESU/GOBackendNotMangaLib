package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"notmangalib.com/internal/models"
	"os"
	"path/filepath"
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
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	var newManga = &models.MangaCreate{}
	newManga.Name = r.Form.Get("name")
	newManga.Description = r.Form.Get("description")
	newManga.Type = r.Form.Get("type")
	newManga.Status = r.Form.Get("status")
	newManga.Author = r.Form.Get("author")

	f, h, err := r.FormFile("mangaImg")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	id, err := app.manga.Insert(newManga.Name, newManga.Description, newManga.Author, newManga.Type, newManga.Status)
	if err != nil {
		app.serverError(w, err)
	}

	defer f.Close()
	path := filepath.Join(".", fmt.Sprintf("public/manga/%d", id))

	_ = os.MkdirAll(path, os.ModePerm)

	fullPath := path + "/" + "mangaImg" + filepath.Ext(h.Filename)
	staticPath := fmt.Sprintf("/static/manga/%d/", id) + "mangaImg" + filepath.Ext(h.Filename)
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	defer file.Close()

	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	log.Print(staticPath)
	err = app.manga.ChangeImg(id, staticPath)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	resp := make(map[string]any)
	resp["id"] = id
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

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

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	result, err := app.user.Get(id)

	resp := make(map[string]any)
	resp["user"] = result
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

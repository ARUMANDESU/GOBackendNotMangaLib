package pkg

import (
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"notmangalib.com/pkg/handler"
	"notmangalib.com/pkg/repository"
	"os"
)

func Routes() http.Handler {
	_ = godotenv.Load(".env")

	var db, _ = repository.OpenDB(os.Getenv("dbName"), os.Getenv("dbPassword"))

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	router.HandlerFunc(http.MethodGet, "/", handler)
	router.HandlerFunc(http.MethodPost, "/manga/create", app.createManga)
	router.HandlerFunc(http.MethodGet, "/manga/:id", app.getManga)

	standard := alice.New(app.recoverPanic, app.logRequest, handler.secureHeaders)

	return standard.Then(router)
}

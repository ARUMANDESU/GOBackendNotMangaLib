package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *main.application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *main.application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *main.application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

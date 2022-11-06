package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"notmangalib.com/internal/models"
	"runtime/debug"
	"strconv"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (app *application) NewJWT(user models.User) (string, error) {
	claims := &models.Claims{
		Name:  user.Name,
		Id:    user.Id,
		Role:  user.Role,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(250 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("SecretString"))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func (app *application) VerifyToken(accesstoken string) (jwt.MapClaims, *models.ErrorHandlerJwt) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accesstoken, claims,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, models.TokenError
			}
			return []byte("SecretString"), nil
		})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, models.HandleJWTError(claims, models.ExpiredToken)
		}
		return nil, models.HandleJWTError(nil, models.NotAuthorized)
	}
	return claims, nil
}

func (app *application) RefreshAccessToken(payload jwt.MapClaims) (string, error) {
	newUserModel := new(models.User)
	newUserModel.Name = fmt.Sprint(payload["Username"])
	newUserModel.Role = fmt.Sprint(payload["Role"])
	newUserModel.Email = fmt.Sprint(payload["Email"])
	id := fmt.Sprint(payload["Id"])
	fmt.Println(id)

	newUserModel.Id, _ = strconv.Atoi(id)
	token, err := app.NewJWT(*newUserModel)
	if err != nil {
		return "", models.TokenError
	}
	return token, nil
}

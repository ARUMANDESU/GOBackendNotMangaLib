package main

import "notmangalib.com/internal/models"

func (app *application) SignUpService(name string, email string, password string) (*models.User, string, error) {
	hashedpPassword, err := app.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	newUser, err := app.user.Insert(name, email, string(hashedpPassword))
	if err != nil {
		return nil, "", err
	}
	jwt, err := app.NewJWT(*newUser)

	return newUser, jwt, nil

}

func (app *application) SignINService(email string, password string) (*models.User, string, error) {

	user, err := app.user.FindCheckUser(email, password)
	if err != nil {
		return nil, "", err
	}

	jwt, err := app.NewJWT(*user)

	return user, jwt, nil

}

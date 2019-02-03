package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UName    string
	FName    *string
	LName    *string
	Email    string
	Password string
	Descript *string
	Role     *string
	Banned   *int
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail string `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
}

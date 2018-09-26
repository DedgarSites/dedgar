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

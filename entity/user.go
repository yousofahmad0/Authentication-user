package entity

import (
	"database/sql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName         string `json:"firstname"`
	LastName          string `json:"lastname"`
	Email             string `json:"email"`
	Username          string `json:"username"`
	PasswordHash      string
	Verification_code string
	verified          bool
	Nationality       sql.NullString `json:"nationality"`
}

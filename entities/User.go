package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" gorm:"type:varchar(60);not null"`
	Email      string `json:"email" gorm:"type:varchar(255);unique;not null"`
}

type UserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=60,alphaspace"`
	Email string `json:"email" validate:"required,lowercase,email"`
}

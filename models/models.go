package models

import (
	"gorm.io/gorm"
	"time"
)

type Laboratory struct {
	gorm.Model
	ID   int
	Info string `json:"info" gorm:"text;not null;default:null"`
}

type Article struct {
	gorm.Model
	ID          int    `gorm:"type:uuid;primary_key;"`
	Title       string `json:"title" gorm:"text;not null;default:null"`
	Description string `json:"description" gorm:"text;not null;default:null"`
	Edited      time.Time
	//Url         string `json:"url"`
}

type Photo struct {
	gorm.Model
	Data     []byte `json:"data"`
	Filename string `json:"filename"`
}

type Logo struct {
	gorm.Model
	StrData  string `json:"strdata"`
	Data     []byte `json:"data"`
	Filename string `json:"filename"`
}

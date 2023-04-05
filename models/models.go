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
	Title       string `json:"title" gorm:"text;not null;default:null"`
	Description string `json:"description" gorm:"text;not null;default:null"`
	Edited      time.Time
	Photos      []Photo
}

type Photo struct {
	gorm.Model
	Data      []byte `json:"data"`
	StrData   string `json:"strdata"`
	Filename  string `json:"filename"`
	ArticleID int    `json:"article_id,string,omitempty"`
}

type Logo struct {
	gorm.Model
	StrData  string `json:"strdata"`
	Data     []byte `json:"data"`
	Filename string `json:"filename"`
}

type User struct {
	gorm.Model
	Email string `json:"email"`
	Role  string `json:"role"`
}

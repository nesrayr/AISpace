package models

import (
	"encoding/base64"
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
}

type Photo struct {
	gorm.Model
	ID      int `gorm:"type:uuid;primary_key;"`
	Content base64.Encoding
}

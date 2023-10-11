package models

import "time"

type Word struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Word      string    `json:"word"`
	CreatedAt time.Time `json:"created_at"`
}

package models

import "time"

type Car struct {
	ID        int       `json:"id" db:"id"`
	RegNum    string    `json:"regNum" db:"reg_num"`
	Mark      string    `json:"mark" db:"mark"`
	Model     string    `json:"model" db:"model"`
	Year      int       `json:"year" db:"year"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Owner     People    `json:"owner"`
}

package models

type Car struct {
	ID      int      `json:"id" db:"id"`
	RegNums []string `json:"regNums"`
	Mark    string   `json:"mark" db:"mark"`
	Model   string   `json:"model" db:"model"`
	Year    int      `json:"year" db:"year"`
}

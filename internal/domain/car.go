package domain

import "github.com/Verce11o/effective-mobile-test/internal/models"

type CreateCarsRequest struct {
	RegNums []string `json:"regNums" binding:"required,gt=0"`
}

type GetCarsRequest struct {
	Cursor  string `form:"cursor"`
	RegNum  string `form:"regNum"`
	Mark    string `form:"mark"`
	Model   string `form:"model"`
	Year    int    `form:"year"`
	OwnerID int    `form:"owner_id"`
}

type Car struct {
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Owner  People `json:"owner"`
	RegNum string `json:"regNum"`
	Year   *int   `json:"year,omitempty"`
}

type People struct {
	Name       string  `json:"name"`
	Patronymic *string `json:"patronymic,omitempty"`
	Surname    string  `json:"surname"`
}

type CarList struct {
	Cursor string       `json:"cursor"`
	Total  int          `json:"total"`
	Cars   []models.Car `json:"cars"`
}

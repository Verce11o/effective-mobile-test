package domain

type CreateCarsRequest struct {
	RegNums []string `json:"regNums" binding:"required,gt=0"`
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

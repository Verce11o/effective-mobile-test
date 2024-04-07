package domain

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

type UpdateCarsRequest struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
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

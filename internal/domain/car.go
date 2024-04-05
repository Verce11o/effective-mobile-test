package domain

type CreateCarsRequest struct {
	RegNums []string `json:"regNums" binding:"required,gt=0"`
}

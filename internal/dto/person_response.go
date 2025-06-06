package dto

// PersonResponse represents a person with enriched information
// @Description Person information with age, gender and country
type PersonResponse struct {
	ID         int     `json:"id" example:"1"`
	Name       string  `json:"name" example:"Иван"`
	Surname    string  `json:"surname" example:"Иванов"`
	Patronymic *string `json:"patronymic,omitempty" example:"Иванович"`
	Age        *int    `json:"age,omitempty" example:"30"`
	Gender     *string `json:"gender,omitempty" example:"мужской"`
	Country    *string `json:"country,omitempty" example:"RU"`
	CreatedAt  string  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  string  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ListPersonsResponse represents paginated list of persons
// @Description Paginated list of persons with total count
type ListPersonsResponse struct {
	TotalCount int64            `json:"total" example:"100"` // total number of records
	Page       int              `json:"page" example:"1"`    // current page
	Size       int              `json:"size" example:"10"`   // page size
	Data       []PersonResponse `json:"data"`                // array of records
}

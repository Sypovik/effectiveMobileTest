package dto

// PersonResponse — DTO, который мы возвращаем в ответе на клиент.
// Содержит «обогащённую» информацию, в том числе age/gender/country.
type PersonResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
	Age        *int    `json:"age,omitempty"`
	Gender     *string `json:"gender,omitempty"`
	Country    *string `json:"country,omitempty"`
	CreatedAt  string  `json:"created_at"` // RFC3339 формат
	UpdatedAt  string  `json:"updated_at"`
}

// ListPersonsResponse — DTO для GET /people со списком и пагинацией.
type ListPersonsResponse struct {
	TotalCount int64            `json:"total"` // общее число записей
	Page       int              `json:"page"`  // текущая страница
	Size       int              `json:"size"`  // размер страницы
	Data       []PersonResponse `json:"data"`  // массив записей
}

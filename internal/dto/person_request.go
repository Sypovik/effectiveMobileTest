package dto

// CreatePersonRequest — DTO для POST /people.
// Клиент отправляет:
//
//	{
//	  "name": "Dmitriy",
//	  "surname": "Ushakov",
//	  "patronymic": "Vasilevich" // необязательно
//	}
type CreatePersonRequest struct {
	Name       string  `json:"name" binding:"required"`
	Surname    string  `json:"surname" binding:"required"`
	Patronymic *string `json:"patronymic,omitempty"`
}

// UpdatePersonRequest — DTO для PUT/PATCH /people/{id}.
// Все поля опциональны для частичного обновления (PATCH) или обязательны для PUT:
type UpdatePersonRequest struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
}

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
	Name       string  `json:"name" binding:"required" example:"Ivan"`
	Surname    string  `json:"surname" binding:"required" example:"Ivanov"`
	Patronymic *string `json:"patronymic,omitempty" example:"Ivanovich"`
}

// UpdatePersonRequest — DTO для PUT/PATCH /people/{id}.
// Все поля опциональны для частичного обновления (PATCH) или обязательны для PUT:
type UpdatePersonRequest struct {
	Name       *string `json:"name,omitempty" example:"Ivan"`
	Surname    *string `json:"surname,omitempty" example:"Ivanov"`
	Patronymic *string `json:"patronymic,omitempty" example:"Ivanovich"`
}

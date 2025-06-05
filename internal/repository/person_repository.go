package repository

import (
	"context"

	"github.com/Sypovik/effectiveMobileTest/internal/models"
)

// type Person struct {
// 	ID         int       `gorm:"primaryKey" json:"id"`
// 	Name       string    `gorm:"size:100;not null" json:"name"`
// 	Surname    string    `gorm:"size:100;not null" json:"surname"`
// 	Patronymic *string   `gorm:"size:100" json:"patronymic,omitempty"`
// 	Age        *int      `json:"age,omitempty"`
// 	Gender     *string   `gorm:"size:10" json:"gender,omitempty"`
// 	Country    *string   `gorm:"size:5" json:"country,omitempty"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// }

// PersonFilter — для фильтрации и пагинации в методе List.
type PersonFilter struct {
	Name       *string // частичное/точное совпадение по имени
	Surname    *string // частичное/точное совпадение по фамилии
	Age        *int    // равенство по возрасту (или nil, чтобы не фильтровать)
	Gender     *string
	Patronymic *string
	Country    *string
	Limit      int // число записей на страницу
	Offset     int // смещение (для пагинации)
}

// PersonRepository описывает CRUD-операции и получение списка с фильтрами.
type PersonRepository interface {
	// Create сохраняет нового человека в БД.
	Create(ctx context.Context, person *models.Person) error

	// GetByID возвращает одного человека по ID.
	GetByID(ctx context.Context, id int) (*models.Person, error)

	// Update обновляет существующую запись (все поля, кроме ID).
	Update(ctx context.Context, person *models.Person) error

	// Delete удаляет запись по ID.
	Delete(ctx context.Context, id int) error

	// List возвращает список людей по фильтрам + общее количество (для пагинации).
	List(ctx context.Context, filter PersonFilter) (results []models.Person, totalCount int64, err error)
}

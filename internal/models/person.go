package models

import "time"

// Person — внутренняя модель, соответствует таблице в Postgres.
type Person struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Surname    string    `gorm:"size:100;not null" json:"surname"`
	Patronymic *string   `gorm:"size:100" json:"patronymic,omitempty"`
	Age        *int      `json:"age,omitempty"`
	Gender     *string   `gorm:"size:10" json:"gender,omitempty"`
	Country    *string   `gorm:"size:5" json:"country,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

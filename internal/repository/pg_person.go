package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sypovik/effectiveMobileTest/internal/models"
	"gorm.io/gorm"
)

type PgPersonRepository struct {
	db *gorm.DB
}

// NewPgPersonRepository — конструктор репозитория
func NewPgPersonRepository(db *gorm.DB) PersonRepository {
	return &PgPersonRepository{db: db}
}

func (r *PgPersonRepository) Create(ctx context.Context, person *models.Person) error {
	return r.db.WithContext(ctx).Create(person).Error
}

func (r *PgPersonRepository) GetByID(ctx context.Context, id int) (*models.Person, error) {
	var person models.Person
	result := r.db.WithContext(ctx).First(&person, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &person, result.Error
}

func (r *PgPersonRepository) Update(ctx context.Context, person *models.Person) error {
	return r.db.WithContext(ctx).Save(person).Error
}

func (r *PgPersonRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Person{}, id).Error
}

func (r *PgPersonRepository) List(ctx context.Context, filter PersonFilter) (results []models.Person, totalCount int64, err error) {
	var persons []models.Person

	query := r.db.WithContext(ctx).Model(&models.Person{})

	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}
	if filter.Surname != nil {
		query = query.Where("surname = ?", *filter.Surname)
	}
	if filter.Gender != nil {
		query = query.Where("gender = ?", *filter.Gender)
	}
	if filter.Country != nil {
		query = query.Where("country = ?", *filter.Country)
	}
	if filter.Age != nil {
		query = query.Where("age = ?", *filter.Age)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(10) // default limit
	}

	query = query.Offset(filter.Offset)

	if err := query.Find(&persons).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list persons: %w", err)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count persons: %w", err)
	}
	totalCounts := totalCount
	return persons, totalCounts, nil
}

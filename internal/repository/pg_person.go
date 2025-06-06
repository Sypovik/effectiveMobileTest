package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sypovik/effectiveMobileTest/internal/models"
	"github.com/rs/zerolog/log"
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
	// Здесь уже не указываем отдельное поле "source":
	// просто пишем “PersonRepository.Create: ...” в тексте сообщения.
	logger := log.Ctx(ctx)

	err := r.db.WithContext(ctx).Create(person).Error
	if err != nil {
		logger.
			Error().
			Err(err).
			Msgf("PersonRepository.Create: Ошибка при создании пользователя — имя=%s фамилия=%s", person.Name, person.Surname)
		return err
	}

	logger.
		Debug().
		Int("ID", person.ID).
		Msg("PersonRepository.Create: Пользователь успешно создан")
	return nil
}

func (r *PgPersonRepository) GetByID(ctx context.Context, id int) (*models.Person, error) {
	logger := log.Ctx(ctx)

	var person models.Person
	result := r.db.WithContext(ctx).First(&person, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.
			Debug().
			Int("ID", id).
			Msg("PersonRepository.GetByID: Пользователь не найден")
		return nil, nil
	}
	if result.Error != nil {
		logger.
			Error().
			Err(result.Error).
			Int("ID", id).
			Msg("PersonRepository.GetByID: Ошибка при получении пользователя")
		return nil, result.Error
	}

	logger.
		Debug().
		Int("ID", id).
		Msg("PersonRepository.GetByID: Пользователь получен")
	return &person, nil
}

func (r *PgPersonRepository) Update(ctx context.Context, person *models.Person) error {
	logger := log.Ctx(ctx)

	err := r.db.WithContext(ctx).Save(person).Error
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("ID", person.ID).
			Msg("PersonRepository.Update: Ошибка при обновлении пользователя")
		return err
	}

	logger.
		Debug().
		Int("ID", person.ID).
		Msg("PersonRepository.Update: Пользователь успешно обновлён")
	return nil
}

func (r *PgPersonRepository) Delete(ctx context.Context, id int) error {
	logger := log.Ctx(ctx)

	err := r.db.WithContext(ctx).Delete(&models.Person{}, id).Error
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("ID", id).
			Msg("PersonRepository.Delete: Ошибка при удалении пользователя")
		return err
	}

	logger.
		Debug().
		Int("ID", id).
		Msg("PersonRepository.Delete: Пользователь успешно удалён")
	return nil
}

func (r *PgPersonRepository) List(ctx context.Context, filter PersonFilter) (results []models.Person, totalCount int64, err error) {
	logger := log.Ctx(ctx)

	logger.
		Debug().
		Int("limit", filter.Limit).
		Int("offset", filter.Offset).
		Msg("PersonRepository.List: Запрос к БД на получение списка")

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
	if filter.Patronymic != nil {
		query = query.Where("patronymic = ?", *filter.Patronymic)
	}
	if filter.Country != nil {
		query = query.Where("country = ?", *filter.Country)
	}
	if filter.Age != nil {
		query = query.Where("age = ?", *filter.Age)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonRepository.List: Ошибка при подсчёте общего количества пользователей")
		return nil, 0, fmt.Errorf("ошибка при подсчёте количества: %w", err)
	}

	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	query = query.Limit(limit).Offset(filter.Offset)

	if err := query.Find(&persons).Error; err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonRepository.List: Ошибка при получении списка пользователей")
		return nil, 0, fmt.Errorf("ошибка при получении списка: %w", err)
	}

	logger.
		Debug().
		Int("count", len(persons)).
		Int64("total", totalCount).
		Any("first person", persons[:1]).
		Msg("PersonRepository.List: Результат получен из БД")

	return persons, totalCount, nil
}

package services

import (
	"context"
	"time"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/models"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
	"github.com/rs/zerolog/log"
)

type PersonService struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

// ===== Create =====
func (s *PersonService) Create(ctx context.Context, req dto.CreatePersonRequest) (*dto.PersonResponse, error) {
	logger := log.Ctx(ctx)
	logger.
		Debug().
		Str("name", req.Name).
		Str("surname", req.Surname).
		Msg("PersonService.Create: Входящие данные")

	// Обогащение данных
	age, gender, country := enrichData(req.Name)

	person := models.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
		Age:        age,
		Gender:     gender,
		Country:    country,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, &person); err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonService.Create: Ошибка при создании персоны в репозитории")
		return nil, err
	}

	logger.
		Debug().
		Int("id", person.ID).
		Msg("PersonService.Create: Персона успешно создана")

	return toDTO(&person), nil
}

// ===== GetByID =====
func (s *PersonService) GetByID(ctx context.Context, id int) (*dto.PersonResponse, error) {
	logger := log.Ctx(ctx)
	logger.
		Debug().
		Int("id", id).
		Msg("PersonService.GetByID: Запрос пользователя по ID")

	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("id", id).
			Msg("PersonService.GetByID: Ошибка получения персоны из репозитория")
		return nil, err
	}

	if person == nil {
		logger.
			Debug().
			Int("id", id).
			Msg("PersonService.GetByID: Персона не найдена")
		return nil, nil
	}

	logger.
		Debug().
		Int("id", id).
		Msg("PersonService.GetByID: Персона найдена")

	return toDTO(person), nil
}

// ===== List =====
func (s *PersonService) List(ctx context.Context, filter repository.PersonFilter) (*dto.ListPersonsResponse, error) {
	logger := log.Ctx(ctx)
	logger.
		Debug().
		Int("limit", filter.Limit).
		Int("offset", filter.Offset).
		Msg("PersonService.List: Запрос списка пользователей")

	persons, total, err := s.repo.List(ctx, filter)
	if err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonService.List: Ошибка получения списка из репозитория")
		return nil, err
	}

	var dtos []dto.PersonResponse
	for _, p := range persons {
		dtos = append(dtos, *toDTO(&p))
	}

	page := filter.Offset/filter.Limit + 1

	logger.
		Debug().
		Int("получено", len(dtos)).
		Int64("общее_количество", total).
		Int("страница", page).
		Int("размер_страницы", filter.Limit).
		Msg("PersonService.List: Список сформирован")

	return &dto.ListPersonsResponse{
		TotalCount: total,
		Page:       page,
		Size:       filter.Limit,
		Data:       dtos,
	}, nil
}

// ===== Update =====
func (s *PersonService) Update(ctx context.Context, id int, req dto.UpdatePersonRequest) (*dto.PersonResponse, error) {
	logger := log.Ctx(ctx)
	logger.
		Debug().
		Int("id", id).
		Msg("PersonService.Update: Входящие данные на обновление")

	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("id", id).
			Msg("PersonService.Update: Ошибка при получении персоны")
		return nil, err
	}

	if req.Name != nil {
		person.Name = *req.Name
	}
	if req.Surname != nil {
		person.Surname = *req.Surname
	}
	if req.Patronymic != nil {
		person.Patronymic = req.Patronymic
	}
	if req.Name != nil {
		person.Age, person.Gender, person.Country = enrichData(person.Name)
	}

	person.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, person); err != nil {
		logger.
			Error().
			Err(err).
			Int("id", id).
			Msg("PersonService.Update: Ошибка при обновлении персоны")
		return nil, err
	}

	logger.
		Debug().
		Int("id", id).
		Msg("PersonService.Update: Персона успешно обновлена")

	return toDTO(person), nil
}

// ===== Delete =====
func (s *PersonService) Delete(ctx context.Context, id int) error {
	logger := log.Ctx(ctx)
	logger.
		Debug().
		Int("id", id).
		Msg("PersonService.Delete: Запрос на удаление персоны")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("id", id).
			Msg("PersonService.Delete: Ошибка при удалении персоны")
		return err
	}

	logger.
		Debug().
		Int("id", id).
		Msg("PersonService.Delete: Персона успешно удалена")

	return nil
}

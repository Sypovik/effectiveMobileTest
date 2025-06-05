package services

import (
	"context"
	"time"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/models"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
)

type PersonService struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

// ===== Create =====

func (s *PersonService) Create(ctx context.Context, req dto.CreatePersonRequest) (*dto.PersonResponse, error) {
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
		return nil, err
	}

	return toDTO(&person), nil
}

// ===== GetByID =====

func (s *PersonService) GetByID(ctx context.Context, id int) (*dto.PersonResponse, error) {
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(person), nil
}

// ===== List =====

func (s *PersonService) List(ctx context.Context, filter repository.PersonFilter) (*dto.ListPersonsResponse, error) {
	persons, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	var dtos []dto.PersonResponse
	for _, p := range persons {
		dtos = append(dtos, *toDTO(&p))
	}

	page := filter.Offset/filter.Limit + 1
	return &dto.ListPersonsResponse{
		TotalCount: total,
		Page:       page,
		Size:       filter.Limit,
		Data:       dtos,
	}, nil
}

// ===== Update =====

func (s *PersonService) Update(ctx context.Context, id int, req dto.UpdatePersonRequest) (*dto.PersonResponse, error) {
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
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

	// Повторное обогащение, если имя изменилось
	if req.Name != nil {
		person.Age, person.Gender, person.Country = enrichData(person.Name)
	}

	person.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, person); err != nil {
		return nil, err
	}

	return toDTO(person), nil
}

// ===== Delete =====

func (s *PersonService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

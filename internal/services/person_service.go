package services

import (
	"context"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
)

// PersonService описывает бизнес‑логику (обогащение + CRUD).
type PersonService interface {
	// Create обогащает и сохраняет нового человека.
	// В качестве входа — dto.CreatePersonRequest, завершается dto.PersonResponse
	Create(ctx context.Context, req dto.CreatePersonRequest) (*dto.PersonResponse, error)

	// GetByID получает одного человека по ID и возвращает dto.PersonResponse
	GetByID(ctx context.Context, id int) (*dto.PersonResponse, error)

	// List возвращает список с фильтрами/пагинацией (dto.ListPersonsResponse)
	List(ctx context.Context, filters repository.PersonFilter) (*dto.ListPersonsResponse, error)

	// Update обновляет запись (если меняется имя — заново обогащает)
	Update(ctx context.Context, id int, req dto.UpdatePersonRequest) (*dto.PersonResponse, error)

	// Delete удаляет запись по ID
	Delete(ctx context.Context, id int) error
}

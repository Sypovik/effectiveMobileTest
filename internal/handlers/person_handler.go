package handlers

import (
	"github.com/Sypovik/effectiveMobileTest/internal/services"
	"github.com/gin-gonic/gin"
)

// PersonHandler хранит зависимость на service (интерфейс).
type PersonHandler struct {
	Service services.PersonService // см. ниже: интерфейс в service
}

// RegisterPersonRoutes регистрирует все маршруты, связанные с /people.
func RegisterPersonRoutes(r *gin.Engine, svc services.PersonService) {
	h := &PersonHandler{Service: svc}

	grp := r.Group("/people")
	{
		grp.POST("", h.CreatePerson)       // POST   /people
		grp.GET("", h.ListPersons)         // GET    /people
		grp.GET("/:id", h.GetPersonByID)   // GET    /people/{id}
		grp.PUT("/:id", h.UpdatePerson)    // PUT    /people/{id}    (или PATCH)
		grp.DELETE("/:id", h.DeletePerson) // DELETE /people/{id}
	}
}

// CreatePerson — хендлер для POST /people
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	// 1. Привести JSON → dto.CreatePersonRequest
	// 2. Вызвать h.Service.Create(ctx, dto)
	// 3. Вернуть JSON dto.PersonResponse или ошибку
}

// ListPersons — хендлер для GET /people?name=&surname=&page=&size=&gender=&country=
func (h *PersonHandler) ListPersons(c *gin.Context) {
	// 1. Считать query-параметры (page, size, фильтры)
	// 2. Вызвать h.Service.List(ctx, filters, pagination)
	// 3. Вернуть JSON dto.ListPersonsResponse
}

// GetPersonByID — хендлер для GET /people/{id}
func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	// 1. Прочитать id из path
	// 2. h.Service.GetByID(ctx, id)
	// 3. Вернуть JSON dto.PersonResponse или 404
}

// UpdatePerson — хендлер для PUT /people/{id}
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	// 1. Прочитать id из path
	// 2. Привести JSON → dto.UpdatePersonRequest
	// 3. h.Service.Update(ctx, id, dto)
	// 4. Вернуть JSON обновлённой dto.PersonResponse
}

// DeletePerson — хендлер для DELETE /people/{id}
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	// 1. Прочитать id из path
	// 2. h.Service.Delete(ctx, id)
	// 3. Вернуть статус 204 No Content
}

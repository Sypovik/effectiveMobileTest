package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
	"github.com/Sypovik/effectiveMobileTest/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PersonHandler struct {
	Service services.PersonService
}

func RegisterPersonRoutes(r *gin.Engine, svc services.PersonService) {
	h := &PersonHandler{Service: svc}
	grp := r.Group("/people")
	{
		grp.POST("", h.CreatePerson)
		grp.GET("", h.ListPersons)
		grp.GET("/:id", h.GetPersonByID)
		grp.PUT("/:id", h.UpdatePerson)
		grp.DELETE("/:id", h.DeletePerson)
	}
}

// CreatePerson godoc
// @Summary Создать новую персону
// @Description Создает новую запись о человеке
// @Tags people
// @Accept json
// @Produce json
// @Param input body dto.CreatePersonRequest true "Данные для создания персоны"
// @Success 201 {object} dto.PersonResponse "Успешно создано"
// @Failure 400 {object} object "Неверный запрос"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /people [post]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)

	logger.
		Info().
		Msg("PersonHandler.CreatePerson: Входящий запрос")

	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonHandler.CreatePerson: Неверный JSON в теле запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	logger.
		Info().
		Str("name", req.Name).
		Str("surname", req.Surname).
		Msg("PersonHandler.CreatePerson: Параметры запроса")

	person, err := h.Service.Create(ctx, req)
	if err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonHandler.CreatePerson: Ошибка при создании персоны в сервисе")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.
		Info().
		Int("ID", person.ID).
		Msg("PersonHandler.CreatePerson: Персона успешно создана")
	c.JSON(http.StatusCreated, person)
}

// ListPersons godoc
// @Summary List persons
// @Description Returns a list of persons with pagination and filtering
// @Tags people
// @Produce json
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param patronymic query string false "Filter by patronymic"
// @Param gender query string false "Filter by gender"
// @Param country query string false "Filter by country"
// @Param age query int false "Filter by age"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} dto.ListPersonsResponse "List of persons"
// @Failure 500 {object} object "Internal server error"
// @Router /people [get]
func (h *PersonHandler) ListPersons(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)

	logger.
		Info().
		Msg("PersonHandler.ListPersons: Входящий запрос")

	var filter repository.PersonFilter

	// Параметры фильтрации
	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}
	if surname := c.Query("surname"); surname != "" {
		filter.Surname = &surname
	}
	if patronymic := c.Query("patronymic"); patronymic != "" {
		filter.Patronymic = &patronymic
	}
	if gender := c.Query("gender"); gender != "" {
		filter.Gender = &gender
	}
	if country := c.Query("country"); country != "" {
		filter.Country = &country
	}
	if ageStr := c.Query("age"); ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil {
			filter.Age = &age
		} else {
			logger.
				Error().
				Err(err).
				Str("ageParam", ageStr).
				Msg("PersonHandler.ListPersons: Неверный параметр age")
		}
	}

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	filter.Limit = size
	filter.Offset = (page - 1) * size

	logger.
		Info().
		Int("page", page).
		Int("size", size).
		Msg("PersonHandler.ListPersons: Фильтры и пагинация")

	list, err := h.Service.List(ctx, filter)
	if err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonHandler.ListPersons: Ошибка при получении списка в сервисе")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Привяжем страницу и размер к ответу, если dto поддерживает такие поля
	list.Page = page
	list.Size = size

	logger.
		Info().
		Int("получено", len(list.Data)).
		Int64("общее_количество", list.TotalCount).
		Int("страница", page).
		Int("размер_страницы", size).
		Msg("PersonHandler.ListPersons: Список сформирован")

	c.JSON(http.StatusOK, list)
}

// GetPersonByID godoc
// @Summary Получить персону по ID
// @Description Возвращает информацию о человеке по его идентификатору
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID персоны"
// @Success 200 {object} dto.PersonResponse "Успешный ответ"
// @Failure 400 {object} object "Неверный ID"
// @Failure 404 {object} object "Персона не найдена"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /people/{id} [get]
func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)

	logger.
		Info().
		Msg("PersonHandler.GetPersonByID: Входящий запрос")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		logger.
			Error().
			Err(err).
			Str("idParam", idStr).
			Msg("PersonHandler.GetPersonByID: Неверный id в параметрах")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	logger.
		Info().
		Int("ID", id).
		Msg("PersonHandler.GetPersonByID: Параметр id корректен, ищем пользователя")

	person, err := h.Service.GetByID(ctx, id)
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("ID", id).
			Msg("PersonHandler.GetPersonByID: Ошибка при получении персоны из сервиса")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if person == nil {
		logger.
			Info().
			Int("ID", id).
			Msg("PersonHandler.GetPersonByID: Пользователь не найден")
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	logger.
		Info().
		Int("ID", id).
		Msg("PersonHandler.GetPersonByID: Пользователь найден")
	c.JSON(http.StatusOK, person)
}

// UpdatePerson godoc
// @Summary Обновить данные персоны
// @Description Обновляет информацию о человеке по его идентификатору
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID персоны"
// @Param input body dto.UpdatePersonRequest true "Данные для обновления"
// @Success 200 {object} dto.PersonResponse "Успешный ответ"
// @Failure 400 {object} object "Неверный запрос"
// @Failure 404 {object} object "Персона не найдена"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /people/{id} [put]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)

	logger.
		Info().
		Msg("PersonHandler.UpdatePerson: Входящий запрос")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		logger.
			Error().
			Err(err).
			Str("idParam", idStr).
			Msg("PersonHandler.UpdatePerson: Неверный id в параметрах")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.
			Error().
			Err(err).
			Msg("PersonHandler.UpdatePerson: Неверный JSON в теле запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	logger.
		Info().
		Int("ID", id).
		Str("Name", *req.Name).
		Str("Surname", *req.Surname).
		Msg("PersonHandler.UpdatePerson: Параметры для обновления")

	person, err := h.Service.Update(ctx, id, req)
	if err != nil {
		logger.
			Error().
			Err(err).
			Int("ID", id).
			Msg("PersonHandler.UpdatePerson: Ошибка при обновлении персоны в сервисе")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.
		Info().
		Int("ID", id).
		Msg("PersonHandler.UpdatePerson: Персона успешно обновлена")
	c.JSON(http.StatusOK, person)
}

// DeletePerson godoc
// @Summary Удалить персону
// @Description Удаляет запись о человеке по его ID
// @Tags people
// @Param id path int true "ID персоны"
// @Success 204 "Успешно удалено"
// @Failure 400 {object} object "Неверный запрос"
// @Failure 404 {object} object "Персона не найдена"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /people/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)

	logger.
		Info().
		Msg("PersonHandler.DeletePerson: Входящий запрос")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		logger.
			Error().
			Err(err).
			Str("idParam", idStr).
			Msg("PersonHandler.DeletePerson: Неверный id в параметрах")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	logger.
		Info().
		Int("ID", id).
		Msg("PersonHandler.DeletePerson: Пытаемся удалить пользователя")

	if err := h.Service.Delete(ctx, id); err != nil {
		logger.
			Error().
			Err(err).
			Int("ID", id).
			Msg("PersonHandler.DeletePerson: Ошибка при удалении персоны в сервисе")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.
		Info().
		Int("ID", id).
		Msg("PersonHandler.DeletePerson: Персона успешно удалена")
	c.Status(http.StatusNoContent)
}

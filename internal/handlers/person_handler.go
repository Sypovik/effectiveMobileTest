package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
	"github.com/Sypovik/effectiveMobileTest/internal/services"
	"github.com/gin-gonic/gin"
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

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	person, err := h.Service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (h *PersonHandler) ListPersons(c *gin.Context) {
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

	list, err := h.Service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list.Page = page
	list.Size = size
	c.JSON(http.StatusOK, list)
}

func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	person, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	person, err := h.Service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

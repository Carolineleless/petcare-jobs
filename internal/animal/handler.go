package animal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service}
}

func (h *Handler) Create(c *gin.Context) {
	var input animal

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	animal := NewAnimal(req)

	h.repo.Save(animal)

	c.JSON(http.StatusCreated, animal)
}

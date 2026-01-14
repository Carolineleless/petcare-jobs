package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Carolineleless/petcare-jobs/internal/animal"
	"github.com/Carolineleless/petcare-jobs/internal/exam"
)

func NewRouter(
	animalHandler *animal.Handler,
	examHandler *exam.Handler,
) *gin.Engine {

	r := gin.Default()

	r.POST("/animals", animalHandler.Create)
	r.POST("/exams", examHandler.Create)

	return r
}

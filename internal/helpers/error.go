package helpers

import (
	"github.com/gin-gonic/gin"
	"library-backend/internal/dto/response"
)

func HandleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, response.ErrorResponseDTO{
		Message: errorMessage,
	})
}

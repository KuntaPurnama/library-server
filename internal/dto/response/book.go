package response

import "library-backend/internal/models"

type GetBooksBySubjectResponseDTO struct {
	Books []models.Book `json:"books"`
}

type SubmitPickupScheduleResponseDTO struct {
	Message string `json:"message"`
}

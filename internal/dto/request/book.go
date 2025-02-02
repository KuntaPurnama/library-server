package request

import "library-backend/internal/models"

type SubmitPickupScheduleRequestDTO struct {
	User       models.UserDetails `json:"user" validate:"required"`
	Book       models.Book        `json:"book" validate:"required"`
	PickupDate string             `json:"pickup_date" validate:"required"`
	PickupTime string             `json:"pickup_time" validate:"required"`
}

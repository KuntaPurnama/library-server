package models

type UserDetails struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

type PickupSchedule struct {
	ID         string      `json:"id"`
	User       UserDetails `json:"user"`
	Book       Book        `json:"book"`
	PickupDate string      `json:"pickup_date"`
	PickupTime string      `json:"pickup_time"`
}

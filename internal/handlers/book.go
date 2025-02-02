package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"library-backend/config"
	"library-backend/internal/dto/request"
	"library-backend/internal/dto/response"
	"library-backend/internal/helpers"
	"library-backend/internal/models"
	"library-backend/internal/service"
	"library-backend/pkg/openlibrary"
	"net/http"
	"strconv"
	"sync"
)

var validate = validator.New()

type BookHandler struct {
	openLibraryClient *openlibrary.Client
	books             []models.Book
	pickupSchedules   map[string]models.PickupSchedule
	mu                sync.RWMutex
}

func NewBookHandler() *BookHandler {
	return &BookHandler{
		openLibraryClient: &openlibrary.Client{},
		books:             []models.Book{},
		pickupSchedules:   make(map[string]models.PickupSchedule),
	}
}

func (h *BookHandler) GetBooksBySubject(c *gin.Context) {
	subject := c.Query("subject")
	if subject == "" {
		helpers.HandleError(c, http.StatusBadRequest, "Subject is required")
		return
	}

	limit := config.Config.OpenLibLimit
	page := 0

	if o := c.Query("page"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			page = parsed
		}
	}

	books, err := h.openLibraryClient.FetchBooksBySubject(subject, limit, page)
	if err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, fmt.Sprintf("Error fetching books: %v", err))
		return
	}

	h.mu.Lock()
	h.books = books
	h.mu.Unlock()

	responseData := response.GetBooksBySubjectResponseDTO{
		Books: books,
	}

	c.JSON(http.StatusOK, responseData)
}

func (h *BookHandler) SubmitPickupSchedule(c *gin.Context) {
	var scheduleInput request.SubmitPickupScheduleRequestDTO
	if err := c.ShouldBindJSON(&scheduleInput); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	err := validate.Struct(scheduleInput)
	if err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is %s", e.Field(), e.Tag()))
		}

		helpers.HandleError(c, http.StatusBadRequest, fmt.Sprintf("Invalid input format: %v", errorMessages))
		return
	}

	id := uuid.New().String()
	newSchedule := models.PickupSchedule{
		ID:         id,
		User:       scheduleInput.User,
		Book:       scheduleInput.Book,
		PickupDate: scheduleInput.PickupDate,
		PickupTime: scheduleInput.PickupTime,
	}

	h.mu.Lock()
	h.pickupSchedules[id] = newSchedule
	h.mu.Unlock()

	go func() {
		service.SendEmailConfirmation(newSchedule)
	}()

	responseData := response.SubmitPickupScheduleResponseDTO{
		Message: "Pickup schedule created successfully",
	}

	c.JSON(http.StatusCreated, responseData)
}

func (h *BookHandler) ListPickupSchedules(c *gin.Context) {
	var schedules []models.PickupSchedule
	for _, schedule := range h.pickupSchedules {
		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules)
}

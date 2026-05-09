package handlers

import (
	"net/http"
	"ocs-room-booking/services"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	svc *services.BookingService
}

func NewBookingHandler(svc *services.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

type CreateBookingRequest struct {
	RoomID           string `json:"room_id" binding:"required"`
	Date             string `json:"date" binding:"required"`
	StartTime        string `json:"start_time" binding:"required"`
	EndTime          string `json:"end_time" binding:"required"`
	Purpose          string `json:"purpose" binding:"required"`
	ParticipantCount int    `json:"participant_count" binding:"required,min=1"`
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid user context")
		return
	}

	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room_id")
		return
	}

	booking, err := h.svc.CreateBooking(userID, roomID, req.Date, req.StartTime, req.EndTime, req.Purpose, req.ParticipantCount)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "room is already booked for the selected time slot" {
			status = http.StatusConflict
		}
		utils.Error(c, status, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, booking)
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid user context")
		return
	}

	bookings, err := h.svc.GetMyBookings(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	utils.Success(c, http.StatusOK, bookings)
}

func (h *BookingHandler) CancelMyBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.svc.CancelBooking(id, userID, false); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "Booking cancelled"})
}

func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	bookings, err := h.svc.GetAllBookings()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	utils.Success(c, http.StatusOK, bookings)
}

func (h *BookingHandler) AdminCancelBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	if err := h.svc.CancelBooking(id, uuid.Nil, true); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "Booking cancelled"})
}

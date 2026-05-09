package services

import (
	"errors"
	"ocs-room-booking/models"
	"ocs-room-booking/repository"
	"time"

	"github.com/google/uuid"
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
	roomRepo    *repository.RoomRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository, roomRepo *repository.RoomRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo, roomRepo: roomRepo}
}

func (s *BookingService) CreateBooking(userID uuid.UUID, roomID uuid.UUID, date, startTime, endTime, purpose string, participantCount int) (*models.Booking, error) {

	// 1. Validate purpose
	validPurposes := map[string]bool{"OA": true, "Interview": true, "PPT": true}
	if !validPurposes[purpose] {
		return nil, errors.New("purpose must be one of: OA, Interview, PPT")
	}

	// 2. Validate time ordering
	startT, err := time.Parse("15:04", startTime)
	if err != nil {
		startT, err = time.Parse("15:04:05", startTime)
		if err != nil {
			return nil, errors.New("invalid start_time format, use HH:MM")
		}
	}
	endT, err := time.Parse("15:04", endTime)
	if err != nil {
		endT, err = time.Parse("15:04:05", endTime)
		if err != nil {
			return nil, errors.New("invalid end_time format, use HH:MM")
		}
	}
	if !endT.After(startT) {
		return nil, errors.New("end_time must be strictly after start_time")
	}

	// 3. Date not in the past
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}
	today := time.Now().Truncate(24 * time.Hour)
	if bookingDate.Before(today) {
		return nil, errors.New("booking date cannot be in the past")
	}

	// 4. Fetch room
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// 5. Room availability flag
	if !room.IsAvailable {
		return nil, errors.New("room is not available for booking")
	}

	// 6. Capacity check
	if participantCount > room.Capacity {
		return nil, errors.New("participant count exceeds room capacity")
	}

	// 7. Purpose allowed for this room
	purposeAllowed := false
	for _, p := range room.AllowedPurposes {
		if p == purpose {
			purposeAllowed = true
			break
		}
	}
	if !purposeAllowed {
		return nil, errors.New("this purpose is not allowed for the selected room")
	}

	// 8. Conflict check inside a transaction with SELECT FOR UPDATE
	db := s.bookingRepo.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	conflict, err := s.bookingRepo.CheckConflict(tx, roomID, date, startTime, endTime)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to check for conflicts")
	}
	if conflict {
		tx.Rollback()
		return nil, errors.New("room is already booked for the selected time slot")
	}

	booking := &models.Booking{
		RoomID:           roomID,
		UserID:           userID,
		Date:             date,
		StartTime:        startTime,
		EndTime:          endTime,
		Purpose:          purpose,
		ParticipantCount: participantCount,
		Status:           "confirmed",
	}
	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create booking")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit booking")
	}

	return booking, nil
}

func (s *BookingService) GetMyBookings(userID uuid.UUID) ([]models.Booking, error) {
	return s.bookingRepo.FindByUserID(userID)
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	return s.bookingRepo.FindAll()
}

func (s *BookingService) CancelBooking(bookingID uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	if !isAdmin && booking.UserID != userID {
		return errors.New("you can only cancel your own bookings")
	}
	if booking.Status == "cancelled" {
		return errors.New("booking is already cancelled")
	}
	booking.Status = "cancelled"
	return s.bookingRepo.Update(booking)
}

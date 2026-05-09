package repository

import (
	"ocs-room-booking/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Room.Block").Preload("User").First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindByUserID(userID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Room.Block").Where("user_id = ?", userID).Order("created_at desc").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindAll() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Room.Block").Preload("User").Order("created_at desc").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

// CheckConflict uses SELECT FOR UPDATE to safely detect overlaps
func (r *BookingRepository) CheckConflict(tx *gorm.DB, roomID uuid.UUID, date, startTime, endTime string) (bool, error) {
	var count int64
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Model(&models.Booking{}).
		Where(
			"room_id = ? AND date = ? AND status = 'confirmed' AND start_time < ? AND end_time > ?",
			roomID, date, endTime, startTime,
		).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BookingRepository) GetDB() *gorm.DB {
	return r.db
}

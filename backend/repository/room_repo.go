package repository

import (
	"ocs-room-booking/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) FindAll() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Preload("Block").Find(&rooms).Error
	return rooms, err
}

func (r *RoomRepository) FindByID(id uuid.UUID) (*models.Room, error) {
	var room models.Room
	err := r.db.Preload("Block").First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *RoomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *RoomRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Room{}, "id = ?", id).Error
}

func (r *RoomRepository) Search(blockID, minCapacity, purpose, date, startTime, endTime string) ([]models.Room, error) {
	query := r.db.Preload("Block").Where("is_available = true")

	if blockID != "" {
		query = query.Where("block_id = ?", blockID)
	}
	if minCapacity != "" {
		query = query.Where("capacity >= ?", minCapacity)
	}
	if purpose != "" {
		query = query.Where("? = ANY(allowed_purposes)", purpose)
	}

	var rooms []models.Room
	if err := query.Find(&rooms).Error; err != nil {
		return nil, err
	}

	// Filter out rooms with conflicting bookings
	if date != "" && startTime != "" && endTime != "" {
		var available []models.Room
		for _, room := range rooms {
			var count int64
			r.db.Model(&models.Booking{}).
				Where("room_id = ? AND date = ? AND status = 'confirmed' AND start_time < ? AND end_time > ?",
					room.ID, date, endTime, startTime).
				Count(&count)
			if count == 0 {
				available = append(available, room)
			}
		}
		return available, nil
	}

	return rooms, nil
}

func (r *RoomRepository) FindAllBlocks() ([]models.Block, error) {
	var blocks []models.Block
	err := r.db.Find(&blocks).Error
	return blocks, err
}

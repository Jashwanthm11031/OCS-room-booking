package services

import (
	"errors"
	"ocs-room-booking/models"
	"ocs-room-booking/repository"

	"github.com/google/uuid"
)

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) SearchRooms(blockID, minCapacity, purpose, date, startTime, endTime string) ([]models.Room, error) {
	if purpose != "" {
		valid := map[string]bool{"OA": true, "Interview": true, "PPT": true}
		if !valid[purpose] {
			return nil, errors.New("purpose must be OA, Interview, or PPT")
		}
	}
	return s.repo.Search(blockID, minCapacity, purpose, date, startTime, endTime)
}

func (s *RoomService) GetRoomByID(id uuid.UUID) (*models.Room, error) {
	return s.repo.FindByID(id)
}

func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	return s.repo.FindAll()
}

func (s *RoomService) CreateRoom(blockID uuid.UUID, roomName string, capacity int, allowedPurposes []string, notes string) (*models.Room, error) {
	if roomName == "" || capacity <= 0 {
		return nil, errors.New("room name and capacity are required")
	}
	if allowedPurposes == nil {
		allowedPurposes = []string{"OA", "Interview", "PPT"}
	}
	room := &models.Room{
		BlockID:         blockID,
		RoomName:        roomName,
		Capacity:        capacity,
		IsAvailable:     true,
		AllowedPurposes: allowedPurposes,
		Notes:           notes,
	}
	if err := s.repo.Create(room); err != nil {
		return nil, errors.New("failed to create room")
	}
	return room, nil
}

func (s *RoomService) UpdateRoom(id uuid.UUID, roomName string, capacity int, isAvailable *bool, allowedPurposes []string, notes string) (*models.Room, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if roomName != "" {
		room.RoomName = roomName
	}
	if capacity > 0 {
		room.Capacity = capacity
	}
	if isAvailable != nil {
		room.IsAvailable = *isAvailable
	}
	if allowedPurposes != nil {
		room.AllowedPurposes = allowedPurposes
	}
	if notes != "" {
		room.Notes = notes
	}
	if err := s.repo.Update(room); err != nil {
		return nil, errors.New("failed to update room")
	}
	return room, nil
}

func (s *RoomService) DeleteRoom(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("room not found")
	}
	return s.repo.Delete(id)
}

func (s *RoomService) GetAllBlocks() ([]models.Block, error) {
	return s.repo.FindAllBlocks()
}

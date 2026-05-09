package handlers

import (
	"net/http"
	"ocs-room-booking/services"
	"ocs-room-booking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomHandler struct {
	svc *services.RoomService
}

func NewRoomHandler(svc *services.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) SearchRooms(c *gin.Context) {
	block := c.Query("block")
	capacity := c.Query("capacity")
	purpose := c.Query("purpose")
	date := c.Query("date")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	rooms, err := h.svc.SearchRooms(block, capacity, purpose, date, startTime, endTime)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room ID")
		return
	}
	room, err := h.svc.GetRoomByID(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Room not found")
		return
	}
	utils.Success(c, http.StatusOK, room)
}

func (h *RoomHandler) GetAllBlocks(c *gin.Context) {
	blocks, err := h.svc.GetAllBlocks()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch blocks")
		return
	}
	utils.Success(c, http.StatusOK, blocks)
}

type CreateRoomRequest struct {
	BlockID         string   `json:"block_id" binding:"required"`
	RoomName        string   `json:"room_name" binding:"required"`
	Capacity        int      `json:"capacity" binding:"required,min=1"`
	AllowedPurposes []string `json:"allowed_purposes"`
	Notes           string   `json:"notes"`
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	blockID, err := uuid.Parse(req.BlockID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid block_id")
		return
	}
	room, err := h.svc.CreateRoom(blockID, req.RoomName, req.Capacity, req.AllowedPurposes, req.Notes)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, room)
}

type UpdateRoomRequest struct {
	RoomName        string   `json:"room_name"`
	Capacity        string   `json:"capacity"`
	IsAvailable     *bool    `json:"is_available"`
	AllowedPurposes []string `json:"allowed_purposes"`
	Notes           string   `json:"notes"`
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room ID")
		return
	}
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}
	cap := 0
	if req.Capacity != "" {
		cap, _ = strconv.Atoi(req.Capacity)
	}
	room, err := h.svc.UpdateRoom(id, req.RoomName, cap, req.IsAvailable, req.AllowedPurposes, req.Notes)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room ID")
		return
	}
	if err := h.svc.DeleteRoom(id); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

// ListAllRooms returns ALL rooms including unavailable ones (admin use)
func (h *RoomHandler) ListAllRooms(c *gin.Context) {
	rooms, err := h.svc.GetAllRooms()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch rooms")
		return
	}
	utils.Success(c, http.StatusOK, rooms)
}

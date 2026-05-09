package handlers

import (
	"net/http"
	"ocs-room-booking/services"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	utils.Success(c, http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	user, err := h.svc.CreateUser(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, user)
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}
	user, err := h.svc.UpdateUser(id, req.Name, req.Role, req.IsActive)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, user)
}

package handlers

import (
	"net/http"
	"ocs-room-booking/repository"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: email and password are required")
		return
	}

	// Step 1: Find user by email
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Step 2: Check is_active
	if !user.IsActive {
		utils.Error(c, http.StatusUnauthorized, "Account has been deactivated")
		return
	}

	// Step 3: Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Step 4: Issue JWT
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

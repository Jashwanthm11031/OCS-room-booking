package services

import (
	"errors"
	"ocs-room-booking/models"
	"ocs-room-booking/repository"
	"ocs-room-booking/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(name, email, password, role string) (*models.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	validRoles := map[string]bool{"admin": true, "core": true, "viewer": true}
	if !validRoles[role] {
		return nil, errors.New("role must be one of: admin, core, viewer")
	}
	if s.repo.EmailExists(email) {
		return nil, errors.New("email already in use")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}
	return user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, name, role string, isActive *bool) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if name != "" {
		user.Name = name
	}
	if role != "" {
		validRoles := map[string]bool{"admin": true, "core": true, "viewer": true}
		if !validRoles[role] {
			return nil, errors.New("invalid role")
		}
		user.Role = role
	}
	if isActive != nil {
		user.IsActive = *isActive
	}
	if err := s.repo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}
	return user, nil
}

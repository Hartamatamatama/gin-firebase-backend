package services

import (
	"github.com/hartamatamatama/gin-firebase-backend/repositories"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}
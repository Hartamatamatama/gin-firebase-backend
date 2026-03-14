package repositories

import (
	"github.com/hartamatamatama/gin-firebase-backend/config"
	"github.com/hartamatamatama/gin-firebase-backend/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByFirebaseUID mencari user berdasarkan Firebase UID
func (r *UserRepository) FindByFirebaseUID(uid string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("firebase_uid = ?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}


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

// VerifyFirebaseToken verifikasi token dari Firebase,
// pastikan email sudah verified, lalu return Backend JWT
func (s *AuthService) VerifyFirebaseToken(firebaseToken string) (string, *models.User, error) {
	// 1. Verifikasi Firebase ID Token ke server Google
	token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), firebaseToken)
	if err != nil {
		return "", nil, errors.New("firebase token tidak valid atau kadaluarsa")
	}

	
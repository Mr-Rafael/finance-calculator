package service

import (
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func TestValidateAccessToken(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: nil,
	}
	mockUsersRepo := &MockUsersRepo{
		CreateUserFunc: nil,
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	want := "001"
	claims := auth.AccessClaims{
		UserID: want,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "savings-app",
			Subject:   want,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mockAccessSecret))
	if err != nil {
		log.Fatalf("Failed to generate test access token: %v", err)
	}

	got, err := service.ValidateAccessToken(signedToken)
	if err != nil {
		log.Fatalf("Failed to validate the access token: %v", err)
	}

	if want != got {
		log.Fatalf("The expected User ID (%v) did not match the one obtained from validation (%v)", want, got)
	}
}

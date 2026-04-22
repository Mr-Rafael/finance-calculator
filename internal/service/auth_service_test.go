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

func TestValidateExpiredToken(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: nil,
	}
	mockUsersRepo := &MockUsersRepo{
		CreateUserFunc: nil,
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	userID := "001"
	claims := auth.AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "savings-app",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mockAccessSecret))
	if err != nil {
		log.Fatalf("Failed to generate test access token: %v", err)
	}

	_, err = service.ValidateAccessToken(signedToken)
	if err == nil {
		log.Fatalf("Validate function returned the expired token as valid.")
	}
}

func TestValidateInvalidSignature(t *testing.T) {
	mockAccessSecret := "ACCESS"
	badAccessSecret := "NOACCESS"
	mockRefreshSecret := "REFRESH"
	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: nil,
	}
	mockUsersRepo := &MockUsersRepo{
		CreateUserFunc: nil,
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	userID := "001"
	claims := auth.AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "savings-app",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(badAccessSecret))
	if err != nil {
		log.Fatalf("Failed to generate test access token: %v", err)
	}

	_, err = service.ValidateAccessToken(signedToken)
	if err == nil {
		log.Fatalf("Validate function wrongly signed token as valid.")
	}
}

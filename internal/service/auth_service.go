package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/auth"
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo      *repository.AuthRepo
	usersRepo     *repository.UsersRepo
	accessSecret  string
	refreshSecret string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginInfo struct {
	ID           pgtype.UUID
	Email        string
	UserName     string
	AccessToken  string
	RefreshToken string
}

func NewAuthService(authRepo *repository.AuthRepo, usersRepo *repository.UsersRepo, accessSecret string, refreshSecret string) AuthService {
	return AuthService{
		authRepo:      authRepo,
		usersRepo:     usersRepo,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (LoginInfo, error) {
	userInfo, err := s.usersRepo.GetUserByEmail(context.Background(), input.Email)

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.PasswordHash), []byte(input.Password))
	if err != nil {
		return LoginInfo{}, fmt.Errorf("password hash mismatch")
	}

	accessToken, err := auth.GenerateAccessToken(userInfo.ID.String(), s.accessSecret)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("error generating access token: %v", err)
	}

	refreshToken, expDate, err := auth.GenerateRefreshToken(userInfo.ID.String(), s.refreshSecret)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("error generating refresh token: %v", err)
	}
	refTokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))

	createParams := ToRefreshTokenCreateParams(userInfo.ID, refTokenHash, expDate)
	_, err = s.authRepo.CreateRefreshToken(context.Background(), createParams)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("error storing the refresh token: %v", err)
	}

	return LoginInfo{
		ID:           userInfo.ID,
		Email:        userInfo.Email,
		UserName:     userInfo.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ToRefreshTokenCreateParams(user pgtype.UUID, tokenHash string, expDate time.Time) db.CreateRefreshTokenParams {
	return db.CreateRefreshTokenParams{
		UserID:    user,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expDate,
			Valid: true,
		},
		Revoked: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}
}

func ToLoginInfoModel(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
	}
}

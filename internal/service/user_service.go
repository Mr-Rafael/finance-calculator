package service

import (
	"context"
	"fmt"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UsersRepo
}

type RegisterUserInput struct {
	Email    string
	Password string
	Username string
}

type User struct {
	ID        pgtype.UUID
	Email     string
	Username  string
	CreatedAt pgtype.Timestamp
}

func (s *UserService) NewUserService(repo *repository.UsersRepo) {
	s.repo = repo
}

func (s *UserService) RegisterUser(ctx context.Context, input RegisterUserInput) (User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to generate password hash: %v", err)
	}

	params := db.CreateUserParams{
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		Username:     input.Username,
	}

	user, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return ToUserModel(user), nil
}

func ToUserModel(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
	}
}

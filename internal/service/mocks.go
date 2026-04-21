package service

import (
	"context"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type MockAuthRepo struct {
	CreateRefreshTokenFunc  func(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error)
	GetTokenByHashFunc      func(ctx context.Context, hash string) (db.RefreshToken, error)
	RevokeTokenByUserIDFunc func(ctx context.Context, id pgtype.UUID) error
}

type MockUsersRepo struct {
	CreateUserFunc     func(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (db.User, error)
	GetUserByIDFunc    func(ctx context.Context, id pgtype.UUID) (db.User, error)
	DeleteUserFunc     func(ctx context.Context, id pgtype.UUID) error
}

func (m *MockAuthRepo) CreateRefreshToken(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
	if m.CreateRefreshTokenFunc != nil {
		return m.CreateRefreshTokenFunc(ctx, userID, tokenHash, expDate)
	}
	return db.RefreshToken{}, nil
}

func (m *MockAuthRepo) GetTokenByHash(ctx context.Context, hash string) (db.RefreshToken, error) {
	if m.CreateRefreshTokenFunc != nil {
		return m.GetTokenByHashFunc(ctx, hash)
	}
	return db.RefreshToken{}, nil
}

func (m *MockAuthRepo) RevokeTokenByUserID(ctx context.Context, id pgtype.UUID) error {
	if m.CreateRefreshTokenFunc != nil {
		return m.RevokeTokenByUserIDFunc(ctx, id)
	}
	return nil
}

func (m *MockUsersRepo) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, params)
	}
	return db.User{}, nil
}

func (m *MockUsersRepo) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.CreateUserFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return db.User{}, nil
}

func (m *MockUsersRepo) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	if m.CreateUserFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return db.User{}, nil
}

func (m *MockUsersRepo) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

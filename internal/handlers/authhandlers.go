package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/auth"
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *ApiConfig) HandlerUsersLogin(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := models.UserLoginRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("received bad user creation request: %v", err), http.StatusBadRequest)
		return
	}

	user, err := cfg.Queries.GetUserByEmail(context.Background(), reqParams.Email)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("failed to fetch user information: %v", err), http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(reqParams.Password))
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("Login attempt with incorrect credentials for user '%v'", reqParams.Email), http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String(), cfg.AccessSecret)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error generating access token: %v", err), "There was an error generating access token.", http.StatusInternalServerError)
	}

	refreshToken, expDate, err := auth.GenerateRefreshToken(user.ID.String(), cfg.RefreshSecret)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error generating refresh token: %v", err), "There was an error generating the refresh token.", http.StatusInternalServerError)
	}

	hashedRefToken := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))
	refreshTokenInsertParams := db.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: hashedRefToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expDate,
			Valid: true,
		},
		Revoked: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}
	_, err = cfg.Queries.CreateRefreshToken(context.Background(), refreshTokenInsertParams)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error storing refresh token on database: %v", err), "There was an error storing refresh toking on DB.", http.StatusInternalServerError)
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	respondWithJSON(writer, models.UserLoginResponseParams{
		ID:          user.ID.String(),
		Email:       user.Email,
		Username:    user.Username,
		AccessToken: accessToken,
	}, http.StatusOK)
}

func (cfg *ApiConfig) HandlerRefresh(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("refresh_token")
	if err != nil {
		respondWithError(writer, "request is missing refresh token", "missing refresh token", http.StatusUnauthorized)
		return
	}

	refreshTokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(cookie.Value)))

	tokenData, err := cfg.Queries.GetTokenByHash(context.Background(), refreshTokenHash)

	if tokenData.Revoked.Bool {
		respondWithErrorCode(writer, fmt.Sprintf("attempt to refresh with revoked token for user: %v", tokenData.UserID), http.StatusUnauthorized)
		return
	}

	if tokenData.ExpiresAt.Time.Before(time.Now()) {
		respondWithErrorCode(writer, fmt.Sprintf("attempt to refresh with expired token for user: %v", tokenData.UserID), http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateAccessToken(tokenData.ID.String(), cfg.AccessSecret)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error generating access token: %v", err), "There was an error generating access token.", http.StatusInternalServerError)
	}

	respondWithJSON(writer, models.RefreshResponseParams{
		AccessToken: accessToken,
	}, http.StatusOK)
}

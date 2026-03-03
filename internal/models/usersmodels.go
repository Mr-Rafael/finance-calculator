package models

type UserCreateRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserCreateResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserLoginRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponseParams struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

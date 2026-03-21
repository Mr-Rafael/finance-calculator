package mapper

import (
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

func ToLoginInput(reqParams dto.UserLoginRequestParams) service.LoginInput {
	return service.LoginInput{
		Email:    reqParams.Email,
		Password: reqParams.Password,
	}
}

func ToLoginResponse(loginInfo service.LoginInfo) dto.UserLoginResponseParams {
	return dto.UserLoginResponseParams{
		ID:          loginInfo.ID.String(),
		Email:       loginInfo.Email,
		Username:    loginInfo.UserName,
		AccessToken: loginInfo.AccessToken,
	}
}

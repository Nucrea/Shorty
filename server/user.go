package server

import (
	"shorty/src/services/users"

	"github.com/gin-gonic/gin"
)

type UserRegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterOutput struct {
	Id string `json:"id"`
}

func (s *server) RegisterUser(ctx *gin.Context, input UserRegisterInput) (*UserRegisterOutput, error) {
	user, err := s.UserService.Create(ctx, users.CreateUserParams{
		Email:    input.Email,
		Password: input.Password,
	})
	if err == users.ErrBadInputEmail || err == users.ErrBadInputPassword || err == users.ErrUserExists {
		return nil, &ErrorBadRequest{err.Error()}
	}
	if err != nil {
		return nil, err
	}

	return &UserRegisterOutput{Id: user.Id}, nil
}

type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	Token string `json:"token"`
}

func (s *server) LoginUser(ctx *gin.Context, input UserLoginInput) (*UserLoginOutput, error) {
	result, err := s.UserService.Login(ctx, input.Email, input.Password)
	if err == users.ErrUserNotExists || err == users.ErrWrongPassword {
		return nil, &ErrorBadRequest{err.Error()}
	}
	if err != nil {
		return nil, err
	}

	return &UserLoginOutput{result.Token}, nil
}

type UserLogoutInput struct {
	Token string `json:"token"`
}

func (s *server) LogoutUser(ctx *gin.Context, input UserLogoutInput) (interface{}, error) {
	err := s.UserService.DeleteSession(ctx, input.Token)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

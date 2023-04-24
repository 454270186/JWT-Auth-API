package service

import "jwtAuth/dto"

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, error)
	Verify(urlParams map[string]string) (bool, error)
}

type DefaultAuthService struct {
	
}
package service

import (
	"jwtAuth/domain"
	"jwtAuth/dto"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
	//Verify(urlParams map[string]string) (bool, error)
}

type DefaultAuthService struct {
	repo domain.AuthRepo
}

func NewAuthService(repository domain.AuthRepo) AuthService {
	return DefaultAuthService{
		repo: repository,
	}
}

func (d DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
	login, err := d.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err;
	}

	token, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}

	return token, nil
}
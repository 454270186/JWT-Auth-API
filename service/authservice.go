package service

import (
	"errors"
	"jwtAuth/domain"
	"jwtAuth/dto"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, error)
	Verify(urlParams map[string]string) (bool, error)
	Refresh(refreshReq dto.RefreshRequest) (*dto.LoginResponse, error)
}

type DefaultAuthService struct {
	repo domain.AuthRepo
	rolePermissions domain.RolePermissions
}

func NewAuthService(repository domain.AuthRepo) AuthService {
	return DefaultAuthService{
		repo: repository,
		rolePermissions: domain.GetRolePermissions(),
	}
}

func (d DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	login, err := d.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	userToken, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}

	err = d.repo.StoreRefreshToken(*userToken)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}, nil
}

func (d DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {
	log.Println(urlParams)
	jwtToken, err := getJwtFromSting(urlParams["token"])
	if err != nil {
		return false, err
	}

	if jwtToken.Valid {
		mapClaims := jwtToken.Claims.(jwt.MapClaims)
		
		claims, err := domain.BuildClaims(mapClaims)
		if err != nil {
			return false, err
		}
		log.Println(claims)
		if claims.IsUser() {
			if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
				log.Println("Failed IsRequestVerifiedWithTokenClaims")
				return false, nil
			}
		}

		return true, nil
	} else {
		return false, errors.New("invalid token")
	}
}

// 1. access token is valid and has expired
// 2. refresh token should exist
// 3. refresh token should be valid and not be expired
func (d DefaultAuthService) Refresh(refreshReq dto.RefreshRequest) (*dto.LoginResponse, error) {
	if jwtErr := refreshReq.IsAccessTokenValid(); jwtErr != nil {
		if errors.Is(jwtErr, jwt.ErrTokenExpired) {
			// jwt refresh function
			if err := d.repo.IsRefreshTokenExist(refreshReq); err != nil {
				return nil, errors.New("refresh token does not exist")
			}
			// generate a new access token from refresh token
			username, err := refreshReq.GetUserName()
			if err != nil {
				return nil, err
			}

			login, err := d.repo.FindByName(username)
			if err != nil {
				return nil, err
			}

			userToken, err := login.GenerateToken()
			if err != nil {
				return nil, err
			}
		
			return &dto.LoginResponse{
				AccessToken: userToken.Token,
			}, nil
		}
		
		return nil, jwtErr
	}

	return nil, errors.New("access is valid")
}

func getJwtFromSting(tokenstring string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}

	return token, nil
}

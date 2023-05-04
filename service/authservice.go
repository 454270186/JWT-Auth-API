package service

import (
	"errors"
	"jwtAuth/domain"
	"jwtAuth/dto"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(dto.LoginRequest) (*domain.AuthToken, error)
	Verify(urlParams map[string]string) (bool, error)
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

func (d DefaultAuthService) Login(req dto.LoginRequest) (*domain.AuthToken, error) {
	login, err := d.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	userToken, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}

	return userToken, nil
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

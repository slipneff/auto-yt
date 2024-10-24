package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
)

type accessClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}
type refreshClaims struct {
	jwt.StandardClaims
	TokenId string `json:"token_id"`
	UserId  string `json:"user_id"`
}
type organisationClaims struct {
	jwt.StandardClaims
	OrgId       string   `json:"org_id"`
	TokenId     string   `json:"token_id"`
	Permissions []string `json:"permissions"`
}

func (s *Service) ParseAccessJWT(ctx context.Context, value string) (string, error) {
	token, err := jwt.ParseWithClaims(value, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedMethod
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", ErrInvalidToken
	}
	myClaims := token.Claims.(*accessClaims)
	return myClaims.UserId, nil
}

func (s *Service) CreateAccessJWT(ctx context.Context, id string) (string, error) {
	jwtSecret := []byte(s.cfg.JWTSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   id,
		"ExpiresAt": time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

type ParseOrganisationJWTOutput struct {
	OrgId       string
	Permissions []string
}

func (s *Service) ParseOrganisationJWT(ctx context.Context, value string) (*ParseOrganisationJWTOutput, error) {
	token, err := jwt.ParseWithClaims(value, &organisationClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedMethod
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	myClaims := token.Claims.(*organisationClaims)
	// fmt.Println("parseOrganisJWT - ", myClaims.Permissions)
	return &ParseOrganisationJWTOutput{
		OrgId:       myClaims.OrgId,
		Permissions: myClaims.Permissions,
	}, nil
}

type CreateOrganisationJWTInput struct {
	// TokenId     string
	OrgId       string
	Permissions []string
}

func (s *Service) CreateOrganisationJWT(ctx context.Context, input CreateOrganisationJWTInput) (string, error) {
	jwtSecret := []byte(s.cfg.JWTSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"org_id": input.OrgId,
		// "token_id":    input.TokenId,
		"permissions": input.Permissions,
		// "ExpiresAt":   time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
func (s *Service) ParseRefreshJWT(ctx context.Context, value string) (string, string, error) {
	token, err := jwt.ParseWithClaims(value, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedMethod
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", "", ErrInvalidToken
	}
	myClaims := token.Claims.(*refreshClaims)
	return myClaims.TokenId, myClaims.UserId, nil
}

func (s *Service) CreateRefreshJWT(ctx context.Context, tokenId string, userId string) (string, error) {
	jwtSecret := []byte(s.cfg.JWTSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_id":  tokenId,
		"user_id":   userId,
		"ExpiresAt": 31536000,
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

type ResetPasswordClaims struct {
	jwt.StandardClaims
	Email  string `json:"Email"`
	UserId string `json:"UserId"`
}

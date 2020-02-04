package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type (
	Jwt struct {
		svcAccount *Account
	}

	Claims struct {
		UserID string
		jwt.StandardClaims
	}
)

func NewJwt(account *Account) *Jwt {
	return &Jwt{account}
}

func (s *Jwt) GenToken(userID string, salt []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{UserID: userID}).SignedString(salt)
}

func (s *Jwt) ParseClaims(token string) (*Claims, error) {
	// Parsing only header and payload parts of token
	// and getting ID from last
	claims := new(Claims)
	_, _, err := new(jwt.Parser).ParseUnverified(token, claims)
	if err != nil {
		return nil, fmt.Errorf("parse token failed: %v", err)
	}
	return claims, nil
}

func (s *Jwt) ParseToken(token string, salt []byte) error {
	parseFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse token failed: unexpected signing method %v", token.Header["alg"])
		}
		return salt, nil
	}
	_, err := jwt.Parse(token, parseFunc)
	return err
}

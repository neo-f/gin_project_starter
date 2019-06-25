package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type JWTClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func Create(email string) (string, error) {
	claims := JWTClaims{
		Email: email,
	}
	expireTime := viper.GetDuration("JWT.EXPIRE_TIME")
	claims.ExpiresAt = time.Now().Add(expireTime).Unix()
	claims.IssuedAt = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("JWT.SECRET")))
}

func Refresh(t string) (string, error) {
	claim, err := parse(t)
	if err != nil {
		return "", err
	}
	expireTime := viper.GetDuration("JWT.EXPIRE_TIME")
	claim.ExpiresAt = time.Now().Add(expireTime).Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, _ := newToken.SignedString([]byte(viper.GetString("JWT.SECRET")))
	return signedToken, nil
}

func parse(t string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(t, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT.SECRET")), nil
	})
	if err != nil {
		return &JWTClaims{}, err
	}
	claims := token.Claims.(*JWTClaims)
	err = token.Claims.Valid()
	return claims, err
}

func Verify(t string) bool {
	if _, err := parse(t); err != nil {
		return false
	}
	return true
}

func GetAccountEmail(t string) (string, error) {
	claim, err := parse(t)
	if err != nil {
		return "", err
	}
	return claim.Email, nil
}

package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
)

var secretKey = viper.GetString("JWT_SECRET")

type MyClaims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	Plan   string `json:"plan"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uuid.UUID, role, plan string) (string, error) {
	iMyClaims := MyClaims{
		UserId: userId.String(),
		Role:   role,
		Plan:   plan,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(GetExpiresIn())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "Token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iMyClaims)

	return token.SignedString([]byte(viper.GetString(secretKey)))
}

func ParseToken(tokenString string) (MyClaims, error) {
	iMyClaims := MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &iMyClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString(secretKey)), nil
	})

	if err == nil && !token.Valid {
		//如果解析沒錯  但是token為非法的 也創一個錯誤給他
		err = errors.New("invalid token")
	}
	return iMyClaims, err
}

func GetExpiresIn() time.Duration {
	return viper.GetDuration("jwt.tokenExpired_min") * time.Minute
}

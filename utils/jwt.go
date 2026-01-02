package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

var secretKey = viper.GetString("JWT_SECRET")

type MyClaims struct {
	Uid  uint   `json:"uid"`
	Role string `json:"role"`
	Plan string `json:"plan"`
	jwt.RegisteredClaims
}

func GenerateToken(uid uint, role, plan string) (string, error) {
	iMyClaims := MyClaims{
		Uid:  uid,
		Role: role,
		Plan: plan,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.tokenExpired_min") * time.Minute)),
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

	//switch {
	//case err == nil && token.Valid:
	//	//這裡代表是合法token 並且驗證通過
	//case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
	//	// Token is either expired or not active yet
	//	// 這裡表示token超時 要換token
	//
	//default:
	//	fmt.Println("Couldn't handle this token:", err)
	//}

	return iMyClaims, err
}

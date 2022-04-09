package jwt

import (
	"errors"
	"fmt"
	"github.com/SunMaybo/zero/common/err_status"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type ThmusJWTClaims struct {
	Payload map[string]interface{}
	jwt.StandardClaims
}

type JsonWebToken struct {
	secretKey    string
	expireSecond int64
}

func New(secretKey string, expireSecond int64) *JsonWebToken {
	return &JsonWebToken{
		secretKey:    secretKey,
		expireSecond: expireSecond,
	}
}
func NewParseToken(secretKey string) *JsonWebToken {
	return &JsonWebToken{
		secretKey: secretKey,
	}
}

func (j *JsonWebToken) GenerateToken(payload map[string]interface{}) (string, int64, error) {
	expire := time.Now().Unix() + j.expireSecond
	// 将 uid，用户角色， 过期时间作为数据写入 token 中
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ThmusJWTClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
		},
	})
	// SecretKey 用于对用户数据进行签名，不能暴露
	code, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", 0, err
	}

	return code, expire, nil
}

func (j *JsonWebToken) ParseToken(tokenString string) (*ThmusJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ThmusJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ThmusJWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("parse token error")
	}
}

func (j *JsonWebToken) Verify(token string) (map[string]interface{}, *err_status.ErrMsg) {
	if strings.HasPrefix(token, "Bearer") {
		token = strings.TrimSpace(token[len("Bearer"):])
	}
	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, err_status.NewWithMsg(2001, "token is not available")
	}

	return claims.Payload, nil
}

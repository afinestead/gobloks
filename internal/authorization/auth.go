package authorization

import (
	"errors"
	"gobloks/internal/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateAccessToken(pid types.PlayerID, gid types.GameID, ttl uint) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["pid"] = pid
	claims["gid"] = gid
	claims["exp"] = time.Now().Add(time.Duration(ttl) * time.Second)

	secretKey, err := os.ReadFile("key.priv")
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyAccessToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return "", nil
	})
}

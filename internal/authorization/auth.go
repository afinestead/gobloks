package authorization

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gobloks/internal/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateAccessToken(pid types.PlayerID, gid types.GameID, ttl uint) (string, error) {

	secretKey, err := os.ReadFile("key.priv")
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"pid": pid,
			"gid": gid,
			"exp": time.Now().Add(time.Duration(ttl) * time.Second),
		},
	)

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

func GenerateKey() (string, error) {

	generateRandomBytes := func(length int) ([]byte, error) {
		randomBytes := make([]byte, length)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return nil, err
		}
		return randomBytes, nil
	}

	randomBytes, err := generateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

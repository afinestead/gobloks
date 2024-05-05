package authorization

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"gobloks/internal/types"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const KeyLocation = "key.priv"
const AccessTokenHeader = "access_token"

var secretKey []byte

func AuthMiddleware(noAuth []gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, handler := range noAuth {
			// go can't compare function pointers??
			if reflect.ValueOf(c.Handler()).Pointer() == reflect.ValueOf(handler).Pointer() {
				c.Next()
				return
			}
		}

		token, err := verifyAccessToken(c.GetHeader(AccessTokenHeader))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "token invalid"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("pid", types.PlayerID(claims["PlayerId"].(float64)))
		c.Set("gid", types.GameID(claims["GameId"].(string)))
		c.Next()
	}
}

func SetupKeys() {
	if _, err := os.ReadFile(KeyLocation); err != nil {
		file, err := os.OpenFile(KeyLocation, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		key, err := generateKey()
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(key)
		if err != nil {
			log.Fatal(err)
		}
	}

	secretKey, _ = os.ReadFile(KeyLocation)
}

func CreateAccessToken(pid types.PlayerID, gid types.GameID, ttl int) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"PlayerId":  pid,
			"GameId":    gid,
			"ExpiresAt": time.Now().Add(time.Duration(ttl) * time.Second),
		},
	)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyAccessToken(token string) (*jwt.Token, error) {
	fmt.Println(token)
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return secretKey, nil
	})
}

func generateKey() (string, error) {

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

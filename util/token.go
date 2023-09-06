package util

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateJWT(c *gin.Context, userId uuid.UUID) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_TOKEN")))
	if err != nil {
		return err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt-token", tokenString, (3600 * 24), "", "", false, true)

	return nil
}

func ParsesAndValidateJWT(c *gin.Context) (string, float64, error) {
	tokenString, err := c.Cookie("jwt-token")
	if err != nil {
		return "", 0, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", 0, err
	}

	return claims["user_id"].(string), claims["exp"].(float64), nil
}

package baseClass

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("goodAndgood")

// Claims 结构体，用于生成JWT
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT 生成JWT
func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(99999999 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// ValidateJWT 验证JWT
func ValidateJWT(c *gin.Context) error {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return errors.New("Authorization missing")
	}

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Bad token"})
		c.Abort()
		return err
	}

	c.Set("userID", claims.UserID)
	c.Next()

	return nil
}

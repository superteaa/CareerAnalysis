package baseClass

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("goodAndgood")

// Claims 结构体，用于生成JWT
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT 生成JWT
func GenerateJWT(userID int) (string, error) {
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
func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		db := GetDB()
		result := map[string]interface{}{
			"id":       0,
			"token":    "",
			"username": "",
		}
		db_result := db.Raw("SELECT * FROM user WHERE token = ?", tokenString).Scan(&result)
		if db_result.Error != nil {
			log.Println("jwt处理有问题:", db_result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
			c.Abort()
			return
		}
		if db_result.RowsAffected != 1 {
			log.Println("token无效", tokenString)
			c.JSON(http.StatusOK, gin.H{"error": "Bad token"})
			c.Abort()
			return
		}

		claims := &Claims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "Bad token"})
			c.Abort()
			return
		}

		c.Set("userID", result["id"])
		log.Println("token鉴权", result["id"], tokenString)
		c.Next()
		// return
	}
}

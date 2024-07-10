package baseClass

import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

var jwtKey = []byte("goodAndgood")

// Claims 结构体，用于生成JWT
type Claims struct {
    UserID uint `json:"user_id"`
    jwt.StandardClaims
}

// GenerateJWT 生成JWT
func GenerateJWT(userID uint) (string, error) {
    expirationTime := time.Now().Add(30 * time.Minute)
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
func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        return nil, err
    }
    if !token.Valid {
        return nil, err
    }
    return claims, nil
}

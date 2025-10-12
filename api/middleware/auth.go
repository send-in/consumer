package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(passkey string) gin.HandlerFunc {
	key := []byte(passkey)

	return func(context *gin.Context){
			
		token := context.GetHeader("Authorization")
		
		if token == "" {
			context.AbortWithStatusJSON(
				http.StatusUnauthorized, 
				gin.H{"error": "Unauthorized"},
			)
			return
		}
		
		token = token[len("Bearer "):]
		err := verifyToken(token, key)
		
		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusUnauthorized, 
				gin.H{"error": err.Error()},
			)
			return
		}
			
		context.Next()
	}
}

func verifyToken(token string, key []byte) error {
	parsed, err := jwt.Parse(
		token, 
		func(parsed *jwt.Token) (any, error) {
			return key, nil
		},
	)

	if err != nil {
		return err
	}

	if !parsed.Valid {
		return errors.New("invalid token")
	}

	return nil
}
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/utils"
)

func JWTMiddleware(jwtSecretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Default().Println("This is the header: " + authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// parse the jwt
		// 注意要 pass the address
		// p := &x        // p 是 *int，指向 x 的記憶體位置
		claim := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		// Valid specifies if the token is valid.  Populated when you Parse/Verify a token
		if err != nil || !token.Valid {
			log.Default().Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// 這邊完成後 token.Claims 會變成我們的 claims 實際行別是 *utils.Claims

		if claims, ok := token.Claims.(*utils.Claims); ok {
			c.Set("email", claims.Email)
			c.Set("userID", claims.ID)
			log.Default().Println(fmt.Sprintf("email %s:", claims.Email))
			log.Default().Println(fmt.Sprintf("userID %d:", (claims.UserID)))
		}

		c.Next()
	}
}

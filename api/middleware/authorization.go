package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-fast/internal/tools"
	"net/http"
)

var jwtSecretKey string

func init() {
	jwtSecretKey = viper.GetString("app.jwtSecretKey")
}

func AuthorizationJwt(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	accessToken, err := tools.ParseAuthHeader(auth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

	loginJwt := tools.JWT{SigningKey: []byte(jwtSecretKey)}
	token, err := loginJwt.ParseToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "token无效",
		})
		c.Abort()
		return
	}

	userId := int(claims["user_id"].(float64))
	c.Set("user_id", userId)
	c.Next()
}

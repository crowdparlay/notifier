package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func AuthorizeJWT(c *gin.Context) (string, bool) {
	id := c.Param("id")

	header := strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(header, jwt.MapClaims{})

	if err != nil {
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", false
	}

	if claims["sub"] == id {
		return id, false
	}

	return "", false
}

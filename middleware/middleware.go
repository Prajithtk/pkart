package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	ID    uint
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func JwtToken(c *gin.Context, id uint, email string, role string) (string, error) {
	claims := Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 240).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signedToken, err := c.Cookie("JwtToken" + requiredRole)
		// fmt.Println(signedToken)
		if err != nil {
			c.JSON(401, gin.H{
				"Status":  "Unauthorized",
				"Code":    401,
				"Message": "can't find cookie try again",
				"Data":    gin.H{},
			})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"Status":  "Unauthorized",
				"Code":    401,
				"Message": "Invalid token",
				"Data":    gin.H{},
			})
			c.Abort()
		}
		if claims.Role != requiredRole {
			c.JSON(403, gin.H{
				"Status":  "Forbidden",
				"Code":    401,
				"Message": "No permission",
				"Data":    gin.H{},
			})
			c.Abort()
		}

		c.Set("userid", claims.ID)
		c.Set("useremail", claims.Email)
		c.Next()
	}
}

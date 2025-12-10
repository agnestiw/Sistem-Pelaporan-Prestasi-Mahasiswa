package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID, roleID string, roleName string, studentID *string, advisor_id *string ,permissions []string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"role_id":     roleID,
		"role_name":   roleName,
		"permissions": permissions,
		"student_id":  studentID,
		"advisor_id":  advisor_id,
		"exp":         time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("API_SECRET")
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("API_SECRET")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

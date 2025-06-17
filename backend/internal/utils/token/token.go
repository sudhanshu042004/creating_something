package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sudhanshu042004/sandbox/internal/config"
)

func CreateToken(id int64, email string) (string, error) {
	secretByte := []byte(config.Cfg.TokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    id,
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

	tokenString, err := token.SignedString(secretByte)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return config.Cfg.TokenSecret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid token")
	}
	return nil
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "authCookie",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
}

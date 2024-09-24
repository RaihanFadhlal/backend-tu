package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

var JwtKey = []byte("di93mji439rj3489r")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var rdb *redis.Client

type contextKey string

const UserEmailKey contextKey = "userEmail"

func GenerateToken(email string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		}

		isBlacklisted, err := IsTokenBlacklisted(tokenString)
		if err != nil {
			http.Error(w, "Error checking token blacklist", http.StatusInternalServerError)
			return
		}
		if isBlacklisted {
			http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}

func IsTokenBlacklisted(tokenString string) (bool, error) {
	ctx := context.Background()
	result, err := rdb.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return result == "blacklisted", nil
}

func BlacklistToken(tokenString string, duration time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, tokenString, "blacklisted", duration).Err()
}



package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
    secretKey        []byte
    refreshSecretKey []byte
}

func NewJWTService() services.IJWTService {
    return &jwtService{
        secretKey:        []byte(os.Getenv("JWT_SECRET_KEY")),
        refreshSecretKey: []byte(os.Getenv("JWT_REFRESH_SECRET_KEY")),
    }
}

func (s *jwtService) GenerateToken(ctx context.Context, userID string, roles []string, exp int) (*string, error) {
    claims := jwt.MapClaims{
        "sub":   userID,
        "roles": strings.Join(roles, ","),
        "type":  "access",
        "iat":   time.Now().Unix(),
        "exp":   time.Now().Add(time.Second * time.Duration(exp)).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.secretKey)
    if err != nil {
        return nil, fmt.Errorf("error creating access token: %w", err)
    }

    return &tokenString, nil
}

func (s *jwtService) GenerateRefreshToken(ctx context.Context, userID string, exp int) (*string, error) {
    claims := jwt.MapClaims{
        "sub":  userID,
        "type": "refresh",
        "iat":  time.Now().Unix(),
        "exp":  time.Now().Add(time.Second * time.Duration(exp)).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.refreshSecretKey)
    if err != nil {
        return nil, fmt.Errorf("error creating refresh token: %w", err)
    }

    return &tokenString, nil
}

func (s *jwtService) ExtractClaims(ctx context.Context, token string) (map[string]interface{}, error) {
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, fmt.Errorf("token expired")
        }

        return nil, fmt.Errorf("error parsing token: %w", err)
    }

    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

func (s *jwtService) ExtractRefreshClaims(ctx context.Context, token string) (map[string]interface{}, error) {
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.refreshSecretKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("error parsing refresh token: %w", err)
    }

    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
        if tokenType, exists := claims["type"]; !exists || tokenType != "refresh" {
            return nil, fmt.Errorf("invalid token type")
        }
        return claims, nil
    }

    return nil, fmt.Errorf("invalid refresh token")
}
package services

import "context"

type IJWTService interface {
	GenerateToken(ctx context.Context, userID string, roles []string, exp int) (*string, error)
	GenerateRefreshToken(ctx context.Context, userID string, exp int) (*string, error)
	ExtractClaims(ctx context.Context, token string) (map[string]interface{}, error)
	ExtractRefreshClaims(ctx context.Context, token string) (map[string]interface{}, error)
}
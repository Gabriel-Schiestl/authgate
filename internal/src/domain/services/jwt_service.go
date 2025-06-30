package services

type IJWTService interface {
	GenerateToken(userID string, roles []string, exp int) (*string, error)
	GenerateRefreshToken(userID string, exp int) (*string, error)
	ExtractClaims(token string) (map[string]interface{}, error)
	ExtractRefreshClaims(token string) (map[string]interface{}, error)
}
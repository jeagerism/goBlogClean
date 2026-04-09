package jwtclaims

import "github.com/golang-jwt/jwt/v5"

// AccessClaims is the signed payload for API access tokens.
type AccessClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     bool   `json:"role"`
	jwt.RegisteredClaims
}

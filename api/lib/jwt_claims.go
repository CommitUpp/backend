package lib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type accessTokenClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
}

func jwtClaims(accessToken string) (*accessTokenClaims, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid jwt format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims accessTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

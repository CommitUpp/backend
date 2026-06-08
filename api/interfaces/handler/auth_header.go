package handler

import "strings"

func bearerToken(authHeader string) string {
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" || token == authHeader {
		return ""
	}

	return token
}

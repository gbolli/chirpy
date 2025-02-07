package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {

	header := headers.Get("Authorization")
	
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer "), nil
	}

	return "", fmt.Errorf("no token found")
}
package serve

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Claims struct {
	Exp int64 `json:"exp"`
}

var Jwt string

func IsJWTExpired(clientSession string) (Jwt string, err error) {
	if Jwt == "" {
		Jwt, err = GetJwtToken(clientSession)
		return
	}
	parts := strings.Split(Jwt, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT format. Expected format: header.payload.signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Print(err)
		return "", err
	}
	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		log.Print(err)
		return "", err
	}
	expTime := time.Unix(claims.Exp, 0)
	if time.Now().After(expTime) {
		Jwt, err = GetJwtToken(clientSession)
		return
	}
	return Jwt, nil
}

package utils

import (
	"context"
	"fmt"
	"google.golang.org/api/idtoken"
)

type GooglePayload struct {
	Sub     string
	Email   string
	Name    string
	Picture string
}

func VerifyGoogleIDToken(c context.Context, idToken string) (*GooglePayload, error) {
	payload, err := idtoken.Validate(c, idToken, "")
	if err != nil {
		return nil, err
	}

	fmt.Println(payload.Claims)

	return &GooglePayload{
		Sub:     payload.Subject,
		Email:   payload.Claims["email"].(string),
		Name:    payload.Claims["name"].(string),
		Picture: payload.Claims["picture"].(string),
	}, nil
}

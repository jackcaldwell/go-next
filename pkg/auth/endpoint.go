package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

/////////////
// Sign Up //
/////////////

type signUpRequest struct {
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type signUpResponse struct {
	SessionID string `json:"session_id"`
}

func makeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return s.SignUp(request.(signUpRequest))
	}
}

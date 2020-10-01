package auth

import (
	"github.com/jackcaldwell/go-next/pkg/session"
	"github.com/jackcaldwell/go-next/pkg/user"
)

// Service is the interface that providers auth service methods.
type Service interface {
	SignUp(r signUpRequest) (*signUpResponse, error)
}

type service struct {
	sessions session.Repository
	users    user.Repository
}

// NewService creates an auth service with the necessary dependencies
func NewService(sessions session.Repository, users user.Repository) Service {
	return &service{
		sessions: sessions,
		users:    users,
	}
}

func (s *service) SignUp(r signUpRequest) (*signUpResponse, error) {
	user, err := s.users.Create(r.Email, r.Password)

	if err != nil {
		return nil, err
	}

	res, err := s.sessions.Create(user.ID)
	if err != nil {
		return nil, err
	}

	return &signUpResponse{
		SessionID: res.ID,
	}, nil
}

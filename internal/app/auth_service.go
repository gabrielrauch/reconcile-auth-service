package app

import (
	"errors"

	"github.com/gabrielrauch/reconcile-auth-service/internal/domain/model"
	"github.com/gabrielrauch/reconcile-auth-service/internal/domain/ports"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   ports.UserRepository
	tokens ports.TokenProvider
	logger zerolog.Logger
}

func NewAuthService(r ports.UserRepository, t ports.TokenProvider, logger zerolog.Logger) *AuthService {
	return &AuthService{repo: r, tokens: t, logger: logger}
}

func (s *AuthService) Register(first_name, email, password, role string, requestID string) error {
	errChan := make(chan error, 1)

	go func() {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
		user := &model.User{FirstName: first_name, Email: email, Password: string(hash), Role: role, IsActive: true}

		s.logger.Info().Str("request_id", requestID).Str("email", email).Msg("registering user")

		err := s.repo.Create(user)
		errChan <- err
	}()

	err := <-errChan
	if err != nil {
		s.logger.Error().Err(err).Str("request_id", requestID).Msg("failed to register user")
		return err
	}

	s.logger.Info().Str("request_id", requestID).Str("email", email).Msg("user registered successfully")
	return nil
}

func (s *AuthService) Login(email, password string, requestID string) (string, error) {
	tokenChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		user, err := s.repo.FindByEmail(email)
		if err != nil {
			errChan <- errors.New("invalid credentials")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			errChan <- errors.New("invalid credentials")
			return
		}

		token, err := s.tokens.Generate(user.Email)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- token
	}()

	select {
	case token := <-tokenChan:
		s.logger.Info().Str("request_id", requestID).Str("email", email).Msg("user logged in successfully")
		return token, nil
	case err := <-errChan:
		s.logger.Error().Err(err).Str("request_id", requestID).Msg("failed to login user")
		return "", err
	}
}

func (s *AuthService) ValidateToken(token string, requestID string) bool {
	isValid := s.tokens.Validate(token)
	if !isValid {
		s.logger.Warn().Str("request_id", requestID).Msg("invalid token")
	}
	return isValid
}

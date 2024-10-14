package service

import (
	"context"
	"time"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type AuthService struct {
	repo AuthRepo
	log  *zerolog.Logger
}

func NewAuthService(repo AuthRepo, log *zerolog.Logger) *AuthService {
	return &AuthService{
		repo: repo,
		log:  log,
	}
}

type AuthRepo interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Get(ctx context.Context, user models.User) (models.User, error)
}

func (s *AuthService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	user, err := s.repo.Create(ctx, user)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create new user")
		return user, err
	}
	s.log.Debug().Msg("created new user")
	return user, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, user models.User) (string, error) {
	user, err := s.repo.Get(ctx, user)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create new user")
		return "", err
	}
	s.log.Debug().Msg("created new user")

	token, err := NewToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

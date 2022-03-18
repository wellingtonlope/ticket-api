package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Payload struct {
	ID        string
	Name      string
	Profile   string
	ExpiresAt time.Time
}

func (payload Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return security.ErrUnauthorized
	}
	return nil
}

type Authenticator struct {
	userRepository repository.UserRepository
	secret         string
	duration       time.Duration
}

func NewAuthenticator(userRepository repository.UserRepository, secret string, duration time.Duration) *Authenticator {
	return &Authenticator{
		userRepository: userRepository,
		secret:         secret,
		duration:       duration,
	}
}

func (j *Authenticator) Generate(userId string) (string, error) {
	user, err := j.userRepository.GetByID(userId)
	if err != nil {
		return "", security.ErrUnauthorized
	}
	payload := Payload{
		ID:        user.ID,
		Name:      user.Name,
		Profile:   string(user.Profile),
		ExpiresAt: time.Now().Add(j.duration),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := at.SignedString([]byte(j.secret))
	if err != nil {
		return "", security.ErrUnauthorized
	}

	return token, nil
}

func (j *Authenticator) Validate(token string) (*domain.User, error) {
	payload := Payload{}
	_, err := jwt.ParseWithClaims(token, &payload, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, security.ErrUnauthorized
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, security.ErrUnauthorized
	}

	user, err := j.userRepository.GetByID(payload.ID)
	if err != nil {
		return nil, security.ErrUnauthorized
	}

	return user, nil
}

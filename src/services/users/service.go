package users

import (
	"context"
	"crypto/md5"
	"fmt"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBadInputPassword = fmt.Errorf("bad bassword")
	ErrBadInputEmail    = fmt.Errorf("bad email")
	ErrInternal         = fmt.Errorf("internal error")
	ErrUserExists       = fmt.Errorf("user with this email already exists")
	ErrUserNotExists    = fmt.Errorf("user with this email not exists")
	ErrWrongPassword    = fmt.Errorf("wrong user password")
	ErrAuthorization    = fmt.Errorf("bad auth token")
)

const SessionKeySalt = "SomeKeySalt"

func NewService(
	userRepo UserRepo, sessionRepo SessionRepo,
	log logging.Logger, tracer trace.Tracer, meter metrics.Meter,
) *Service {
	return &Service{userRepo, sessionRepo, log, tracer}
}

type Service struct {
	userRepo    UserRepo
	sessionRepo SessionRepo
	log         logging.Logger
	tracer      trace.Tracer
}

func (s *Service) Create(ctx context.Context, params CreateUserParams) (*UserDTO, error) {
	email, err := ValidateEmail(params.Email)
	if err != nil {
		return nil, ErrBadInputEmail
	}
	password, err := ValidatePassword(params.Password)
	if err != nil {
		return nil, ErrBadInputPassword
	}

	existsingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInternal
	}
	if existsingUser != nil {
		return nil, ErrUserExists
	}

	secret, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	result, err := s.userRepo.CreateUser(ctx, UserDTO{
		Email:  email,
		Secret: string(secret),
	})
	if err != nil {
		return nil, ErrInternal
	}

	return result, nil
}

func (s *Service) SendVerify(ctx context.Context, email string) error {
	return fmt.Errorf("not implemented")
}

func (s *Service) SendRestore(ctx context.Context, email string) error {
	return fmt.Errorf("not implemented")
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInternal
	}
	if user == nil {
		return nil, ErrUserNotExists
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(password))
	if err != nil {
		return nil, ErrWrongPassword
	}

	sessionKey := uuid.New().String()
	hashedKey := s.hashSessionKey(sessionKey)
	err = s.sessionRepo.PutSession(ctx, hashedKey, SessionDTO{UserId: user.Id}, 12*time.Hour)
	if err != nil {
		return nil, ErrInternal
	}

	return &LoginResult{
		User:  *user,
		Token: sessionKey,
	}, nil
}

func (s *Service) hashSessionKey(key string) string {
	key = key + SessionKeySalt
	keyHashed := md5.Sum([]byte(key))
	return string(keyHashed[:])
}

func (s *Service) Authorize(ctx context.Context, token string) (*SessionDTO, error) {
	hashedKey := s.hashSessionKey(token)
	user, err := s.sessionRepo.GetSession(ctx, hashedKey)
	if err != nil {
		return nil, ErrInternal
	}
	if user == nil {
		return nil, ErrAuthorization
	}
	return user, nil
}

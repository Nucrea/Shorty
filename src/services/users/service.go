package users

import (
	"context"
	"crypto/md5"
	"fmt"
	"shorty/src/common"
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
	logger logging.Logger, tracer trace.Tracer, meter metrics.Meter,
) *Service {
	return &Service{userRepo, sessionRepo, logger.WithService("users"), tracer}
}

type Service struct {
	userRepo    UserRepo
	sessionRepo SessionRepo
	logger      logging.Logger
	tracer      trace.Tracer
}

func (s *Service) Create(ctx context.Context, params CreateUserParams) (*UserDTO, error) {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "users::Create")
	defer span.End()

	email, err := common.ValidateEmail(params.Email)
	if err != nil {
		return nil, ErrBadInputEmail
	}
	password, err := common.ValidatePassword(params.Password)
	if err != nil {
		return nil, ErrBadInputPassword
	}

	existsingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting user with email=%s from repo", email)
		return nil, ErrInternal
	}
	if existsingUser != nil {
		log.Info().Msgf("user with email=%s already exists", email)
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
	log := s.logger.WithContext(ctx)
	ctx, span := s.tracer.Start(ctx, "users::Login")
	defer span.End()

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting user with email=%s from repo", email)
		return nil, ErrInternal
	}
	if user == nil {
		log.Info().Msgf("user with email=%s dost not exist", email)
		return nil, ErrUserNotExists
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(password))
	if err != nil {
		log.Info().Msgf("wrong password for user with id=%s", user.Id)
		return nil, ErrWrongPassword
	}

	sessionKey := uuid.New().String()
	hashedKey := s.hashSessionKey(sessionKey)
	err = s.sessionRepo.PutSession(ctx, hashedKey, SessionDTO{UserId: user.Id}, 12*time.Hour)
	if err != nil {
		log.Error().Err(err).Msgf("failed saving auth session for user with id=%s", user.Id)
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
	log := s.logger.WithContext(ctx)
	ctx, span := s.tracer.Start(ctx, "users::Authorize")
	defer span.End()

	hashedKey := s.hashSessionKey(token)
	session, err := s.sessionRepo.GetSession(ctx, hashedKey)
	if err != nil {
		log.Error().Err(err).Msg("failed getting auth session")
		return nil, ErrInternal
	}
	if session == nil {
		log.Info().Msgf("no such session with key=%s", common.MaskSecret(token))
		return nil, ErrAuthorization
	}
	return session, nil
}

func (s *Service) DeleteSession(ctx context.Context, token string) error {
	log := s.logger.WithContext(ctx)
	ctx, span := s.tracer.Start(ctx, "users::DeleteSession")
	defer span.End()

	err := s.sessionRepo.DelSession(ctx, token)
	if err != nil {
		log.Error().Err(err).Msg("failed clearing session")
	}
	return err
}

func (s *Service) GetById(ctx context.Context, userId string) (*UserDTO, error) {
	log := s.logger.WithContext(ctx)
	ctx, span := s.tracer.Start(ctx, "users::GetById")
	defer span.End()

	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting with id=%s", userId)
		return nil, ErrInternal
	}
	if user == nil {
		log.Info().Msgf("no such user with id=%s", userId)
		return nil, ErrAuthorization
	}

	return user, nil
}

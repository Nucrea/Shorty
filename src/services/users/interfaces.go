package users

import (
	"context"
	"time"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user UserDTO) (*UserDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDTO, error)
	GetUserById(ctx context.Context, id string) (*UserDTO, error)
	SetUserVerified(ctx context.Context, id string) error
}

type SessionRepo interface {
	PutSession(ctx context.Context, key string, session SessionDTO, expiration time.Duration) error
	GetSession(ctx context.Context, key string) (*SessionDTO, error)
	DelSession(ctx context.Context, key string) error
}

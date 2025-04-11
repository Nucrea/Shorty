package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmails(t *testing.T) {
	ctx := context.Background()
	userService := NewService(nil, nil, nil, nil, nil)

	password := "dsghvGHhj34!"
	tc := []string{"@aaa.aa", "bad2 @aaa.aa", "bad3@aa.", "bad4aa.aa", ""}

	for _, email := range tc {
		user, err := userService.Create(ctx, CreateUserParams{
			Email:    email,
			Password: password,
		})
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrBadInputEmail)
	}
}

func TestPasswords(t *testing.T) {
	ctx := context.Background()
	userService := NewService(nil, nil, nil, nil, nil)

	email := "test@example.com"
	tc := []string{"", "dd", "ffdsbdffsg", "ghvsadgGFCFGGD", "setg!@@@"}

	for _, pass := range tc {
		user, err := userService.Create(ctx, CreateUserParams{
			Email:    email,
			Password: pass,
		})
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrBadInputPassword)
	}
}

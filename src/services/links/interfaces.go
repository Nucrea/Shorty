package links

//go:generate mockgen -package links -destination service_mock_test.go -source=$GOFILE Storage

import "context"

type Storage interface {
	SaveLink(ctx context.Context, id, url string) error
	SaveLinkForUser(ctx context.Context, id, userId, url string) error
	GetLinkById(ctx context.Context, id string) (*LinkDTO, error)
	GetLinksByUserId(ctx context.Context, userId string) ([]*LinkDTO, error)
	DeleteLinks(ctx context.Context, ids ...string) error
}

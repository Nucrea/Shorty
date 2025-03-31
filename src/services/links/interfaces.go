package links

import "context"

type Storage interface {
	SaveShortlink(ctx context.Context, id, url string) error
	GetShortlink(ctx context.Context, id string) (string, error)
}

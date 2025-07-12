package broker

import "context"

type Broker interface {
	PutFilesToDelete(ctx context.Context, name ...string) error
}

package links

import (
	"context"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/common/tracing"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newMock(t *testing.T) (*Service, *MockStorage) {
	logger := logging.NewNoop()
	tracer := tracing.NewNoopTracer()
	meter := metrics.NewNoop()
	ctrl := gomock.NewController(t)
	storage := NewMockStorage(ctrl)
	return NewService(storage, logger, tracer, meter), storage
}

func TestSaveLink(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should create link", func(t *testing.T) {
		assert := assert.New(t)
		service, storage := newMock(t)

		linkId := ""
		url := "http://example.com"
		storage.EXPECT().SaveLink(gomock.Any(), gomock.Any(), url).DoAndReturn(
			func(ctx context.Context, id string, url string) (*LinkDTO, error) {
				linkId = id
				return &LinkDTO{Id: id, Url: url}, nil
			},
		)

		link, err := service.Save(ctx, url, nil)
		assert.NoError(err)
		assert.NotNil(link)
		assert.Equal(url, link.Url)
		assert.Equal(linkId, link.Id)
		assert.Empty(link.UserId)
	})
}

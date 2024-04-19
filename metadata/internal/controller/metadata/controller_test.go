package metadata

import (
	"context"
	"errors"
	"movie-micro/metadata/internal/repository"
	"movie-micro/metadata/pkg/model"
	"testing"

	gen "movie-micro/gen/mock/metadata/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	// "github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	testCases := []struct {
		desc       string
		expRepoRes *model.Metadata
		expRepoErr error
		wantRes    *model.Metadata
		wantErr    error
	}{
		{
			desc:       "not found",
			expRepoErr: repository.ErrNotFound,
			wantErr:    ErrNotFound,
		},
		{
			desc:       "unexpected error",
			expRepoErr: errors.New("unexpected error"),
			wantErr:    errors.New("unexpected error"),
		},
		{
			desc:       "success",
			expRepoRes: &model.Metadata{},
			wantRes:    &model.Metadata{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := gen.NewMockmetadataRepository(ctrl)
			c := New(repoMock)
			ctx := context.Background()
			id := "id"
			repoMock.EXPECT().Get(ctx, id).Return(tt.expRepoRes, tt.expRepoErr)
			res, err := c.Get(ctx, id)
			assert.Equal(t, tt.wantRes, res, tt.desc)
			assert.Equal(t, tt.wantErr, err, tt.desc)
		})
	}
}

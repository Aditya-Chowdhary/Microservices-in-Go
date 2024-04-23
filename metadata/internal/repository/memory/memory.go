package memory

import (
	"context"
	"sync"

	"github.com/Aditya-Chowdhary/micro-movies/metadata/internal/repository"
	"github.com/Aditya-Chowdhary/micro-movies/metadata/pkg/model"

	"go.opentelemetry.io/otel"
)

const tracerID = "metadata-repository-memory"

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Get")
	defer span.End()

	r.RLock()
	defer r.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Put")
	defer span.End()
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}

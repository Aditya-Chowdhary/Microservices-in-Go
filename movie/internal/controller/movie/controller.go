package movie

import (
	"context"
	"errors"
	metadatamodel "movie-micro/metadata/pkg/model"
	"movie-micro/movie/internal/gateway"
	"movie-micro/movie/pkg/model"
	ratingmodel "movie-micro/rating/pkg/model"
	"sync"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	// PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// Controller defines a movie service controller
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller
func New(ratingGateway ratingGateway, metametadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway,
		metametadataGateway,
	}
}

// Get returns the movies details including the aggregated rating and movie metadata
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	var metadata *metadatamodel.Metadata
	var getMetadataErr error
	var rating float64
	var getRatingErr error
	go func() {
		defer wg.Done()
		metadata, getMetadataErr = c.metadataGateway.Get(ctx, id)
	}()
	go func() {
		defer wg.Done()
		rating, getRatingErr = c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	}()
	wg.Wait()

	if err := getMetadataErr; err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{Metadata: *metadata}
	if err := getRatingErr; err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Proceed. It is ok to not have ratings
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}

	return details, nil
}

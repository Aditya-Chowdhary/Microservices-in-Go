package mysql

import (
	"context"
	"database/sql"
	"movie-micro/rating/internal/repository"
	"movie-micro/rating/pkg/model"

	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MYSQL-based rating repository
type Repository struct {
	db *sql.DB
}

// New creates a new MYSQL-based rating repository
func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Get retrieves all ratings for a given record
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	query := "SELECT user_id, value FROM ratings WHERE record_id = ? AND record_type = ?"

	rows, err := r.db.QueryContext(ctx, query, recordID, recordType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Rating
	for rows.Next() {
		var (
			user_id string
			value   int32
		)
		if err := rows.Scan(&user_id, &value); err != nil {
			return nil, err
		}
		res = append(res, model.Rating{
			UserID: model.UserID(user_id),
			Value:  model.RatingValue(value),
		})
	}
	if len(res) == 0 {
		return nil, repository.ErrNotFound
	}
	return res, nil
}

// Put adds rating for a given record
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	query := `INSERT INTO ratings (record_id, record_type, user_id, value)
	VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, recordID, recordType, rating.UserID, rating.Value)
	return err
}

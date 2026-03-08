package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type RatingRepository struct {
	db *sql.DB
}

func NewRatingRepository(db *sql.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) Create(rating *domain.Rating) error {
	// Generate UUID v7 (time-ordered for better index performance)
	ratingID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	rating.ID = ratingID
	
	query := `INSERT INTO ratings (id, user_id, profile_id, content_id, rating, review, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	now := time.Now()
	_, err = r.db.Exec(query, rating.ID, rating.UserID, rating.ProfileID, rating.ContentID,
		rating.Rating, rating.Review, now, now)
	if err != nil {
		return err
	}
	
	rating.CreatedAt = now
	rating.UpdatedAt = now
	return nil
}

func (r *RatingRepository) GetByContentID(contentID string) ([]*domain.Rating, error) {
	query := `SELECT id, user_id, profile_id, content_id, rating, review, created_at, updated_at 
	          FROM ratings WHERE content_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*domain.Rating
	for rows.Next() {
		rating := &domain.Rating{}
		err := rows.Scan(
			&rating.ID, &rating.UserID, &rating.ProfileID, &rating.ContentID,
			&rating.Rating, &rating.Review, &rating.CreatedAt, &rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

func (r *RatingRepository) GetByUserID(userID string, limit, offset int) ([]*domain.Rating, error) {
	query := `SELECT id, user_id, profile_id, content_id, rating, review, created_at, updated_at 
	          FROM ratings WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*domain.Rating
	for rows.Next() {
		rating := &domain.Rating{}
		err := rows.Scan(
			&rating.ID, &rating.UserID, &rating.ProfileID, &rating.ContentID,
			&rating.Rating, &rating.Review, &rating.CreatedAt, &rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

func (r *RatingRepository) GetUserRating(userID, contentID string) (*domain.Rating, error) {
	query := `SELECT id, user_id, profile_id, content_id, rating, review, created_at, updated_at 
	          FROM ratings WHERE user_id = $1 AND content_id = $2`
	
	rating := &domain.Rating{}
	err := r.db.QueryRow(query, userID, contentID).Scan(
		&rating.ID, &rating.UserID, &rating.ProfileID, &rating.ContentID,
		&rating.Rating, &rating.Review, &rating.CreatedAt, &rating.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rating, nil
}

func (r *RatingRepository) Update(rating *domain.Rating) error {
	query := `UPDATE ratings SET rating = $1, review = $2, updated_at = $3 WHERE id = $4`
	rating.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, rating.Rating, rating.Review, rating.UpdatedAt, rating.ID)
	return err
}

func (r *RatingRepository) Delete(id string) error {
	query := `DELETE FROM ratings WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *RatingRepository) GetAverageRating(contentID string) (float64, error) {
	query := `SELECT COALESCE(AVG(rating), 0) FROM ratings WHERE content_id = $1`
	var avg float64
	err := r.db.QueryRow(query, contentID).Scan(&avg)
	return avg, err
}

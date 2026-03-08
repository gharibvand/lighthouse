package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type ContentRepository struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) Create(content *domain.Content) error {
	// Generate UUID v7 (time-ordered for better index performance)
	contentID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	content.ID = contentID
	
	query := `INSERT INTO content (id, type, title, description, release_year, duration, rating, poster_url, 
	          thumbnail_url, trailer_url, video_path, is_active, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	
	now := time.Now()
	_, err = r.db.Exec(query, content.ID, content.Type, content.Title, content.Description, content.ReleaseYear,
		content.Duration, content.Rating, content.PosterURL, content.ThumbnailURL, content.TrailerURL,
		content.VideoPath, content.IsActive, now, now)
	if err != nil {
		return err
	}
	
	content.CreatedAt = now
	content.UpdatedAt = now
	return nil
}

func (r *ContentRepository) GetByID(id string) (*domain.Content, error) {
	query := `SELECT id, type, title, description, release_year, duration, rating, poster_url, 
	          thumbnail_url, trailer_url, video_path, is_active, created_at, updated_at 
	          FROM content WHERE id = $1 AND is_active = true`
	content := &domain.Content{}
	err := r.db.QueryRow(query, id).Scan(
		&content.ID, &content.Type, &content.Title, &content.Description, &content.ReleaseYear,
		&content.Duration, &content.Rating, &content.PosterURL, &content.ThumbnailURL,
		&content.TrailerURL, &content.VideoPath, &content.IsActive, &content.CreatedAt, &content.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *ContentRepository) GetMovies(limit, offset int) ([]*domain.Content, error) {
	query := `SELECT id, type, title, description, release_year, duration, rating, poster_url, 
	          thumbnail_url, trailer_url, video_path, is_active, created_at, updated_at 
	          FROM content WHERE type = 'movie' AND is_active = true 
	          ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []*domain.Content
	for rows.Next() {
		content := &domain.Content{}
		err := rows.Scan(
			&content.ID, &content.Type, &content.Title, &content.Description, &content.ReleaseYear,
			&content.Duration, &content.Rating, &content.PosterURL, &content.ThumbnailURL,
			&content.TrailerURL, &content.VideoPath, &content.IsActive, &content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *ContentRepository) GetTVShows(limit, offset int) ([]*domain.Content, error) {
	query := `SELECT id, type, title, description, release_year, duration, rating, poster_url, 
	          thumbnail_url, trailer_url, video_path, is_active, created_at, updated_at 
	          FROM content WHERE type = 'tv_show' AND is_active = true 
	          ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []*domain.Content
	for rows.Next() {
		content := &domain.Content{}
		err := rows.Scan(
			&content.ID, &content.Type, &content.Title, &content.Description, &content.ReleaseYear,
			&content.Duration, &content.Rating, &content.PosterURL, &content.ThumbnailURL,
			&content.TrailerURL, &content.VideoPath, &content.IsActive, &content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *ContentRepository) Search(query string, limit, offset int) ([]*domain.Content, error) {
	searchQuery := `SELECT id, type, title, description, release_year, duration, rating, poster_url, 
	                thumbnail_url, trailer_url, video_path, is_active, created_at, updated_at 
	                FROM content 
	                WHERE is_active = true AND (title ILIKE $1 OR description ILIKE $1)
	                ORDER BY rating DESC, created_at DESC LIMIT $2 OFFSET $3`
	
	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []*domain.Content
	for rows.Next() {
		content := &domain.Content{}
		err := rows.Scan(
			&content.ID, &content.Type, &content.Title, &content.Description, &content.ReleaseYear,
			&content.Duration, &content.Rating, &content.PosterURL, &content.ThumbnailURL,
			&content.TrailerURL, &content.VideoPath, &content.IsActive, &content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *ContentRepository) GetTrending(limit int) ([]*domain.Content, error) {
	query := `SELECT id, type, title, description, release_year, duration, rating, poster_url, 
	          thumbnail_url, trailer_url, video_path, is_active, created_at, updated_at 
	          FROM content WHERE is_active = true 
	          ORDER BY rating DESC, created_at DESC LIMIT $1`
	
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []*domain.Content
	for rows.Next() {
		content := &domain.Content{}
		err := rows.Scan(
			&content.ID, &content.Type, &content.Title, &content.Description, &content.ReleaseYear,
			&content.Duration, &content.Rating, &content.PosterURL, &content.ThumbnailURL,
			&content.TrailerURL, &content.VideoPath, &content.IsActive, &content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *ContentRepository) GetByGenre(genreID int, limit, offset int) ([]*domain.Content, error) {
	query := `SELECT c.id, c.type, c.title, c.description, c.release_year, c.duration, c.rating, 
	          c.poster_url, c.thumbnail_url, c.trailer_url, c.video_path, c.is_active, c.created_at, c.updated_at 
	          FROM content c 
	          INNER JOIN content_genres cg ON c.id = cg.content_id 
	          WHERE cg.genre_id = $1 AND c.is_active = true 
	          ORDER BY c.rating DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, genreID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []*domain.Content
	for rows.Next() {
		content := &domain.Content{}
		err := rows.Scan(
			&content.ID, &content.Type, &content.Title, &content.Description, &content.ReleaseYear,
			&content.Duration, &content.Rating, &content.PosterURL, &content.ThumbnailURL,
			&content.TrailerURL, &content.VideoPath, &content.IsActive, &content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *ContentRepository) Update(content *domain.Content) error {
	query := `UPDATE content SET type = $1, title = $2, description = $3, release_year = $4, 
	          duration = $5, rating = $6, poster_url = $7, thumbnail_url = $8, trailer_url = $9, 
	          video_path = $10, is_active = $11, updated_at = $12 WHERE id = $13`
	content.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, content.Type, content.Title, content.Description, content.ReleaseYear,
		content.Duration, content.Rating, content.PosterURL, content.ThumbnailURL, content.TrailerURL,
		content.VideoPath, content.IsActive, content.UpdatedAt, content.ID)
	return err
}

func (r *ContentRepository) Delete(id string) error {
	query := `UPDATE content SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

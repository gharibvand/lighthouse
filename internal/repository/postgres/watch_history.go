package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type WatchHistoryRepository struct {
	db *sql.DB
}

func NewWatchHistoryRepository(db *sql.DB) *WatchHistoryRepository {
	return &WatchHistoryRepository{db: db}
}

func (r *WatchHistoryRepository) Create(history *domain.WatchHistory) error {
	// Generate UUID v7 (time-ordered for better index performance)
	historyID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	history.ID = historyID
	
	query := `INSERT INTO watch_history (id, user_id, profile_id, content_id, episode_id, progress, duration, watched_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	now := time.Now()
	_, err = r.db.Exec(query, history.ID, history.UserID, history.ProfileID, history.ContentID,
		history.EpisodeID, history.Progress, history.Duration, now, now)
	if err != nil {
		return err
	}
	
	history.WatchedAt = now
	history.UpdatedAt = now
	return nil
}

func (r *WatchHistoryRepository) GetByUserID(userID string, limit, offset int) ([]*domain.WatchHistory, error) {
	query := `SELECT id, user_id, profile_id, content_id, episode_id, progress, duration, watched_at, updated_at 
	          FROM watch_history WHERE user_id = $1 ORDER BY watched_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*domain.WatchHistory
	for rows.Next() {
		history := &domain.WatchHistory{}
		err := rows.Scan(
			&history.ID, &history.UserID, &history.ProfileID, &history.ContentID,
			&history.EpisodeID, &history.Progress, &history.Duration,
			&history.WatchedAt, &history.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (r *WatchHistoryRepository) GetByProfileID(profileID string, limit, offset int) ([]*domain.WatchHistory, error) {
	query := `SELECT id, user_id, profile_id, content_id, episode_id, progress, duration, watched_at, updated_at 
	          FROM watch_history WHERE profile_id = $1 ORDER BY watched_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, profileID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*domain.WatchHistory
	for rows.Next() {
		history := &domain.WatchHistory{}
		err := rows.Scan(
			&history.ID, &history.UserID, &history.ProfileID, &history.ContentID,
			&history.EpisodeID, &history.Progress, &history.Duration,
			&history.WatchedAt, &history.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (r *WatchHistoryRepository) GetByContentID(userID, contentID string) (*domain.WatchHistory, error) {
	query := `SELECT id, user_id, profile_id, content_id, episode_id, progress, duration, watched_at, updated_at 
	          FROM watch_history WHERE user_id = $1 AND content_id = $2 ORDER BY watched_at DESC LIMIT 1`
	
	history := &domain.WatchHistory{}
	err := r.db.QueryRow(query, userID, contentID).Scan(
		&history.ID, &history.UserID, &history.ProfileID, &history.ContentID,
		&history.EpisodeID, &history.Progress, &history.Duration,
		&history.WatchedAt, &history.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *WatchHistoryRepository) Update(history *domain.WatchHistory) error {
	query := `UPDATE watch_history SET progress = $1, duration = $2, episode_id = $3, updated_at = $4 WHERE id = $5`
	history.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, history.Progress, history.Duration, history.EpisodeID, history.UpdatedAt, history.ID)
	return err
}

func (r *WatchHistoryRepository) Delete(id string) error {
	query := `DELETE FROM watch_history WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

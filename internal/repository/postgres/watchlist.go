package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type WatchlistRepository struct {
	db *sql.DB
}

func NewWatchlistRepository(db *sql.DB) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}

func (r *WatchlistRepository) Create(watchlist *domain.Watchlist) error {
	// Generate UUID v7 (time-ordered for better index performance)
	watchlistID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	watchlist.ID = watchlistID
	
	query := `INSERT INTO watchlist (id, user_id, profile_id, content_id, added_at) 
	          VALUES ($1, $2, $3, $4, $5)`
	
	now := time.Now()
	_, err = r.db.Exec(query, watchlist.ID, watchlist.UserID, watchlist.ProfileID, watchlist.ContentID, now)
	if err != nil {
		return err
	}
	
	watchlist.AddedAt = now
	return nil
}

func (r *WatchlistRepository) GetByUserID(userID string, limit, offset int) ([]*domain.Watchlist, error) {
	query := `SELECT id, user_id, profile_id, content_id, added_at 
	          FROM watchlist WHERE user_id = $1 ORDER BY added_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var watchlists []*domain.Watchlist
	for rows.Next() {
		watchlist := &domain.Watchlist{}
		err := rows.Scan(
			&watchlist.ID, &watchlist.UserID, &watchlist.ProfileID,
			&watchlist.ContentID, &watchlist.AddedAt,
		)
		if err != nil {
			return nil, err
		}
		watchlists = append(watchlists, watchlist)
	}
	return watchlists, nil
}

func (r *WatchlistRepository) GetByProfileID(profileID string, limit, offset int) ([]*domain.Watchlist, error) {
	query := `SELECT id, user_id, profile_id, content_id, added_at 
	          FROM watchlist WHERE profile_id = $1 ORDER BY added_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, profileID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var watchlists []*domain.Watchlist
	for rows.Next() {
		watchlist := &domain.Watchlist{}
		err := rows.Scan(
			&watchlist.ID, &watchlist.UserID, &watchlist.ProfileID,
			&watchlist.ContentID, &watchlist.AddedAt,
		)
		if err != nil {
			return nil, err
		}
		watchlists = append(watchlists, watchlist)
	}
	return watchlists, nil
}

func (r *WatchlistRepository) Delete(userID, contentID string) error {
	query := `DELETE FROM watchlist WHERE user_id = $1 AND content_id = $2`
	_, err := r.db.Exec(query, userID, contentID)
	return err
}

func (r *WatchlistRepository) Exists(userID, contentID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM watchlist WHERE user_id = $1 AND content_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, userID, contentID).Scan(&exists)
	return exists, err
}

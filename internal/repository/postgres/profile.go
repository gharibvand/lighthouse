package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(profile *domain.Profile) error {
	// Generate UUID v7 (time-ordered for better index performance)
	profileID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	profile.ID = profileID
	
	query := `INSERT INTO profiles (id, user_id, name, avatar, is_child, language, subtitle_lang, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	now := time.Now()
	_, err = r.db.Exec(query, profile.ID, profile.UserID, profile.Name, profile.Avatar, 
		profile.IsChild, profile.Language, profile.SubtitleLang, now, now)
	if err != nil {
		return err
	}
	
	profile.CreatedAt = now
	profile.UpdatedAt = now
	return nil
}

func (r *ProfileRepository) GetByID(id string) (*domain.Profile, error) {
	query := `SELECT id, user_id, name, avatar, is_child, language, subtitle_lang, created_at, updated_at 
	          FROM profiles WHERE id = $1`
	profile := &domain.Profile{}
	err := r.db.QueryRow(query, id).Scan(
		&profile.ID, &profile.UserID, &profile.Name, &profile.Avatar,
		&profile.IsChild, &profile.Language, &profile.SubtitleLang,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileRepository) GetByUserID(userID string) ([]*domain.Profile, error) {
	query := `SELECT id, user_id, name, avatar, is_child, language, subtitle_lang, created_at, updated_at 
	          FROM profiles WHERE user_id = $1 ORDER BY created_at`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*domain.Profile
	for rows.Next() {
		profile := &domain.Profile{}
		err := rows.Scan(
			&profile.ID, &profile.UserID, &profile.Name, &profile.Avatar,
			&profile.IsChild, &profile.Language, &profile.SubtitleLang,
			&profile.CreatedAt, &profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (r *ProfileRepository) Update(profile *domain.Profile) error {
	query := `UPDATE profiles SET name = $1, avatar = $2, is_child = $3, language = $4, 
	          subtitle_lang = $5, updated_at = $6 WHERE id = $7`
	profile.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, profile.Name, profile.Avatar, profile.IsChild,
		profile.Language, profile.SubtitleLang, profile.UpdatedAt, profile.ID)
	return err
}

func (r *ProfileRepository) Delete(id string) error {
	query := `DELETE FROM profiles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

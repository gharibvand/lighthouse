package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     *string   `json:"email,omitempty"` // Optional - can be null if phone is provided
	Phone     *string   `json:"phone,omitempty"` // Optional - can be null if email is provided
	Password  string    `json:"-"` // Never return password in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	IsChild      bool      `json:"is_child"`
	Language     string    `json:"language"`
	SubtitleLang string    `json:"subtitle_lang"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type WatchHistory struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	ProfileID   string     `json:"profile_id"`
	ContentID   string     `json:"content_id"`
	EpisodeID   *string    `json:"episode_id,omitempty"` // For TV shows
	Progress    int        `json:"progress"`             // Seconds watched
	Duration    int        `json:"duration"`             // Total duration in seconds
	WatchedAt   time.Time  `json:"watched_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Watchlist struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProfileID string    `json:"profile_id"`
	ContentID string    `json:"content_id"`
	AddedAt   time.Time `json:"added_at"`
}

type Rating struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProfileID string    `json:"profile_id"`
	ContentID string    `json:"content_id"`
	Rating    int       `json:"rating"` // 1-5 stars
	Review    string    `json:"review,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package domain

import "time"

type ContentType string

const (
	ContentTypeMovie   ContentType = "movie"
	ContentTypeTVShow  ContentType = "tv_show"
)

type Content struct {
	ID          string      `json:"id"`
	Type        ContentType `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ReleaseYear int         `json:"release_year"`
	Duration    int         `json:"duration"` // Minutes for movies, total for TV shows
	Rating      float64     `json:"rating"`  // Average rating
	PosterURL   string      `json:"poster_url"`
	ThumbnailURL string     `json:"thumbnail_url"`
	TrailerURL  string      `json:"trailer_url"`
	VideoPath   string      `json:"-"` // Internal path, not exposed
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	
	// Relations
	Genres    []Genre    `json:"genres,omitempty"`
	Actors    []Actor    `json:"actors,omitempty"`
	Directors []Director `json:"directors,omitempty"`
	Episodes  []Episode  `json:"episodes,omitempty"` // For TV shows
}

type Episode struct {
	ID          string    `json:"id"`
	ContentID   string    `json:"content_id"`
	Season      int       `json:"season"`
	EpisodeNum  int       `json:"episode_num"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"` // Minutes
	ThumbnailURL string   `json:"thumbnail_url"`
	VideoPath   string    `json:"-"` // Internal path
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type Director struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type Subtitle struct {
	ID        string  `json:"id"`
	ContentID string  `json:"content_id"`
	EpisodeID *string `json:"episode_id,omitempty"`
	Language  string  `json:"language"`
	FilePath  string  `json:"file_path"`
}

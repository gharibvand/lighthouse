package domain

type StreamingQuality string

const (
	Quality480p  StreamingQuality = "480p"
	Quality720p  StreamingQuality = "720p"
	Quality1080p StreamingQuality = "1080p"
	Quality4K    StreamingQuality = "4k"
)

type Playlist struct {
	ContentID string            `json:"content_id"`
	EpisodeID *string           `json:"episode_id,omitempty"`
	Qualities []StreamingQuality `json:"qualities"`
	Segments  map[StreamingQuality][]Segment `json:"segments"`
}

type Segment struct {
	URL      string `json:"url"`
	Duration float64 `json:"duration"`
	Sequence int    `json:"sequence"`
}

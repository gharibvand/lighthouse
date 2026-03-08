package domain

type RecommendationType string

const (
	RecommendationTypeTrending     RecommendationType = "trending"
	RecommendationTypePersonalized RecommendationType = "personalized"
	RecommendationTypeSimilar      RecommendationType = "similar"
	RecommendationTypeBecauseWatched RecommendationType = "because_watched"
)

type Recommendation struct {
	Type     RecommendationType `json:"type"`
	Title    string             `json:"title"`
	Contents []Content          `json:"contents"`
}

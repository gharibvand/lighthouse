package service

import (
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
)

type RecommendationService struct {
	contentRepo      *postgres.ContentRepository
	watchHistoryRepo *postgres.WatchHistoryRepository
}

func NewRecommendationService(contentRepo *postgres.ContentRepository, watchHistoryRepo *postgres.WatchHistoryRepository) *RecommendationService {
	return &RecommendationService{
		contentRepo:      contentRepo,
		watchHistoryRepo: watchHistoryRepo,
	}
}

func (s *RecommendationService) GetRecommendations(userID string, limit int) ([]*domain.Content, error) {
	// Get user's watch history
	histories, err := s.watchHistoryRepo.GetByUserID(userID, 10, 0)
	if err != nil {
		return nil, err
	}

	// If no history, return trending
	if len(histories) == 0 {
		return s.contentRepo.GetTrending(limit)
	}

	// Simple recommendation: get trending content
	// In production, this would use more sophisticated algorithms
	return s.contentRepo.GetTrending(limit)
}

func (s *RecommendationService) GetBecauseYouWatched(userID string, contentID string, limit int) ([]*domain.Content, error) {
	// Get similar content based on genres, actors, etc.
	// For now, return trending
	return s.contentRepo.GetTrending(limit)
}

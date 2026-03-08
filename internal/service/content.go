package service

import (
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
)

type ContentService struct {
	contentRepo *postgres.ContentRepository
}

func NewContentService(contentRepo *postgres.ContentRepository) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
	}
}

func (s *ContentService) GetContentByID(id string) (*domain.Content, error) {
	return s.contentRepo.GetByID(id)
}

func (s *ContentService) GetMovies(limit, offset int) ([]*domain.Content, error) {
	return s.contentRepo.GetMovies(limit, offset)
}

func (s *ContentService) GetTVShows(limit, offset int) ([]*domain.Content, error) {
	return s.contentRepo.GetTVShows(limit, offset)
}

func (s *ContentService) Search(query string, limit, offset int) ([]*domain.Content, error) {
	return s.contentRepo.Search(query, limit, offset)
}

func (s *ContentService) GetTrending(limit int) ([]*domain.Content, error) {
	return s.contentRepo.GetTrending(limit)
}

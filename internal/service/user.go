package service

import (
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
)

type UserService struct {
	userRepo   *postgres.UserRepository
	profileRepo *postgres.ProfileRepository
}

func NewUserService(userRepo *postgres.UserRepository, profileRepo *postgres.ProfileRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		profileRepo: profileRepo,
	}
}

func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetProfiles(userID string) ([]*domain.Profile, error) {
	return s.profileRepo.GetByUserID(userID)
}

func (s *UserService) CreateProfile(profile *domain.Profile) error {
	return s.profileRepo.Create(profile)
}

func (s *UserService) UpdateProfile(profile *domain.Profile) error {
	return s.profileRepo.Update(profile)
}

func (s *UserService) DeleteProfile(id string) error {
	return s.profileRepo.Delete(id)
}

package service

import (
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
	"time"
)

type SubscriptionService struct {
	planRepo       *postgres.SubscriptionPlanRepository
	subscriptionRepo *postgres.UserSubscriptionRepository
}

func NewSubscriptionService(planRepo *postgres.SubscriptionPlanRepository, subscriptionRepo *postgres.UserSubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		planRepo:       planRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *SubscriptionService) GetPlans() ([]*domain.SubscriptionPlan, error) {
	return s.planRepo.GetAll()
}

func (s *SubscriptionService) GetPlanByID(id int) (*domain.SubscriptionPlan, error) {
	return s.planRepo.GetByID(id)
}

func (s *SubscriptionService) GetUserActiveSubscription(userID string) (*domain.UserSubscription, error) {
	return s.subscriptionRepo.GetActiveByUserID(userID)
}

func (s *SubscriptionService) GetUserSubscriptions(userID string) ([]*domain.UserSubscription, error) {
	return s.subscriptionRepo.GetByUserID(userID)
}

func (s *SubscriptionService) CreateSubscription(userID string, planID int, paymentProvider, paymentID string) (*domain.UserSubscription, error) {
	plan, err := s.planRepo.GetByID(planID)
	if err != nil || plan == nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.AddDate(0, 0, plan.DurationDays)

	subscription := &domain.UserSubscription{
		UserID:          userID,
		PlanID:          planID,
		Status:          domain.SubscriptionStatusActive,
		StartedAt:       now,
		ExpiresAt:       expiresAt,
		AutoRenew:       true,
		PaymentProvider: paymentProvider,
		PaymentID:       paymentID,
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, err
	}

	subscription.Plan = plan
	return subscription, nil
}

func (s *SubscriptionService) CancelSubscription(subscriptionID string) error {
	return s.subscriptionRepo.Cancel(subscriptionID)
}

func (s *SubscriptionService) IsUserSubscribed(userID string) (bool, error) {
	subscription, err := s.subscriptionRepo.GetActiveByUserID(userID)
	if err != nil {
		return false, err
	}
	return subscription != nil, nil
}

func (s *SubscriptionService) GetUserMaxQuality(userID string) (string, error) {
	subscription, err := s.subscriptionRepo.GetActiveByUserID(userID)
	if err != nil || subscription == nil {
		return "480p", nil // Default quality for non-subscribers
	}

	plan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil || plan == nil {
		return "480p", nil
	}

	return plan.MaxQuality, nil
}

func (s *SubscriptionService) GetUserMaxProfiles(userID string) (int, error) {
	subscription, err := s.subscriptionRepo.GetActiveByUserID(userID)
	if err != nil || subscription == nil {
		return 1, nil // Default 1 profile for non-subscribers
	}

	plan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil || plan == nil {
		return 1, nil
	}

	return plan.MaxProfiles, nil
}

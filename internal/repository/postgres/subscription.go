package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type SubscriptionPlanRepository struct {
	db *sql.DB
}

func NewSubscriptionPlanRepository(db *sql.DB) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{db: db}
}

func (r *SubscriptionPlanRepository) GetAll() ([]*domain.SubscriptionPlan, error) {
	query := `SELECT id, name, display_name, price, currency, duration_days, max_profiles, 
	          max_quality, is_active, features, created_at, updated_at 
	          FROM subscription_plans WHERE is_active = true ORDER BY price`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*domain.SubscriptionPlan
	for rows.Next() {
		plan := &domain.SubscriptionPlan{}
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.DisplayName, &plan.Price, &plan.Currency,
			&plan.DurationDays, &plan.MaxProfiles, &plan.MaxQuality, &plan.IsActive,
			&plan.Features, &plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func (r *SubscriptionPlanRepository) GetByID(id int) (*domain.SubscriptionPlan, error) {
	query := `SELECT id, name, display_name, price, currency, duration_days, max_profiles, 
	          max_quality, is_active, features, created_at, updated_at 
	          FROM subscription_plans WHERE id = $1`
	
	plan := &domain.SubscriptionPlan{}
	err := r.db.QueryRow(query, id).Scan(
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Price, &plan.Currency,
		&plan.DurationDays, &plan.MaxProfiles, &plan.MaxQuality, &plan.IsActive,
		&plan.Features, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *SubscriptionPlanRepository) GetByName(name string) (*domain.SubscriptionPlan, error) {
	query := `SELECT id, name, display_name, price, currency, duration_days, max_profiles, 
	          max_quality, is_active, features, created_at, updated_at 
	          FROM subscription_plans WHERE name = $1 AND is_active = true`
	
	plan := &domain.SubscriptionPlan{}
	err := r.db.QueryRow(query, name).Scan(
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Price, &plan.Currency,
		&plan.DurationDays, &plan.MaxProfiles, &plan.MaxQuality, &plan.IsActive,
		&plan.Features, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return plan, nil
}

type UserSubscriptionRepository struct {
	db *sql.DB
}

func NewUserSubscriptionRepository(db *sql.DB) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{db: db}
}

func (r *UserSubscriptionRepository) Create(subscription *domain.UserSubscription) error {
	// Generate UUID v7 (time-ordered for better index performance)
	subscriptionID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	subscription.ID = subscriptionID
	
	query := `INSERT INTO user_subscriptions (id, user_id, plan_id, status, started_at, expires_at, 
	          auto_renew, payment_provider, payment_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	
	now := time.Now()
	_, err = r.db.Exec(query, subscription.ID, subscription.UserID, subscription.PlanID, subscription.Status,
		subscription.StartedAt, subscription.ExpiresAt, subscription.AutoRenew,
		subscription.PaymentProvider, subscription.PaymentID, now, now)
	if err != nil {
		return err
	}
	
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	return nil
}

func (r *UserSubscriptionRepository) GetActiveByUserID(userID string) (*domain.UserSubscription, error) {
	query := `SELECT id, user_id, plan_id, status, started_at, expires_at, cancelled_at, 
	          auto_renew, payment_provider, payment_id, created_at, updated_at 
	          FROM user_subscriptions 
	          WHERE user_id = $1 AND status = 'active' AND expires_at > CURRENT_TIMESTAMP 
	          ORDER BY expires_at DESC LIMIT 1`
	
	subscription := &domain.UserSubscription{}
	var cancelledAt sql.NullTime
	err := r.db.QueryRow(query, userID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.PlanID, &subscription.Status,
		&subscription.StartedAt, &subscription.ExpiresAt, &cancelledAt,
		&subscription.AutoRenew, &subscription.PaymentProvider, &subscription.PaymentID,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if cancelledAt.Valid {
		subscription.CancelledAt = &cancelledAt.Time
	}
	return subscription, nil
}

func (r *UserSubscriptionRepository) GetByUserID(userID string) ([]*domain.UserSubscription, error) {
	query := `SELECT id, user_id, plan_id, status, started_at, expires_at, cancelled_at, 
	          auto_renew, payment_provider, payment_id, created_at, updated_at 
	          FROM user_subscriptions 
	          WHERE user_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*domain.UserSubscription
	for rows.Next() {
		subscription := &domain.UserSubscription{}
		var cancelledAt sql.NullTime
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.PlanID, &subscription.Status,
			&subscription.StartedAt, &subscription.ExpiresAt, &cancelledAt,
			&subscription.AutoRenew, &subscription.PaymentProvider, &subscription.PaymentID,
			&subscription.CreatedAt, &subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if cancelledAt.Valid {
			subscription.CancelledAt = &cancelledAt.Time
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func (r *UserSubscriptionRepository) Update(subscription *domain.UserSubscription) error {
	query := `UPDATE user_subscriptions SET status = $1, expires_at = $2, cancelled_at = $3, 
	          auto_renew = $4, updated_at = $5 WHERE id = $6`
	subscription.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, subscription.Status, subscription.ExpiresAt, subscription.CancelledAt,
		subscription.AutoRenew, subscription.UpdatedAt, subscription.ID)
	return err
}

func (r *UserSubscriptionRepository) Cancel(subscriptionID string) error {
	query := `UPDATE user_subscriptions SET status = 'cancelled', cancelled_at = CURRENT_TIMESTAMP, 
	          auto_renew = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, subscriptionID)
	return err
}

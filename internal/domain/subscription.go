package domain

import "time"

type SubscriptionPlan struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	DurationDays int      `json:"duration_days"`
	MaxProfiles int       `json:"max_profiles"`
	MaxQuality  string    `json:"max_quality"`
	IsActive    bool      `json:"is_active"`
	Features    string    `json:"features"` // JSON string
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
)

type UserSubscription struct {
	ID            string            `json:"id"`
	UserID        string            `json:"user_id"`
	PlanID        int               `json:"plan_id"`
	Status        SubscriptionStatus `json:"status"`
	StartedAt     time.Time         `json:"started_at"`
	ExpiresAt     time.Time         `json:"expires_at"`
	CancelledAt   *time.Time        `json:"cancelled_at,omitempty"`
	AutoRenew     bool              `json:"auto_renew"`
	PaymentProvider string          `json:"payment_provider,omitempty"`
	PaymentID     string            `json:"payment_id,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	
	// Relations
	Plan *SubscriptionPlan `json:"plan,omitempty"`
}

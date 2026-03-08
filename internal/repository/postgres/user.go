package postgres

import (
	"database/sql"
	"lighthouse/internal/domain"
	"lighthouse/pkg/uuid"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	// Generate UUID v7 (time-ordered for better index performance)
	userID, err := uuid.GenerateV7()
	if err != nil {
		return err
	}
	user.ID = userID
	
	query := `INSERT INTO users (id, email, phone, password, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	
	now := time.Now()
	_, err = r.db.Exec(query, user.ID, user.Email, user.Phone, user.Password, now, now)
	if err != nil {
		return err
	}
	
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, email, phone, password, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	var email, phone sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &email, &phone, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if email.Valid {
		user.Email = &email.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, phone, password, created_at, updated_at FROM users WHERE email = $1`
	user := &domain.User{}
	var emailVal, phone sql.NullString
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &emailVal, &phone, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if emailVal.Valid {
		user.Email = &emailVal.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	return user, nil
}

func (r *UserRepository) GetByPhone(phone string) (*domain.User, error) {
	query := `SELECT id, email, phone, password, created_at, updated_at FROM users WHERE phone = $1`
	user := &domain.User{}
	var email, phoneVal sql.NullString
	err := r.db.QueryRow(query, phone).Scan(
		&user.ID, &email, &phoneVal, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if email.Valid {
		user.Email = &email.String
	}
	if phoneVal.Valid {
		user.Phone = &phoneVal.String
	}
	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET email = $1, phone = $2, password = $3, updated_at = $4 WHERE id = $5`
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.Email, user.Phone, user.Password, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

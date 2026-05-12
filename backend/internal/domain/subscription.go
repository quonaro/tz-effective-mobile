package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidDateFormat    = errors.New("invalid date format, expected MM-YYYY")
)

type Subscription struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price" db:"price"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionInput struct {
	ServiceName string    `json:"service_name" validate:"required,min=1,max=255"`
	Price       int       `json:"price" validate:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,min=0"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type ListSubscriptionsInput struct {
	UserID      uuid.UUID `json:"user_id,omitempty" query:"user_id"`
	ServiceName string    `json:"service_name,omitempty" query:"service_name"`
	Limit       int       `json:"limit,omitempty" query:"limit"`
	Offset      int       `json:"offset,omitempty" query:"offset"`
}

type TotalCostInput struct {
	UserID      uuid.UUID `json:"user_id,omitempty" query:"user_id"`
	ServiceName string    `json:"service_name,omitempty" query:"service_name"`
	StartMonth  string    `json:"start_month" query:"start_month" validate:"required"`
	EndMonth    string    `json:"end_month" query:"end_month" validate:"required"`
}

type TotalCostOutput struct {
	TotalCost int `json:"total_cost"`
}

func ParseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return t, nil
}

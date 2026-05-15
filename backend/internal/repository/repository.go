package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"subscriptions/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input domain.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, input domain.ListSubscriptionsInput) ([]domain.Subscription, int, error)
	TotalCost(ctx context.Context, input domain.TotalCostInput) (int, error)
	GetUniqueServiceNames(ctx context.Context) ([]string, error)
}

type subscriptionRepo struct {
	db     *bun.DB
	logger *slog.Logger
}

func NewSubscriptionRepository(db *bun.DB, logger *slog.Logger) SubscriptionRepository {
	return &subscriptionRepo{db: db, logger: logger}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	_, err := r.db.NewInsert().Model(sub).Exec(ctx)
	if err != nil {
		r.logger.Error("failed to create subscription", slog.String("error", err.Error()))
		return fmt.Errorf("create subscription: %w", err)
	}
	r.logger.Info("subscription created", slog.String("id", sub.ID.String()))
	return nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.NewSelect().Table("subscriptions").Where("id = ?", id).Scan(ctx, &sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
		r.logger.Error("failed to get subscription", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return &sub, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, input domain.UpdateSubscriptionInput) error {
	query := r.db.NewUpdate().Table("subscriptions").Where("id = ?", id)

	if input.ServiceName != nil {
		query = query.Set("service_name = ?", *input.ServiceName)
	}
	if input.Price != nil {
		query = query.Set("price = ?", *input.Price)
	}
	if input.StartDate != nil {
		query = query.Set("start_date = ?", *input.StartDate)
	}
	if input.EndDate != nil {
		query = query.Set("end_date = ?", *input.EndDate)
	}

	res, err := query.Set("updated_at = NOW()").Exec(ctx)
	if err != nil {
		r.logger.Error("failed to update subscription", slog.String("error", err.Error()))
		return fmt.Errorf("update subscription: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrSubscriptionNotFound
	}

	r.logger.Info("subscription updated", slog.String("id", id.String()))
	return nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Table("subscriptions").Where("id = ?", id).Exec(ctx)
	if err != nil {
		r.logger.Error("failed to delete subscription", slog.String("error", err.Error()))
		return fmt.Errorf("delete subscription: %w", err)
	}
	r.logger.Info("subscription deleted", slog.String("id", id.String()))
	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, input domain.ListSubscriptionsInput) ([]domain.Subscription, int, error) {
	var subs []domain.Subscription
	var total int

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		countQuery := tx.NewSelect().Table("subscriptions")
		if input.UserID != uuid.Nil {
			countQuery = countQuery.Where("user_id = ?", input.UserID)
		}
		if input.ServiceName != "" {
			countQuery = countQuery.Where("service_name ILIKE ?", "%"+input.ServiceName+"%")
		}
		if err := countQuery.ColumnExpr("COUNT(*)").Scan(ctx, &total); err != nil {
			return fmt.Errorf("count subscriptions: %w", err)
		}

		selectQuery := tx.NewSelect().Table("subscriptions")
		if input.UserID != uuid.Nil {
			selectQuery = selectQuery.Where("user_id = ?", input.UserID)
		}
		if input.ServiceName != "" {
			selectQuery = selectQuery.Where("service_name ILIKE ?", "%"+input.ServiceName+"%")
		}

		sortBy := "created_at"
		if input.SortBy == "price" {
			sortBy = "price"
		}
		sortOrder := "DESC"
		if input.SortOrder == "asc" {
			sortOrder = "ASC"
		}

		limit := input.Limit
		if limit <= 0 {
			limit = 20
		}
		offset := input.Offset
		if offset < 0 {
			offset = 0
		}

		if err := selectQuery.
			OrderExpr("? ?", bun.Ident(sortBy), bun.Safe(sortOrder)).
			Limit(limit).
			Offset(offset).
			Scan(ctx, &subs); err != nil {
			return fmt.Errorf("list subscriptions: %w", err)
		}
		return nil
	})
	if err != nil {
		r.logger.Error("failed to list subscriptions", slog.String("error", err.Error()))
		return nil, 0, err
	}

	return subs, total, nil
}

func (r *subscriptionRepo) TotalCost(ctx context.Context, input domain.TotalCostInput) (int, error) {
	start, err := domain.ParseMonthYear(input.StartMonth)
	if err != nil {
		return 0, err
	}
	end, err := domain.ParseMonthYear(input.EndMonth)
	if err != nil {
		return 0, err
	}

	reqEnd := end.AddDate(0, 1, -1)

	query := r.db.NewSelect().Table("subscriptions").
		Where("start_date <= ?", reqEnd).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("end_date IS NULL").WhereOr("end_date >= ?", start)
		})

	if input.UserID != uuid.Nil {
		query = query.Where("user_id = ?", input.UserID)
	}
	if input.ServiceName != "" {
		query = query.Where("service_name ILIKE ?", "%"+input.ServiceName+"%")
	}

	var subs []domain.Subscription
	if err := query.Scan(ctx, &subs); err != nil {
		r.logger.Error("failed to calculate total cost", slog.String("error", err.Error()))
		return 0, fmt.Errorf("total cost: %w", err)
	}

	total := calculateTotalCost(subs, start, reqEnd)
	return total, nil
}

func (r *subscriptionRepo) GetUniqueServiceNames(ctx context.Context) ([]string, error) {
	var services []string
	err := r.db.NewSelect().Table("subscriptions").
		ColumnExpr("DISTINCT service_name").
		OrderExpr("service_name").
		Scan(ctx, &services)
	if err != nil {
		r.logger.Error("failed to get unique service names", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get unique service names: %w", err)
	}
	return services, nil
}

func calculateTotalCost(subs []domain.Subscription, reqStart, reqEnd time.Time) int {
	total := 0
	for _, sub := range subs {
		overlapMonths := overlappingMonths(sub.StartDate, sub.EndDate, reqStart, reqEnd)
		total += sub.Price * overlapMonths
	}
	return total
}

func overlappingMonths(subStart time.Time, subEnd *time.Time, reqStart, reqEnd time.Time) int {
	var subEffectiveEnd time.Time
	if subEnd != nil && !subEnd.IsZero() {
		subEffectiveEnd = subEnd.AddDate(0, 1, -1)
	} else {
		subEffectiveEnd = reqEnd
	}

	overlapStart := subStart
	if reqStart.After(overlapStart) {
		overlapStart = reqStart
	}

	overlapEnd := subEffectiveEnd
	if reqEnd.Before(overlapEnd) {
		overlapEnd = reqEnd
	}

	if overlapStart.After(overlapEnd) {
		return 0
	}

	months := (overlapEnd.Year()-overlapStart.Year())*12 + int(overlapEnd.Month()-overlapStart.Month()) + 1
	if months < 0 {
		return 0
	}
	return months
}

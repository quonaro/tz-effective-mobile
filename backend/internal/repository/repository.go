package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"subscriptions/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewSubscriptionRepository(db *pgxpool.Pool, logger *slog.Logger) SubscriptionRepository {
	return &subscriptionRepo{db: db, logger: logger}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.CreatedAt, sub.UpdatedAt)
	if err != nil {
		r.logger.Error("failed to create subscription", slog.String("error", err.Error()))
		return fmt.Errorf("create subscription: %w", err)
	}
	r.logger.Info("subscription created", slog.String("id", sub.ID.String()))
	return nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var sub domain.Subscription
	var endDate *time.Time
	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		r.logger.Error("failed to get subscription", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	sub.EndDate = endDate
	return &sub, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, input domain.UpdateSubscriptionInput) error {
	query := `
		UPDATE subscriptions
		SET service_name = COALESCE($2, service_name),
		    price = COALESCE($3, price),
		    start_date = COALESCE($4, start_date),
		    end_date = COALESCE($5, end_date),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, input.ServiceName, input.Price, input.StartDate, input.EndDate)
	if err != nil {
		r.logger.Error("failed to update subscription", slog.String("error", err.Error()))
		return fmt.Errorf("update subscription: %w", err)
	}
	r.logger.Info("subscription updated", slog.String("id", id.String()))
	return nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete subscription", slog.String("error", err.Error()))
		return fmt.Errorf("delete subscription: %w", err)
	}
	r.logger.Info("subscription deleted", slog.String("id", id.String()))
	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, input domain.ListSubscriptionsInput) ([]domain.Subscription, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if input.UserID != uuid.Nil {
		where += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, input.UserID)
		argIdx++
	}
	if input.ServiceName != "" {
		where += fmt.Sprintf(" AND service_name ILIKE $%d", argIdx)
		args = append(args, "%"+input.ServiceName+"%")
		argIdx++
	}

	countQuery := "SELECT COUNT(*) FROM subscriptions " + where
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error("failed to count subscriptions", slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("count subscriptions: %w", err)
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(
		`SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		 FROM subscriptions %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to list subscriptions", slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		var endDate *time.Time
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			r.logger.Error("failed to scan subscription", slog.String("error", err.Error()))
			return nil, 0, fmt.Errorf("scan subscription: %w", err)
		}
		sub.EndDate = endDate
		subs = append(subs, sub)
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

	where := "WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)"
	args := []interface{}{end.AddDate(0, 1, -1), start}
	argIdx := 3

	if input.UserID != uuid.Nil {
		where += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, input.UserID)
		argIdx++
	}
	if input.ServiceName != "" {
		where += fmt.Sprintf(" AND service_name ILIKE $%d", argIdx)
		args = append(args, "%"+input.ServiceName+"%")
		argIdx++
	}

	query := "SELECT COALESCE(SUM(price), 0) FROM subscriptions " + where
	var total int
	err = r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		r.logger.Error("failed to calculate total cost", slog.String("error", err.Error()))
		return 0, fmt.Errorf("total cost: %w", err)
	}

	return total, nil
}

func (r *subscriptionRepo) GetUniqueServiceNames(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT service_name FROM subscriptions ORDER BY service_name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("failed to get unique service names", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get unique service names: %w", err)
	}
	defer rows.Close()

	var services []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			r.logger.Error("failed to scan service name", slog.String("error", err.Error()))
			return nil, fmt.Errorf("scan service name: %w", err)
		}
		services = append(services, name)
	}

	return services, nil
}

package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"subscriptions/internal/domain"
	"subscriptions/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, input domain.CreateSubscriptionInput) (*domain.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, startDateStr, endDateStr *string, input domain.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, input domain.ListSubscriptionsInput) ([]domain.Subscription, int, error)
	TotalCost(ctx context.Context, input domain.TotalCostInput) (int, error)
	GetUniqueServiceNames(ctx context.Context) ([]string, error)
}

type subscriptionService struct {
	repo   repository.SubscriptionRepository
	logger *slog.Logger
}

func NewSubscriptionService(repo repository.SubscriptionRepository, logger *slog.Logger) SubscriptionService {
	return &subscriptionService{repo: repo, logger: logger}
}

func (s *subscriptionService) Create(ctx context.Context, input domain.CreateSubscriptionInput) (*domain.Subscription, error) {
	startDate, err := domain.ParseMonthYear(input.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if input.EndDate != nil {
		parsed, err := domain.ParseMonthYear(*input.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsed
	}

	now := time.Now().UTC()
	sub := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		s.logger.Error("failed to create subscription", slog.String("error", err.Error()))
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return sub, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, startDateStr, endDateStr *string, input domain.UpdateSubscriptionInput) error {
	if startDateStr != nil {
		parsed, err := domain.ParseMonthYear(*startDateStr)
		if err != nil {
			return err
		}
		input.StartDate = &parsed
	}
	if endDateStr != nil {
		parsed, err := domain.ParseMonthYear(*endDateStr)
		if err != nil {
			return err
		}
		input.EndDate = &parsed
	}

	if err := s.repo.Update(ctx, id, input); err != nil {
		s.logger.Error("failed to update subscription", slog.String("error", err.Error()))
		return fmt.Errorf("update subscription: %w", err)
	}
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete subscription", slog.String("error", err.Error()))
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}

func (s *subscriptionService) List(ctx context.Context, input domain.ListSubscriptionsInput) ([]domain.Subscription, int, error) {
	subs, total, err := s.repo.List(ctx, input)
	if err != nil {
		s.logger.Error("failed to list subscriptions", slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("list subscriptions: %w", err)
	}
	return subs, total, nil
}

func (s *subscriptionService) TotalCost(ctx context.Context, input domain.TotalCostInput) (int, error) {
	_, err := domain.ParseMonthYear(input.StartMonth)
	if err != nil {
		return 0, err
	}
	_, err = domain.ParseMonthYear(input.EndMonth)
	if err != nil {
		return 0, err
	}

	total, err := s.repo.TotalCost(ctx, input)
	if err != nil {
		s.logger.Error("failed to calculate total cost", slog.String("error", err.Error()))
		return 0, fmt.Errorf("total cost: %w", err)
	}
	return total, nil
}

func (s *subscriptionService) GetUniqueServiceNames(ctx context.Context) ([]string, error) {
	services, err := s.repo.GetUniqueServiceNames(ctx)
	if err != nil {
		s.logger.Error("failed to get unique service names", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get unique service names: %w", err)
	}
	return services, nil
}

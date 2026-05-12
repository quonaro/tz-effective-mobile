package repository

import (
	"context"

	"subscriptions/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(ctx context.Context, id uuid.UUID, input domain.UpdateSubscriptionInput) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) List(ctx context.Context, input domain.ListSubscriptionsInput) ([]domain.Subscription, int, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]domain.Subscription), args.Get(1).(int), args.Error(2)
}

func (m *MockSubscriptionRepository) TotalCost(ctx context.Context, input domain.TotalCostInput) (int, error) {
	args := m.Called(ctx, input)
	return args.Int(0), args.Error(1)
}

func (m *MockSubscriptionRepository) GetUniqueServiceNames(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

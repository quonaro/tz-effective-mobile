package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"subscriptions/internal/domain"
	"subscriptions/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestService(mockRepo *repository.MockSubscriptionRepository) SubscriptionService {
	logger := slog.New(slog.DiscardHandler)
	return NewSubscriptionService(mockRepo, logger)
}

func TestSubscriptionService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		userID := uuid.New()
		input := domain.CreateSubscriptionInput{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      userID,
			StartDate:   "07-2025",
		}

		mockRepo.On("Create", ctx, mock.Anything).Return(nil)

		sub, err := service.Create(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, sub)
		assert.Equal(t, input.ServiceName, sub.ServiceName)
		assert.Equal(t, input.Price, sub.Price)
		assert.Equal(t, input.UserID, sub.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid date format", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.CreateSubscriptionInput{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      uuid.New(),
			StartDate:   "invalid-date",
		}

		sub, err := service.Create(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, sub)
		assert.Equal(t, domain.ErrInvalidDateFormat, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.CreateSubscriptionInput{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      uuid.New(),
			StartDate:   "07-2025",
		}

		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		sub, err := service.Create(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, sub)
		assert.Contains(t, err.Error(), "create subscription")
		mockRepo.AssertExpectations(t)
	})

	t.Run("with end date", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		userID := uuid.New()
		input := domain.CreateSubscriptionInput{
			ServiceName: "Netflix",
			Price:       500,
			UserID:      userID,
			StartDate:   "01-2025",
			EndDate:     strPtr("12-2025"),
		}

		mockRepo.On("Create", ctx, mock.Anything).Return(nil)

		sub, err := service.Create(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, sub)
		assert.NotNil(t, sub.EndDate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid end date format", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.CreateSubscriptionInput{
			ServiceName: "Netflix",
			Price:       500,
			UserID:      uuid.New(),
			StartDate:   "01-2025",
			EndDate:     strPtr("invalid"),
		}

		sub, err := service.Create(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, sub)
		assert.Equal(t, domain.ErrInvalidDateFormat, err)
	})
}

func TestSubscriptionService_GetByID(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()
		expectedSub := &domain.Subscription{
			ID:          id,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      uuid.New(),
			StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo.On("GetByID", ctx, id).Return(expectedSub, nil)

		sub, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedSub, sub)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()

		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrSubscriptionNotFound)

		sub, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, sub)
		assert.Equal(t, domain.ErrSubscriptionNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionService_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()
		newName := "Netflix"
		newPrice := 500
		input := domain.UpdateSubscriptionInput{
			ServiceName: &newName,
			Price:       &newPrice,
		}

		mockRepo.On("Update", ctx, id, mock.Anything).Return(nil)

		err := service.Update(ctx, id, input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid start date format", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()
		invalidDate := "invalid"
		input := domain.UpdateSubscriptionInput{
			StartDate: &invalidDate,
		}

		err := service.Update(ctx, id, input)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidDateFormat, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()
		input := domain.UpdateSubscriptionInput{}

		mockRepo.On("Update", ctx, id, mock.Anything).Return(errors.New("db error"))

		err := service.Update(ctx, id, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update subscription")
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionService_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()

		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		id := uuid.New()

		mockRepo.On("Delete", ctx, id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete subscription")
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionService_List(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		userID := uuid.New()
		input := domain.ListSubscriptionsInput{
			UserID: userID,
			Limit:  10,
		}
		expectedSubs := []domain.Subscription{
			{ID: uuid.New(), ServiceName: "Yandex Plus"},
			{ID: uuid.New(), ServiceName: "Netflix"},
		}

		mockRepo.On("List", ctx, mock.Anything).Return(expectedSubs, 2, nil)

		subs, total, err := service.List(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedSubs, subs)
		assert.Equal(t, 2, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.ListSubscriptionsInput{}

		mockRepo.On("List", ctx, mock.Anything).Return([]domain.Subscription{}, 0, errors.New("db error"))

		subs, total, err := service.List(ctx, input)

		assert.Error(t, err)
		assert.Empty(t, subs)
		assert.Equal(t, 0, total)
		assert.Contains(t, err.Error(), "list subscriptions")
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionService_TotalCost(t *testing.T) {
	t.Run("successful calculation", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		userID := uuid.New()
		input := domain.TotalCostInput{
			UserID:     userID,
			StartMonth: "01-2025",
			EndMonth:   "12-2025",
		}

		mockRepo.On("TotalCost", ctx, mock.Anything).Return(900, nil)

		total, err := service.TotalCost(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, 900, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid start month", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.TotalCostInput{
			StartMonth: "invalid",
			EndMonth:   "12-2025",
		}

		total, err := service.TotalCost(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, 0, total)
		assert.Equal(t, domain.ErrInvalidDateFormat, err)
	})

	t.Run("invalid end month", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.TotalCostInput{
			StartMonth: "01-2025",
			EndMonth:   "invalid",
		}

		total, err := service.TotalCost(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, 0, total)
		assert.Equal(t, domain.ErrInvalidDateFormat, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSubscriptionRepository)
		service := setupTestService(mockRepo)
		ctx := context.Background()
		input := domain.TotalCostInput{
			StartMonth: "01-2025",
			EndMonth:   "12-2025",
		}

		mockRepo.On("TotalCost", ctx, mock.Anything).Return(0, errors.New("db error"))

		total, err := service.TotalCost(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, 0, total)
		assert.Contains(t, err.Error(), "total cost")
		mockRepo.AssertExpectations(t)
	})
}

func strPtr(s string) *string {
	return &s
}

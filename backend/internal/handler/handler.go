package handler

import (
	"context"
	"net/http"
	"os"

	"subscriptions/internal/domain"
	"subscriptions/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewRouter(subService service.SubscriptionService) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	api := humachi.New(router, huma.DefaultConfig("Subscriptions API", "1.0.0"))

	// Create
	huma.Register(api, huma.Operation{
		OperationID: "create-subscription",
		Method:      http.MethodPost,
		Path:        "/subscriptions",
	}, func(ctx context.Context, input *struct {
		Body domain.CreateSubscriptionInput
	}) (*struct {
		Body *domain.Subscription
	}, error) {
		sub, err := subService.Create(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &struct{ Body *domain.Subscription }{Body: sub}, nil
	})

	// Get
	huma.Register(api, huma.Operation{
		OperationID: "get-subscription",
		Method:      http.MethodGet,
		Path:        "/subscriptions/{id}",
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" validate:"required,uuid"`
	}) (*struct {
		Body *domain.Subscription
	}, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id format", err)
		}
		sub, err := subService.GetByID(ctx, id)
		if err != nil {
			if err == domain.ErrSubscriptionNotFound {
				return nil, huma.Error404NotFound("subscription not found", err)
			}
			return nil, err
		}
		return &struct{ Body *domain.Subscription }{Body: sub}, nil
	})

	// Update
	type UpdateSubscriptionRequest struct {
		ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
		Price       *int    `json:"price,omitempty" validate:"omitempty,min=0"`
		StartDate   *string `json:"start_date,omitempty"`
		EndDate     *string `json:"end_date,omitempty"`
	}

	huma.Register(api, huma.Operation{
		OperationID: "update-subscription",
		Method:      http.MethodPut,
		Path:        "/subscriptions/{id}",
	}, func(ctx context.Context, input *struct {
		ID   string `path:"id" validate:"required,uuid"`
		Body UpdateSubscriptionRequest
	}) (*struct{}, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id format", err)
		}
		updateInput := domain.UpdateSubscriptionInput{
			ServiceName: input.Body.ServiceName,
			Price:       input.Body.Price,
		}
		if err := subService.Update(ctx, id, input.Body.StartDate, input.Body.EndDate, updateInput); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})

	// Delete
	huma.Register(api, huma.Operation{
		OperationID: "delete-subscription",
		Method:      http.MethodDelete,
		Path:        "/subscriptions/{id}",
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" validate:"required,uuid"`
	}) (*struct{}, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id format", err)
		}
		if err := subService.Delete(ctx, id); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})

	// List
	huma.Register(api, huma.Operation{
		OperationID: "list-subscriptions",
		Method:      http.MethodGet,
		Path:        "/subscriptions",
	}, func(ctx context.Context, input *domain.ListSubscriptionsInput) (*struct {
		Body struct {
			Data  []domain.Subscription `json:"data"`
			Total int                   `json:"total"`
		}
	}, error) {
		subs, total, err := subService.List(ctx, *input)
		if err != nil {
			return nil, err
		}
		if subs == nil {
			subs = []domain.Subscription{}
		}
		return &struct {
			Body struct {
				Data  []domain.Subscription `json:"data"`
				Total int                   `json:"total"`
			}
		}{Body: struct {
			Data  []domain.Subscription `json:"data"`
			Total int                   `json:"total"`
		}{Data: subs, Total: total}}, nil
	})

	// Total Cost
	huma.Register(api, huma.Operation{
		OperationID: "total-cost",
		Method:      http.MethodGet,
		Path:        "/subscriptions/total",
	}, func(ctx context.Context, input *domain.TotalCostInput) (*struct {
		Body domain.TotalCostOutput
	}, error) {
		total, err := subService.TotalCost(ctx, *input)
		if err != nil {
			return nil, err
		}
		return &struct{ Body domain.TotalCostOutput }{Body: domain.TotalCostOutput{TotalCost: total}}, nil
	})

	// Get Unique Service Names
	huma.Register(api, huma.Operation{
		OperationID: "get-services",
		Method:      http.MethodGet,
		Path:        "/subscriptions/services",
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Services []string `json:"services"`
		}
	}, error) {
		services, err := subService.GetUniqueServiceNames(ctx)
		if err != nil {
			return nil, err
		}
		return &struct {
			Body struct {
				Services []string `json:"services"`
			}
		}{Body: struct {
			Services []string `json:"services"`
		}{Services: services}}, nil
	})

	// Get TZ (Technical Specification) - reads docs/ТЗ.md directly
	huma.Register(api, huma.Operation{
		OperationID: "get-tz",
		Method:      http.MethodGet,
		Path:        "/tz",
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body string `contentType:"text/markdown"`
	}, error) {
		data, err := os.ReadFile("docs/ТЗ.md")
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to read TZ file", err)
		}
		return &struct {
			Body string `contentType:"text/markdown"`
		}{Body: string(data)}, nil
	})

	return router
}

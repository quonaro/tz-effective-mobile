package repository

import (
	"testing"
	"time"

	"subscriptions/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestOverlappingMonths(t *testing.T) {
	tests := []struct {
		name        string
		subStart    time.Time
		subEnd      *time.Time
		reqStart    time.Time
		reqEnd      time.Time
		expected    int
	}{
		{
			name:     "exact match single month",
			subStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "subscription spans full request range",
			subStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
		{
			name:     "partial overlap at start",
			subStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
		{
			name:     "partial overlap at end",
			subStart: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
		{
			name:     "ongoing subscription",
			subStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   nil,
			reqStart: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
		{
			name:     "no overlap before",
			subStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "no overlap after",
			subStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "multi-year subscription",
			subStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
		{
			name:     "single month request with multi-month sub",
			subStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)),
			reqStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			reqEnd:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := overlappingMonths(tt.subStart, tt.subEnd, tt.reqStart, tt.reqEnd)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateTotalCost(t *testing.T) {
	subs := []domain.Subscription{
		{
			ServiceName: "Netflix",
			Price:       500,
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     ptr(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			ServiceName: "Spotify",
			Price:       300,
			StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     ptr(time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	reqStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	reqEnd := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)

	// Netflix: overlaps Feb, Mar = 2 months * 500 = 1000
	// Spotify: overlaps Feb, Mar = 2 months * 300 = 600
	// Total = 1600
	total := calculateTotalCost(subs, reqStart, reqEnd)
	assert.Equal(t, 1600, total)
}

func ptr(t time.Time) *time.Time {
	return &t
}

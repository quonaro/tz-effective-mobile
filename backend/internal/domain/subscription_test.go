package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseMonthYear(t *testing.T) {
	t.Run("valid formats", func(t *testing.T) {
		tests := []struct {
			input    string
			expected time.Time
		}{
			{"01-2025", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
			{"07-2025", time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)},
			{"12-2025", time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)},
		}

		for _, tt := range tests {
			result, err := ParseMonthYear(tt.input)
			assert.NoError(t, err, "should parse valid date: %s", tt.input)
			assert.Equal(t, tt.expected, result, "parsed date should match for: %s", tt.input)
		}
	})

	t.Run("invalid formats", func(t *testing.T) {
		invalidInputs := []string{
			"2025-01",
			"01/2025",
			"01-25",
			"January-2025",
			"invalid",
			"",
			"13-2025",
			"00-2025",
		}

		for _, input := range invalidInputs {
			result, err := ParseMonthYear(input)
			assert.Error(t, err, "should return error for invalid format: %s", input)
			assert.Equal(t, ErrInvalidDateFormat, err, "should return ErrInvalidDateFormat for: %s", input)
			assert.True(t, result.IsZero(), "result should be zero time for invalid: %s", input)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		// Leap year February
		result, err := ParseMonthYear("02-2024")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), result)

		// Non-leap year February
		result, err = ParseMonthYear("02-2025")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), result)
	})
}

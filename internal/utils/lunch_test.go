package utils_test

import (
	"testing"

	"github.com/Te8va/Gofermarch/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestIsValidLuhn(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid Luhn number",
			input:    "79927398713",
			expected: true,
		},
		{
			name:     "Invalid Luhn number",
			input:    "79927398714",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.IsValidLuhn(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

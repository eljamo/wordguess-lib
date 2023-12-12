package rng

import (
	"testing"
)

type MockRNGService struct{}

func (s *MockRNGService) GenerateWithMax(max int) (int, error) {
	return 1, nil
}

type MockEvenRNGService struct{}

func (s *MockEvenRNGService) GenerateWithMax(max int) (int, error) {
	return 2, nil
}

func TestRNGGenerateWithMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		max       int
		expectErr bool
	}{
		{"ValidMax", 100, false},
		{"NegativeMax", -1, true},
	}

	rngSvc := NewRNGService()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			generated, err := rngSvc.GenerateWithMax(tc.max)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error for max = %v, but got none", tc.max)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for max = %v: %v", tc.max, err)
				}
				if generated < 0 || generated >= tc.max {
					t.Errorf("Generated number is out of bounds for max = %v: got %v", tc.max, generated)
				}
			}
		})
	}
}

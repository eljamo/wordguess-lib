package rng

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// RNGService defines an interface for random number generation.
// It provides a methods for generating random integers
type RNGService interface {
	// Generates a random integer up to the specified maximum value.
	GenerateWithMax(max int) (int, error)
}

// DefaultRNGService is a struct implementing the RNGService interface.
type DefaultRNGService struct{}

// Creates a new instance of DefaultRNGService.
func NewRNGService() *DefaultRNGService {
	return &DefaultRNGService{}
}

// Generates a random integer up to the specified maximum value.
func (s *DefaultRNGService) GenerateWithMax(max int) (int, error) {
	if max < 1 {
		return 0, errors.New("rng max cannot be less than 1")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

package service

// Defines the interface for a service that gets today's word
type WordService interface {
	// GetWord returns a word and returns an error if there is any.
	GetWord() (string, error)
}

package repository

import (
	"github.com/eljamo/libwordle/domain"
)

type RecentWordLogRepository interface {
	// Returns all records in the table
	FindAll() ([]domain.RecentWordLog, error)
	// Returns a single record by date
	FindByDate(date string) (*domain.RecentWordLog, error)
	// Inserts a new record into the table
	InsertWord(date string, word string) error
	// Deletes a record by date
	DeleteByDate(date string) error
}

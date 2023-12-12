package database

import "database/sql"

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

type Scanner interface {
	Scan(dest ...any) error
}

type Rows interface {
	Scanner
	Close() error
	Err() error
	Next() bool
}

func ScanRows(r Rows, scanFunc func(row Scanner) error) error {
	var closeErr error

	defer func() {
		if err := r.Close(); err != nil {
			closeErr = err
		}
	}()

	var scanErr error
	for r.Next() {
		err := scanFunc(r)
		if err != nil {
			scanErr = err
			break
		}
	}
	if r.Err() != nil {
		return r.Err()
	}
	if scanErr != nil {
		return scanErr
	}

	return closeErr
}

package sqlite

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Migrator struct {
	FS   embed.FS
	Path string
	Run  bool
}

type DB struct {
	*sqlx.DB
}

func New(dsn string, migrator *Migrator) (*DB, error) {
	db, err := sqlx.Connect("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, err
	}

	if migrator.Run {
		iofsDriver, err := iofs.New(migrator.FS, migrator.Path)
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, fmt.Sprintf("sqlite3://%s", dsn))
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}

	return &DB{db}, nil
}

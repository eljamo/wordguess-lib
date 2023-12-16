package sqlite

import (
	_ "embed"

	"github.com/eljamo/libwordle/domain"
	"github.com/eljamo/libwordle/internal/database"
	"github.com/eljamo/libwordle/internal/sqlite"
	"github.com/midir99/sqload"
)

func ScanRecentWordLog(row database.Scanner) (*domain.RecentWordLog, error) {
	var e domain.RecentWordLog
	err := row.Scan(&e.Date, &e.Word)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

//go:embed recent_word_log.sql
var RecentWordLogSQL string
var RecentWordLogQuery = sqload.MustLoadFromString[struct {
	FindAll      string `query:"FindAll"`
	FindByDate   string `query:"FindByDate"`
	InsertWord   string `query:"InsertWord"`
	DeleteByDate string `query:"DeleteByDate"`
}](RecentWordLogSQL)

type RecentWordLogRepository struct {
	db *sqlite.DB
}

func NewRecentWordLogRepository(db *sqlite.DB) *RecentWordLogRepository {
	return &RecentWordLogRepository{
		db,
	}
}

func (r *RecentWordLogRepository) FindAll() ([]domain.RecentWordLog, error) {
	rows, err := r.db.Query(RecentWordLogQuery.FindAll)
	if err != nil {
		return nil, err
	}

	var els []domain.RecentWordLog
	scanErr := database.ScanRows(rows, func(row database.Scanner) error {
		el, err := ScanRecentWordLog(rows)
		if err != nil {
			return err
		}

		els = append(els, *el)
		return nil
	})
	if scanErr != nil {
		return nil, scanErr
	}

	return els, nil
}

func (r *RecentWordLogRepository) FindByDate(date string) (*domain.RecentWordLog, error) {
	row := r.db.QueryRow(RecentWordLogQuery.FindByDate, date)

	return ScanRecentWordLog(row)
}

func (r *RecentWordLogRepository) InsertWord(date string, word string) error {
	_, err := r.db.Exec(RecentWordLogQuery.InsertWord, date, word)

	return err
}

func (r *RecentWordLogRepository) DeleteByDate(date string) error {
	_, err := r.db.Exec(RecentWordLogQuery.DeleteByDate, date)

	return err
}

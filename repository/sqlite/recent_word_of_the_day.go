package sqlite

import (
	_ "embed"

	"github.com/eljamo/libwordle/domain"
	"github.com/eljamo/libwordle/internal/database"
	"github.com/eljamo/libwordle/internal/sqlite"
	"github.com/midir99/sqload"
)

func ScanRecentWordOfTheDay(row database.Scanner) (*domain.RecentWordOfTheDay, error) {
	var e domain.RecentWordOfTheDay
	err := row.Scan(&e.Date, &e.Word)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

//go:embed recent_word_of_the_day.sql
var recentWordOfTheDaySQL string
var recentWordOfTheDayQuery = sqload.MustLoadFromString[struct {
	FindAll      string `query:"FindAll"`
	FindByDate   string `query:"FindByDate"`
	InsertWord   string `query:"InsertWord"`
	DeleteByDate string `query:"DeleteByDate"`
}](recentWordOfTheDaySQL)

type RecentWordOfTheDayRepository struct {
	db *sqlite.DB
}

func NewRecentWordOfTheDayRepository(db *sqlite.DB) *RecentWordOfTheDayRepository {
	return &RecentWordOfTheDayRepository{
		db: db,
	}
}

func (r *RecentWordOfTheDayRepository) FindAll() ([]domain.RecentWordOfTheDay, error) {
	rows, err := r.db.Query(recentWordOfTheDayQuery.FindAll)
	if err != nil {
		return nil, err
	}

	var els []domain.RecentWordOfTheDay
	scanErr := database.ScanRows(rows, func(row database.Scanner) error {
		el, err := ScanRecentWordOfTheDay(rows)
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

func (r *RecentWordOfTheDayRepository) FindByDate(date string) (*domain.RecentWordOfTheDay, error) {
	row := r.db.QueryRow(recentWordOfTheDayQuery.FindByDate, date)

	return ScanRecentWordOfTheDay(row)
}

func (r *RecentWordOfTheDayRepository) InsertWord(date string, word string) error {
	_, err := r.db.Exec(recentWordOfTheDayQuery.InsertWord, date, word)

	return err
}

func (r *RecentWordOfTheDayRepository) DeleteByDate(date string) error {
	_, err := r.db.Exec(recentWordOfTheDayQuery.DeleteByDate, date)

	return err
}

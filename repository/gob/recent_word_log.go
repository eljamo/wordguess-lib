package gob

import (
	"fmt"

	"github.com/eljamo/libwordle/config"
	"github.com/eljamo/libwordle/domain"
	"github.com/eljamo/libwordle/internal/gobfile"
)

type RecentWordLogGobRepository struct {
	cfg *config.Settings
}

const gobFileName string = "recent_word_log"

func NewRecentWordLogRepository(cfg *config.Settings) *RecentWordLogGobRepository {
	return &RecentWordLogGobRepository{
		cfg,
	}
}

func (r *RecentWordLogGobRepository) FindAll() ([]domain.RecentWordLog, error) {
	store, err := gobfile.ReadData[[]domain.RecentWordLog](r.cfg.AppName, gobFileName)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (r *RecentWordLogGobRepository) FindByDate(date string) (*domain.RecentWordLog, error) {
	store, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, el := range store {
		if el.Date == date {
			return &el, nil
		}
	}

	return nil, fmt.Errorf("no record found with date %s", date)
}

func (r *RecentWordLogGobRepository) InsertWord(date string, word string) error {
	store, err := r.FindAll()
	if err != nil {
		return err
	}

	store = append(store, domain.RecentWordLog{
		Date: date,
		Word: word,
	})

	return gobfile.WriteData[[]domain.RecentWordLog](r.cfg.AppName, gobFileName, store)
}

func (r *RecentWordLogGobRepository) DeleteByDate(date string) error {
	store, err := r.FindAll()
	if err != nil {
		return err
	}

	for i, el := range store {
		if el.Date == date {
			store = append(store[:i], store[i+1:]...)
			break
		}
	}

	return gobfile.WriteData[[]domain.RecentWordLog](r.cfg.AppName, gobFileName, store)
}

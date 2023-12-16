package service

import (
	"fmt"
	"time"

	"github.com/eljamo/libwordle/asset"
	"github.com/eljamo/libwordle/config"
	"github.com/eljamo/libwordle/domain"
	"github.com/eljamo/libwordle/internal/rng"
	"github.com/eljamo/libwordle/internal/validator"
	"github.com/eljamo/libwordle/repository"
)

// A default implementation of the WordService interface. This can act as a
// example or you can use it as is. Only suitable for development, testing,
// offline applications, or single server applications. Multi-server use
// would need to adjustments to word selection.
type DefaultWordService struct {
	cfg      *config.Settings
	repo     repository.RecentWordLogRepository
	rngSvc   rng.RNGService
	wordList []string
}

const isoDateFmt = "2006-01-02"
const recordMax = 364

// Creates a new instance of DefaultWordService with a custom repository.
func NewDefaultWordService(
	cfg *config.Settings,
	repo repository.RecentWordLogRepository,
	rngSvc rng.RNGService,
) (*DefaultWordService, error) {
	wordList, err := getWordList(cfg.WordList, cfg.WordLength)
	if err != nil {
		return nil, fmt.Errorf("failed to get word list: %v", err)
	}

	return &DefaultWordService{
		cfg,
		repo,
		rngSvc,
		wordList,
	}, nil
}

func getWordList(wordList string, wordLength int) ([]string, error) {
	wl, err := asset.GetFilteredWordList(wordList, wordLength)
	if err != nil {
		return nil, err
	}

	if len(wl) == 0 {
		return nil, fmt.Errorf("no words found in word list %s with a word_length of %d", wordList, wordLength)
	}

	return wl, nil
}

func (s *DefaultWordService) extractWordsFromRecentRecords(slice []domain.RecentWordLog) []string {
	var words []string
	for _, s := range slice {
		words = append(words, s.Word)
	}

	return words
}

func (s *DefaultWordService) findNewWord(pws []string) (string, error) {
	var word string
	foundNewWord := false
	for !foundNewWord {
		idx, err := s.rngSvc.GenerateWithMax(len(s.wordList))
		if err != nil {
			return "", err
		}

		word = s.wordList[idx]

		if !validator.IsElementInSlice(pws, word) {
			foundNewWord = true
		}
	}

	return word, nil
}

func (s *DefaultWordService) getNewWord(recentWords []domain.RecentWordLog) (string, error) {
	pws := s.extractWordsFromRecentRecords(recentWords)
	if validator.UnorderedSlicesAreEqual(s.wordList, pws) {
		return "", fmt.Errorf("all words in word list %s with a word_length of %d have been used", s.cfg.WordList, s.cfg.WordLength)
	}

	return s.findNewWord(pws)
}

func findOldestRecord(records []domain.RecentWordLog) (string, error) {
	if len(records) == 0 {
		return "", fmt.Errorf("no records found")
	}

	var oldestDate time.Time
	var initialized bool
	for _, record := range records {
		parsedDate, err := time.Parse(isoDateFmt, record.Date)
		if err != nil {
			return "", fmt.Errorf("invalid date format: %v", err)
		}

		if !initialized || parsedDate.Before(oldestDate) {
			oldestDate = parsedDate
			initialized = true
		}
	}

	return oldestDate.Format(isoDateFmt), nil
}

func (s *DefaultWordService) deleteOldestRecord(pw []domain.RecentWordLog) error {
	or, err := findOldestRecord(pw)
	if err != nil {
		return err
	}

	err = s.repo.DeleteByDate(or)
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultWordService) setTodaysWord() (string, error) {
	// Find previous words from the last 364 days (1 word per day)
	pw, err := s.repo.FindAll()
	if err != nil {
		return "", err
	}
	// Get new word to be stored in the database and todays game
	word, err := s.getNewWord(pw)
	if err != nil {
		return "", err
	}

	// Insert new word into the database
	err = s.repo.InsertWord(s.cfg.Time.Format(isoDateFmt), word)
	if err != nil {
		return "", err
	}

	// Delete oldest record if there are more than or equal to 364 records
	// in the database, we only store the last 364 words or days
	if len(pw) >= recordMax {
		err = s.deleteOldestRecord(pw)
		if err != nil {
			return "", err
		}
	}

	return word, nil
}

func (s *DefaultWordService) GetWord() (string, error) {
	todaysWord, err := s.repo.FindByDate(s.cfg.Time.Format(isoDateFmt))
	if err != nil {
		// If there is no record for todays date, set a new word
		return s.setTodaysWord()
	}

	return todaysWord.Word, nil
}

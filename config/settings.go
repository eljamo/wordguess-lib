package config

import "time"

type Settings struct {
	// Used for generating a deterministic random number
	Time time.Time `key:"time" json:"time,omitempty"`
	// Word Length to use for the game
	WordLength int `key:"word_length" json:"word_length,omitempty"`
	// The word list to use for the game
	WordList string `key:"word_list" json:"word_list,omitempty"`
}

package main

import (
	"fmt"
	"time"

	"github.com/eljamo/libwordle/config"
	"github.com/eljamo/libwordle/internal/sqlite"
	"github.com/eljamo/libwordle/migration"
	"github.com/eljamo/libwordle/service"
)

func main() {
	migrator := &sqlite.Migrator{
		FS:   migration.EmbeddedFiles,
		Path: "sqlite",
		Run:  true,
	}

	db, err := sqlite.New("db.sqlite", migrator)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	svc, err := service.NewDefaultWordService(
		&config.Settings{
			Time:       time.Now(),
			WordLength: 5,
			WordList:   "EN",
		},
		db,
	)
	if err != nil {
		panic(err)
	}

	word, err := svc.GetWord()
	if err != nil {
		panic(err)
	}

	fmt.Println(word)
}

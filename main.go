package main

import (
	"fmt"
	"time"

	"github.com/eljamo/libwordle/config"
	"github.com/eljamo/libwordle/internal/rng"
	"github.com/eljamo/libwordle/internal/sqlite"
	"github.com/eljamo/libwordle/migration"
	"github.com/eljamo/libwordle/repository/gob"
	sqliteRepo "github.com/eljamo/libwordle/repository/sqlite"
	"github.com/eljamo/libwordle/service"
)

func main() {
	rngSvc := rng.NewRNGService()

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

	cfg := &config.Settings{
		AppName:    "wordguess-lib",
		Time:       time.Now(),
		WordLength: 5,
		WordList:   "EN",
	}

	sqlr := sqliteRepo.NewRecentWordLogRepository(db)
	svc, err := service.NewDefaultWordService(cfg, sqlr, rngSvc)
	if err != nil {
		panic(err)
	}

	word, err := svc.GetWord()
	if err != nil {
		panic(err)
	}

	fmt.Println(word)

	grepo := gob.NewRecentWordLogRepository(cfg)
	gsvc, err := service.NewDefaultWordService(cfg, grepo, rngSvc)
	if err != nil {
		panic(err)
	}

	gword, err := gsvc.GetWord()
	if err != nil {
		panic(err)
	}

	fmt.Println(gword)
}

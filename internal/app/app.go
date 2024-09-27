package app

import (
	"os"

	"github.com/VikaPaz/task_tracker/internal/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf := repository.Config{
		DSN: "user=user password=password dbname=users sslmode=disable",
	}

	db := repository.Connection(conf, log.Logger)

	repo := repository.NewTaskRepository(db, log.Logger)

	log.Logger.Printf("Repo: %v", repo)
}

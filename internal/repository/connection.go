package repository

import (
	"database/sql"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Config struct {
	DSN string
}

func Connection(conf Config, logger zerolog.Logger) *bun.DB {
	sqlDB, err := sql.Open("pgx", conf.DSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error connecting to the database.")
	}

	db := bun.NewDB(sqlDB, pgdialect.New())
	return db
}

type TaskRepository struct {
	conn *bun.DB
	log  zerolog.Logger
}

func NewTaskRepository(conn *bun.DB, logger zerolog.Logger) *TaskRepository {
	return &TaskRepository{
		conn: conn,
		log:  logger,
	}
}

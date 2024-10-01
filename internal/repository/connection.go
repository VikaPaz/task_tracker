package repository

import (
	"database/sql"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Config struct {
	DSN string
}

func Connection(conf Config, logger *zerolog.Logger) *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.DSN)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	logger.Debug().Msgf("connection: %s", db)
	return db
}

type TaskRepository struct {
	conn *bun.DB
	log  *zerolog.Logger
}

func NewTaskRepository(conn *bun.DB, logger *zerolog.Logger) *TaskRepository {
	return &TaskRepository{
		conn: conn,
		log:  logger,
	}
}

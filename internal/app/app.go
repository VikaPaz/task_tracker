package app

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/VikaPaz/task_tracker/internal/repository"
	"github.com/VikaPaz/task_tracker/internal/server/grpc"
	"github.com/VikaPaz/task_tracker/internal/server/rest"
	"github.com/VikaPaz/task_tracker/internal/service"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

func Run() {
	// build with "/local/.env"
	err := godotenv.Load("./local/.env")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading local/.env file")
	}

	serverPort := os.Getenv("SERVER_PORT")
	dsn := os.Getenv("DATABASE_URL")

	logger, err := NewLogger()
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating logger")
	}

	conf := repository.Config{
		DSN: dsn,
	}

	db := repository.Connection(conf, logger)

	err = runMigrations(logger, db)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't run migrations")
	}
	logger.Debug().Msg("migrations are applied successfully")

	repo := repository.NewTaskRepository(db, logger)
	logger.Debug().Msg("creaded repository")

	taskService := service.NewTaskService(repo, logger)
	logger.Debug().Msg("creaded sercise")

	taskServer := rest.NewTaskHandler(taskService, logger)
	logger.Debug().Msg("creaded server")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatal().Err(errors.New("panic recovered"))
			}
		}()
		rest.Run(taskServer, serverPort)
	}()
	logger.Info().Msgf("rest server is running on port: %v", serverPort)

	// TODO: add gRPC config
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatal().Err(errors.New("panic recovered"))
			}
		}()
		grpc.NewServer()
	}()
	logger.Info().Msg("grpc server is running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func NewLogger() (*zerolog.Logger, error) {
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	path := os.Getenv("LOG_PATH")

	logFile, err := os.OpenFile(path+"logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	switch loggerLevel {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	return &logger, nil
}

func runMigrations(logger *zerolog.Logger, dbConn *bun.DB) error {
	upMigration, err := strconv.ParseBool(os.Getenv("RUN_MIGRATION"))
	if err != nil {
		return err
	}
	if !upMigration {
		return nil
	}

	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		return errors.New("no migration dir provided; skipping migrations")
	}

	logger.Debug().Msgf("loaded migrations dir: %s", migrationDir)

	err = goose.Up(dbConn.DB, migrationDir)
	if err != nil {
		return err
	}

	logger.Debug().Msg("loaded migrations")

	return nil
}

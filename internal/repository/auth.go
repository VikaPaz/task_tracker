package repository

import (
	"context"
	"database/sql"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type User struct {
	Id       string `bun:"column:pk,type:uuid,default:uuid_generate_v4()"`
	Login    string `bun:"column:notnull"`
	Password string `bun:"column:notnull"`
}

type AuthRepository struct {
	conn *bun.DB
	log  *zerolog.Logger
}

func NewAuthRepository(conn *bun.DB, logger *zerolog.Logger) *AuthRepository {
	return &AuthRepository{
		conn: conn,
		log:  logger,
	}
}

func (r *AuthRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	repoUser := repoUser(user)
	_, err := r.conn.NewInsert().Model(&repoUser).Returning("*").Exec(ctx)
	if err != nil {
		r.log.Error().Err(err).Msgf("can't creating: %v", user)
		return models.User{}, err
	}
	r.log.Debug().Msgf("maked struct %v", user)

	res := modelsUser(repoUser)
	return res, nil
}

func (r *AuthRepository) Get(ctx context.Context, user models.User) (models.User, error) {
	var repoUser User
	err := r.conn.NewSelect().Model(&repoUser).Where("login = ?", user.Login).Where("password = ?", user.Password).Scan(ctx)
	if err != nil {
		r.log.Error().Err(err).Msgf("can't receiving: %v", repoUser)
		if err == sql.ErrNoRows {
			return models.User{}, models.ErrTaskNotFound
		}
		return models.User{}, err
	}
	r.log.Debug().Msgf("received struct %v", repoUser)

	res := modelsUser(repoUser)
	return res, nil
}

func repoUser(user models.User) User {
	return User(user)
}

func modelsUser(user User) models.User {
	return models.User(user)
}

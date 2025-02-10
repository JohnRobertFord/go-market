package storage

import (
	"context"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/JohnRobertFord/go-market/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	createUserQuery = `INSERT INTO users("username", "hash", "created_at") VALUES($1, $2, current_timestamp)`
	getUserByName   = `SELECT hash FROM users WHERE username=$1`
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) NewUser(ctx context.Context, user *model.User) (*model.User, error) {

	hash, _ := util.HashPassword(user.Password)
	_, err := u.db.pgPool.Exec(ctx, createUserQuery, user.Name, hash)
	if e, ok := err.(*pgconn.PgError); ok {
		if e.Code == "23505" {
			return nil, model.ErrDataConflict
		} else {
			return nil, model.ErrInternal
		}
	}
	return user, nil
}

func (u *UserRepository) ValidateUser(ctx context.Context, user *model.User) error {

	var hash string
	err := u.db.pgPool.QueryRow(ctx, getUserByName, user.Name).Scan(&hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrValidate
		} else {
			return model.ErrInternal
		}
	}

	if ok := util.CheckPasswordHash(user.Password, hash); !ok {
		return model.ErrValidate
	}

	return nil
}

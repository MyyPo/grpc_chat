package repositories

import (
	"context"
	"database/sql"

	// "github.com/MyyPo/grpc-chat/db/postgres/public/model"
	. "github.com/MyyPo/grpc-chat/db/postgres/public/table"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	. "github.com/go-jet/jet/v2/postgres"
)

type Auth interface {
	SignIn(context context.Context) error
}

type DBAuth struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) DBAuth {
	return DBAuth{db: db}
}

func (r DBAuth) SignIn(ctx context.Context, req *authpb.SignInRequest) error {
	// !TODO: validation

	stmt := Users.
		SELECT(
			Users.Username,
			Users.Password,
		).FROM(Users).
		WHERE(
			Users.Username.EQ(String(req.GetUsername())).
				AND(Users.Password.EQ(String(req.GetPassword()))),
		)

	_, err := stmt.Exec(r.db)
	return err
}

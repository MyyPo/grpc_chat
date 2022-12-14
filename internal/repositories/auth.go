package repositories

import (
	"context"
	"database/sql"

	"github.com/MyyPo/grpc-chat/db/postgres/public/model"
	. "github.com/MyyPo/grpc-chat/db/postgres/public/table"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	. "github.com/go-jet/jet/v2/postgres"
)

type Auth interface {
	SignIn(ctx context.Context, req *authpb.SignInRequest) (model.Users, error)
	SignUp(ctx context.Context, req *authpb.SignUpRequest) (model.Users, error)
	CheckBlacklistToken(ctx context.Context, req *authpb.RefreshTokenRequest) error
}

type DBAuth struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *DBAuth {
	return &DBAuth{db: db}
}

// fetches userID and password from database for the hasher
func (r DBAuth) SignIn(ctx context.Context, req *authpb.SignInRequest) (model.Users, error) {
	// !TODO: validation

	stmt := Users.
		SELECT(
			Users.UserID,
			Users.Password,
		).FROM(Users).
		WHERE(
			Users.Username.EQ(String(req.GetUsername())),
		)
	var result model.Users
	err := stmt.Query(r.db, &result)
	if err != nil {
		return model.Users{}, err
	}

	return result, nil
}

func (r DBAuth) SignUp(ctx context.Context, req *authpb.SignUpRequest) (model.Users, error) {
	// !TODO: Validation

	stmt := Users.
		INSERT(
			Users.Username,
			Users.Password,
		).VALUES(
		req.GetUsername(),
		req.GetPassword(),
		// "user",
	).RETURNING(
		Users.UserID,
		Users.Username,
	)

	var result model.Users
	err := stmt.Query(r.db, &result)
	if err != nil {
		return model.Users{}, err
	}
	return result, nil
}

// try to insert a new refersh token in the db
// if it fails, return error signaling that the token is invalid
func (r DBAuth) CheckBlacklistToken(ctx context.Context, req *authpb.RefreshTokenRequest) error {
	stmt := BlacklistedTokens.
		INSERT(
			BlacklistedTokens.Value,
		).VALUES(
		req.GetRefreshToken(),
	)

	_, err := stmt.Exec(r.db)
	if err != nil {
		return err
	}

	return nil
}

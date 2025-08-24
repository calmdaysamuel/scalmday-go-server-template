package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	werror "github.com/palantir/witchcraft-go-error"
	"scalmday-go-server-template/config"
	"scalmday-go-server-template/internal/generated/applicationdb"
)

func Init(ctx context.Context, config config.InstallConfig) (db applicationdb.DBTX, closer func(context.Context) error, err error) {
	if config.Database.Type == "postgresql" {
		pgConn, err := pgx.Connect(ctx, config.Database.ConnectionString)
		if err != nil {
			return nil, nil, werror.WrapWithContextParams(ctx, err, "failed to connect to database using connection string")
		}
		closer = pgConn.Close
		db = pgConn
		return db, closer, nil
	}
	return nil, nil, werror.ErrorWithContextParams(ctx, "database type not recognized")
}

package database

import (
	"context"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	werror "github.com/palantir/witchcraft-go-error"
	"go.uber.org/multierr"
	"scalmday-go-server-template/config"
	"scalmday-go-server-template/internal/generated/applicationdb"
)

//go:embed migrations/*.sql
var dbMigrationFolder embed.FS

func Init(ctx context.Context, config config.InstallConfig) (db applicationdb.DBTX, closer func(context.Context) error, err error) {
	dbMigrationSource, err := iofs.New(dbMigrationFolder, "migrations")
	if err != nil {
		return nil, nil, err
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		dbMigrationSource,
		config.Database.ConnectionString)
	if err != nil {
		return nil, nil, werror.WrapWithContextParams(ctx, err, "failed to create migration client")
	}

	if config.Database.PerformDownMigration {
		if err := m.Down(); err != nil {
			return nil, nil, werror.WrapWithContextParams(ctx, err, "failed to perform all up database migrations")
		}
	}
	if config.Database.PerformUpMigration {
		if err := m.Up(); err != nil {
			return nil, nil, werror.WrapWithContextParams(ctx, err, "failed to perform all up database migrations")
		}
	}

	if sourceErr, databaseErr := m.Close(); sourceErr != nil || databaseErr != nil {
		return nil, nil, werror.WrapWithContextParams(ctx, multierr.Combine(sourceErr, databaseErr), "failed to close postgres migration connection")
	}

	if config.Database.Type == "postgres" {
		pgxConn, err := pgx.Connect(ctx, config.Database.ConnectionString)
		if err != nil {
			return nil, nil, werror.WrapWithContextParams(ctx, err, "failed to create primary database connection")
		}
		db = pgxConn
		return db, pgxConn.Close, nil
	}
	return nil, nil, werror.ErrorWithContextParams(ctx, "database type not recognized", werror.SafeParam("db-type", config.Database.Type))
}

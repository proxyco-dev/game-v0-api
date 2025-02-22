package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"game-v0-api/pkg/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func New(config config.Config) *bun.DB {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", config.Database.Host, config.Database.Port)),
		pgdriver.WithUser(config.Database.User),
		pgdriver.WithPassword(config.Database.Password),
		pgdriver.WithDatabase(config.Database.Database),
		pgdriver.WithInsecure(config.Database.Insecure),
		pgdriver.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: false,
			ServerName:         config.Database.Host,
		}),
	)
	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db
}

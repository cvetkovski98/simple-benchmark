package psqlbenchmark

import (
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewPgDb(dns string) (*bun.DB, error) {
	log.Println("Connecting to PostgreSQL database...")
	var connector = pgdriver.NewConnector(pgdriver.WithDSN(dns))
	var db = sql.OpenDB(connector)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return bun.NewDB(db, pgdialect.New()), nil
}

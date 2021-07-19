package drivers

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var con = &DB{}

const maxOpenDBConnection = 10
const maxIdleDBConnection = 5
const maxDbLifeTime = 5 * time.Minute

//ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDBConnection)
	d.SetMaxIdleConns(maxIdleDBConnection)
	d.SetConnMaxLifetime(maxDbLifeTime)

	con.SQL = d

	err = testDB(d)

	if err != nil {
		return nil, err
	}

	return con, nil
}

//testDB tries to ping database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}

	return nil
}

//NewDatabase create a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

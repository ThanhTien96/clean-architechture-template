package database

import (
	"backend-agent-demo/adapter/repository"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type postgresHandler struct {
	db *sql.DB
}

func NewPostgresHandler(c *config) (*postgresHandler, error) {
	ds := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.host,
		c.port,
		c.user,
		c.database,
		c.password,
	)

	log.Println(ds)
	db, err := sql.Open(c.driver, ds)
	if err != nil {
		return &postgresHandler{}, nil
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	return &postgresHandler{db: db}, nil
}


func (pg postgresHandler) BeginTx(ctx context.Context) (repository.Tx, error) {
	tx, err := pg.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return postgresTx{}, nil
	}

	return newPostgresTx(tx), nil
}

func (pg postgresHandler) ExecuteContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg postgresHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (repository.Rows, error) {
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return newPostgresRows(rows), nil
}

func (p postgresHandler) QueryRowContext(ctx context.Context, query string, args ...interface{}) repository.Row {
	row := p.db.QueryRowContext(ctx, query, args...)

	return newPostgresRow(row)
}

// POSTGRES ROW
type postgresRow struct {
	row *sql.Row
}

func newPostgresRow(row *sql.Row) postgresRow {
	return postgresRow{row: row}
}

func (pr postgresRow) Scan(dest ...interface{}) error {
	if err := pr.row.Scan(dest...); err != nil {
		return err
	}
	return nil
}


// POSTGRES ROWS
type postgresRows struct {
	rows *sql.Rows
}

func newPostgresRows(rows *sql.Rows) postgresRows {
	return postgresRows{rows: rows}
}

func (prs postgresRows) Scan(dest ...interface{}) error {
	if err := prs.rows.Scan(dest...); err != nil {
		return err
	}

	return nil
}

func (prs postgresRows) Next() bool {
	return prs.rows.Next()
}

func (prs postgresRows) Err() error {
	return prs.rows.Err()
}

func (prs postgresRows) Close() error {
	return prs.rows.Close()
}


// POSTGRES TRANSACTION
type postgresTx struct {
	tx *sql.Tx
}

func newPostgresTx(tx *sql.Tx) postgresTx {
	return postgresTx{tx: tx}
}

func (p postgresTx) ExecuteContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p postgresTx) QueryContext(ctx context.Context, query string, args ...interface{}) (repository.Rows, error) {
	rows, err := p.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	row := newPostgresRows(rows)

	return row, nil
}

func (p postgresTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) repository.Row {
	row := p.tx.QueryRowContext(ctx, query, args...)

	return newPostgresRow(row)
}

func (p postgresTx) Commit() error {
	return p.tx.Commit()
}

func (p postgresTx) Rollback() error {
	return p.tx.Rollback()
}

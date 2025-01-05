package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"

	"github.com/lib/pq"

	"dblog/sqllog"
)

/*
	Driver implementation.
*/

// Name is the name of the driver.
const Name = "pq-with-logging"

// driver wraps the original database/sql/driver.Conn adding logging functionality.
type driver struct {
	sqldriver.Driver
}

func init() {
	// Register the driver with the name "pq-with-logging", so that it can be used in sql.Open.
	sql.Register(Name, driver{Driver: &pq.Driver{}})
}

// Open wraps the original sqldriver.Conn with conn
func (d driver) Open(name string) (sqldriver.Conn, error) {
	c, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}

	return &conn{Conn: c}, nil
}

/*
	Connection implementation.
*/

type conn struct {
	sqldriver.Conn
}

// Exec logs and executes a query without context
func (c *conn) Exec(query string, args []sqldriver.Value) (sqldriver.Result, error) {
	sqllog.InfoSQL("Exec", query, args)
	if ec, ok := c.Conn.(sqldriver.Execer); ok {
		return ec.Exec(query, args)
	}

	return nil, sqldriver.ErrSkip
}

// ExecContext logs and executes a query with context
func (c *conn) ExecContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Result, error) {
	sqllog.InfoSQLContext(ctx, "ExecContext", query, args)
	if ec, ok := c.Conn.(sqldriver.ExecerContext); ok {
		return ec.ExecContext(ctx, query, args)
	}

	return nil, sqldriver.ErrSkip
}

func (c *conn) Query(query string, args []sqldriver.Value) (sqldriver.Rows, error) {
	sqllog.InfoSQL("Query", query, args)
	if qc, ok := c.Conn.(sqldriver.Queryer); ok {
		return qc.Query(query, args)
	}

	return nil, sqldriver.ErrSkip
}

// QueryContext logs and queries with context
func (c *conn) QueryContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Rows, error) {
	sqllog.InfoSQLContext(ctx, "QueryContext", query, args)
	if qc, ok := c.Conn.(sqldriver.QueryerContext); ok {
		return qc.QueryContext(ctx, query, args)
	}

	return nil, sqldriver.ErrSkip
}

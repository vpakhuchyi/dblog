package main

import (
	"context"
	"database/sql"
	"log"

	"dblog/driver"
	"dblog/trace"
)

func main() {
	dsn := "postgres://tester:password@localhost/postgres?sslmode=disable"

	// Use our custom "pq-with-logging" driver instead of github.com/lib/pq driver.
	db, err := sql.Open(driver.Name, dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Generate and put the correlation ID into the context, so it can be retrieved later.
	// This is a simulation of a middleware that adds the correlation ID to the request context.
	// An example of such middleware can be found in the trace package: trace.WithCorrelationID.
	ctx := trace.SetCorrelationID(context.Background(), "1234")

	query := "SELECT users.id, users.* FROM users WHERE users.id = $1"
	db.Query(query, 1)
	/*
	   2025-01-05T20:21:59.359+0100    INFO    QueryContext correlation_id=
	   Script:
	          SELECT users.id, users.*
	          FROM users
	          WHERE users.id = $1
	   Args:
	         [1]
	*/

	db.QueryContext(ctx, query, 2)
	/*
	   2025-01-05T20:21:59.382+0100    INFO    QueryContext correlation_id=1234
	   Script:
	          SELECT users.id, users.*
	          FROM users
	          WHERE users.id = $1
	   Args:
	         [2]
	*/

	exec := "INSERT INTO users (name, email) VALUES ($1, $2)"
	db.Exec(exec, "Viktor", "v.example@example.com")
	/*
		2025-01-05T20:21:59.405+0100    INFO    ExecContext correlation_id=
		Script:
		       INSERT INTO users (name, email)
		       VALUES ($1, $2)
		Args:
		      [Viktor, v.example@example.com]
	*/

	db.ExecContext(ctx, exec, "Viktor Pakhuchyi", "vp.example@example.com")
	/*
		2025-01-05T20:21:59.406+0100    INFO    ExecContext correlation_id=1234
		Script:
		       INSERT INTO users (name, email)
		       VALUES ($1, $2)
		Args:
		      [Viktor Pakhuchyi, vp.example@example.com]
	*/
}

package main

import (
	"context"
	"database/sql"
	"log"

	"dblog/driver"
	"dblog/trace"
)

func main() {
	// Connection setting to the database.
	dsn := "postgres://tester:password@localhost/postgres?sslmode=disable"

	// Use our custom "pq-with-logging" driver instead of pq's "postgres".
	db, err := sql.Open(driver.Name, dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Generate and put the correlation ID into the context, so it can be retrieved later.
	ctx := trace.SetCorrelationID(context.Background(), "1234")

	query := "SELECT users.id, users.* FROM users WHERE users.id = $1"
	db.Query(query, 1)
	db.QueryContext(ctx, query, 2)

	db.QueryRow(query, 3)
	db.QueryRowContext(ctx, query)

	exec := "INSERT INTO users (name, email) VALUES ($1, $2)"
	db.Exec(exec, "Bill Doe", "b.example@example.com")
	db.ExecContext(ctx, exec, "John Doe", "j.example@example.com")
}

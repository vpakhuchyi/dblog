# SQL Queries Logger Example in Go

This repository provides a very basic example of a SQL queries logger, showcasing how to:

1. **Wrap SQL drivers to add custom logging functionality**
2. **Implement a simple custom SQL query formatter**
3. **Use correlation IDs to trace SQL queries back to specific requests or use cases**

The goal is to provide a starting point for building more advanced solutions for logging and formatting SQL queries.   
Many popular ORMs or SQL-builders provide query logging, but not all of them. This example could be helpful if you plan 
to use a library that covers all your other needs but lacks SQL query logging for debugging.

## Implementation Steps

1. Create a custom logger wrapper that accepts `context.Context` (to pass correlation IDs), query and arguments (optionally).
2. Create/use an SQL formatting function to print queries and arguments in a readable way.
3. Create a custom SQL driver wrapping another suitable SQL driver.
4. Decorate the required SQL driver methods to invoke the logger with `ctx`, `query`, and/or `args`.
5. Use the `trace.WithCorrelationID` function to generate and add a unique correlation ID to the `context.Context` and log it with queries.

## Usage Example

```go
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
```

## Contributing

Feel free to submit issues or pull requests to improve this repository. Suggestions for new features or improvements are always welcome!


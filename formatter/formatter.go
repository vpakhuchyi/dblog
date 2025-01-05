package formatter

import (
	sqldriver "database/sql/driver"
	"fmt"
	"strings"
)

// Store both upper and lower case keywords for faster lookup,
// cos it allows to avoid calling strings.ToUpper for every word.
var keywords = map[string]bool{
	"SELECT": true, "select": true,
	"FROM": true, "from": true,
	"INSERT": true, "insert": true,
	"WHERE": true, "where": true,
	"GROUP": true, "group": true,
	"ORDER": true, "order": true,
	"LIMIT": true, "limit": true,
	"VALUES": true, "values": true,
	"UPDATE": true, "update": true,
	"SET": true, "set": true,
	"DELETE": true, "delete": true,
	"LEFT": true, "left": true,
	"RIGHT": true, "right": true,
	"INNER": true, "inner": true,
	"OUTER": true, "outer": true,
	"FULL": true, "full": true,
	"JOIN": true, "join": true,
	"ON": true, "on": true,
	"USING": true, "using": true,
	"HAVING": true, "having": true,
	"DISTINCT": true, "distinct": true,
	"AS": true, "as": true,
	"CASE": true, "case": true,
	"WHEN": true, "when": true,
	"THEN": true, "then": true,
	"ELSE": true, "else": true,
	"END": true, "end": true,
	"UNION": true, "union": true,
	"INTERSECT": true, "intersect": true,
	"EXCEPT": true, "except": true,
	"ALL": true, "all": true,
	"CREATE": true, "create": true,
	"ALTER": true, "alter": true,
	"DROP": true, "drop": true,
	"TRUNCATE": true, "truncate": true,
	"WITH": true, "with": true,
	"RETURNING": true, "returning": true,
	"OFFSET": true, "offset": true,
}

// Query formats the query string to be more readable with new lines and indentation.
func Query(query string) string {
	var b strings.Builder
	b.WriteString("\nScript:")

	for _, word := range strings.Fields(query) {
		// Check if the word is a keyword to handle it properly.
		_, ok := keywords[word]
		if ok {
			b.WriteString("\n        ")
		}

		b.WriteString(word + " ")
	}

	result := strings.ReplaceAll(b.String(), "\n ", "\n")
	result = strings.ReplaceAll(result, "\n)", ")")

	return result
}

// QueryWithArgs formats the query string with arguments to be more readable with new lines and indentation.
func QueryWithArgs(query string, args any) string {
	var b strings.Builder
	b.WriteString(Query(query))

	if args == nil {
		return b.String()
	}

	b.WriteString("\nArgs:\n      [")

	// Depending on th e type of options, the type of arguments can be different.
	// So, here we need to check possible types and handle them properly.
	switch values := args.(type) {
	case []sqldriver.Value:
		for i, v := range values {
			b.WriteString(fmt.Sprintf("%v", v))
			if i < len(values)-1 {
				b.WriteString(", ")
			}
		}
	case []sqldriver.NamedValue:
		for i, v := range values {
			b.WriteString(fmt.Sprintf("%v", v.Value))
			if i < len(values)-1 {
				b.WriteString(", ")
			}
		}
	default:
		b.WriteString("#unknown_value")
	}

	b.WriteString("]")

	return b.String()
}

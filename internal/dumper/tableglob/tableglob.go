package tableglob

import (
	"database/sql"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
)

// Show tables based on a list of globs.
func Show(connection string, globs []string) ([]string, error) {
	var globbed []string

	tables, err := getTables(connection)
	if err != nil {
		return globbed, errors.Wrap(err, "failed to query for tables")
	}

	for _, query := range globs {
		g := glob.MustCompile(query)

		for _, table := range tables {
			if g.Match(table) {
				globbed = appendIfMissing(globbed, table)
			}
		}
	}

	return globbed, nil
}

// Helper function to get a list of tables.
func getTables(connection string) ([]string, error) {
	var tables []string

	db, err := sql.Open("mysql", connection)
	if err != nil {
		return tables, err
	}

	rows, err := db.Query("SHOW FULL TABLES")
	if err != nil {
		return tables, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, tableType string

		err := rows.Scan(&tableName, &tableType)
		if err != nil {
			return tables, err
		}

		if tableType == "BASE TABLE" {
			tables = append(tables, tableName)
		}
	}

	return tables, nil
}

// Helper function to append to a slice if not present in a slice.
func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}

	return append(slice, i)
}

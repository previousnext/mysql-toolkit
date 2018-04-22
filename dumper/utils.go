package dumper

import (
	"fmt"
)

// String representation of a MySQL connection.
func (c Connection) String() string {
	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s", c.Username, c.Password, c.Protocol, c.Hostname, c.Port, c.Database)
}

// Map list for sanitizing the MySQL dump.
func (s Sanitize) Map() map[string]map[string]string {
	selectMap := make(map[string]map[string]string)

	for _, table := range s.Tables {
		tableMap := make(map[string]string, len(table.Fields))

		for _, field := range table.Fields {
			if field.Value == "" {
				field.Value = DefaultPlaceholder
			}

			tableMap[field.Name] = fmt.Sprintf("'%s'", field.Value)
		}

		selectMap[table.Name] = tableMap
	}

	return selectMap
}

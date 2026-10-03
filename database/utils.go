package database

import (
	"fmt"
	"strings"
)

func buildQuery(filters map[string]string, separator string) (string, []any) {
	allowed := map[string]bool{
		"name":     true,
		"username": true,
		"user_id":  true,
		"start_at": true,
		"end_at":   true,
		"repeat":   true,
		"child_id": true,
	}

	var conditions []string
	var args []any
	i := 1

	for key, val := range filters {
		if !allowed[key] {
			continue
		}

		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, val)
		i++
	}

	query := ""
	if len(conditions) > 0 {
		query += strings.Join(conditions, separator)
	}
	return query, args
}

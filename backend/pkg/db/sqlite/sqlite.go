package sqlite

import (
	"database/sql"
)

var DB *sql.DB

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

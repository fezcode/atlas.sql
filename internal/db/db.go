package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
	"os"
)

type Result struct {
	Columns []string
	Rows    [][]string
	Error   error
}

type Database interface {
	Query(query string) (*Result, error)
	ListTables() (*Result, error)
	ListSchemas() (*Result, error)
	Close() error
}

type sqlDB struct {
	db     *sql.DB
	driver string
}

func Connect(connStr string) (Database, error) {
	var driver, dsn string

	if strings.HasPrefix(connStr, "sqlite://") {
		driver = "sqlite"
		dsn = strings.TrimPrefix(connStr, "sqlite://")
		
		// Check if file exists for sqlite
		if _, err := os.Stat(dsn); os.IsNotExist(err) {
			return nil, fmt.Errorf("sqlite database file not found: %s", dsn)
		}
	} else if strings.HasPrefix(connStr, "postgres://") {
		driver = "postgres"
		dsn = connStr
	} else {
		return nil, fmt.Errorf("unsupported database driver. use sqlite:// or postgres://")
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &sqlDB{db: db, driver: driver}, nil
}

func (s *sqlDB) Query(query string) (*Result, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var resultRows [][]string
	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make([]string, len(cols))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		resultRows = append(resultRows, row)
	}

	return &Result{
		Columns: cols,
		Rows:    resultRows,
	}, nil
}

func (s *sqlDB) ListTables() (*Result, error) {
	var query string
	if s.driver == "sqlite" {
		query = "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';"
	} else {
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema NOT IN ('information_schema', 'pg_catalog') ORDER BY table_name;"
	}
	return s.Query(query)
}

func (s *sqlDB) ListSchemas() (*Result, error) {
	var query string
	if s.driver == "sqlite" {
		// SQLite doesn't really have schemas in the same way, but it has attached databases
		query = "PRAGMA database_list;"
	} else {
		query = "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'pg_catalog') ORDER BY schema_name;"
	}
	return s.Query(query)
}

func (s *sqlDB) Close() error {
	return s.db.Close()
}

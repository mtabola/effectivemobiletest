package pgsql

import (
	"database/sql"
	"effectivemobiletest/internal/config"

	_ "github.com/lib/pq"

	"fmt"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbcp config.DBConfig) (*Storage, error) {
	connlink := fmt.Sprintf("host = %s port = %d user=%s password=%s dbname=%s sslmode=disable", dbcp.Host, dbcp.Port, dbcp.User, dbcp.Password, dbcp.DBName)
	db, err := sql.Open("postgres", connlink)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(dbcp.Tables); i++ {
		var tableStatus bool
		req := db.QueryRow(fmt.Sprintf(`
		SELECT EXISTS (
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name='%s')`,
			dbcp.Tables[i]))

		err = req.Scan(&tableStatus)
		if err != nil {
			return nil, err
		}
		if !tableStatus {
			slog.Error("Table %s doesn't exists. Use migrations files", dbcp.Tables[i])
			return nil, fmt.Errorf("table %s doesn't exist", dbcp.Tables[i])
		}
	}
	return &Storage{db: db}, nil
}

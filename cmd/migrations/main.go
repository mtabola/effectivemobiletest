package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, dbLink string

	flag.StringVar(&dbLink, "dbcl", "", "* database connection link")
	flag.StringVar(&migrationsPath, "mp", "", "* path to migrations")

	flag.Parse()

	if dbLink == "" || migrationsPath == "" {
		panic("sp, mp and mt flags are required")
	}

	m, err := migrate.New("file://"+migrationsPath, dbLink)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("Migrations are applied successfully")
}

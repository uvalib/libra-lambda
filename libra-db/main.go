package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	var newMigrationName string
	var down bool
	flag.StringVar(&newMigrationName, "create", "", "use -create migration_name to create a new migration")
	flag.BoolVar(&down, "down", false, "migrate down")

	flag.Parse()

	if newMigrationName != "" {

		cmd := exec.Command("migrate",
			"create", "-dir", "./migrations", "-ext", ".sql", newMigrationName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		return
	}

	err = godotenv.Load("env.local")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION"))
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if down {
		err = m.Down()
		fmt.Println(err)
		return
	}
	err = m.Up()
	fmt.Println(err)

}

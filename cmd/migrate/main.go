package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/inganta23/wallet/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	direction := flag.String("direction", "up", "migration direction")
	flag.Parse()

	m, err := migrate.New("file://migrations", cfg.DBUrl)
	if err != nil {
		log.Fatal("Migration setup failed: ", err)
	}

	if *direction == "down" {
		fmt.Println("Rolling back...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Migrating up...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	fmt.Println("Done!")
}
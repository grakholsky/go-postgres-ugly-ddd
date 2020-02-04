package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

const (
	AccountRole = "account"
)

type PostgresConfig struct {
	URI string
}

func ConnectPostgres(config *PostgresConfig) *gorm.DB {
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open("postgres", config.URI)
		if err != nil {
			fmt.Printf("connect to postgres failed: %v", err)
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err := Migrate(db.DB()); err != nil {
		log.Fatalf("run migrations to postgres failed: %v", err)
	}

	// Enable Logger, show detailed log
	db.LogMode(true)
	return db
}

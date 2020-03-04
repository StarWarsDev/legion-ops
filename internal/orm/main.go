package orm

import (
	"log"
	"os"

	"github.com/StarWarsDev/legion-ops/internal/orm/migration"

	// postgres dialect
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ORM struct {
	DB *gorm.DB
}

func Factory() (ORM, error) {
	host := os.Getenv("DATABASE")
	if host == "" {
		host = "host=localhost dbname=postgres sslmode=disable"
	}
	db, err := gorm.Open("postgres", host)
	if err != nil {
		log.Fatal("[ORM] err: ", err)
	}

	orm := ORM{
		DB: db,
	}

	db.LogMode(true)
	err = migration.ServiceAutoMigration(orm.DB)

	log.Println("[ORM] Database connection initialized")

	return orm, err
}

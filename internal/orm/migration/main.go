package migration

import (
	"fmt"
	"log"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	"github.com/StarWarsDev/legion-ops/internal/orm/migration/jobs"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func updateMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&event.Event{},
		&event.Day{},
		&event.Round{},
		&event.Match{},
	).Error
}

func ServiceAutoMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
	m.InitSchema(func(db *gorm.DB) error {
		log.Println("[Migration.InitSchema] Initializing database schema")
		switch db.Dialect().GetName() {
		case "postgres":
			db.Exec("create extension \"uuid-ossp\";")
		}
		if err := updateMigration(db); err != nil {
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}

		return nil
	})
	m.Migrate()

	if err := updateMigration(db); err != nil {
		return err
	}
	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		jobs.SeedUsers,
		jobs.SeedEvents,
	})

	return m.Migrate()
}

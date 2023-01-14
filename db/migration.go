package db

import (
	"embed"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, embedMigrations embed.FS) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(sqlDb, "migrations"); err != nil {
		return err
	}
	return nil
}

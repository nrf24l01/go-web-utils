package pgkit

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

// ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
// ░░░░ЗАПУСКАЕМ░ГУСЕЙ-РАЗВЕДЧИКОВ░░░░
// ░░░░░▄▀▀▀▄░░░▄▀▀▀▀▄░░░▄▀▀▀▄░░░░░
// ▄███▀░◐░░░▌░▐0░░░░0▌░▐░░░◐░▀███▄
// ░░░░▌░░░░░▐░▌░▐▀▀▌░▐░▌░░░░░▐░░░░
// ░░░░▐░░░░░▐░▌░▌▒▒▐░▐░▌░░░░░▌░░░░
// ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
func RunMigrations(db *sql.DB, migrationsDir string) error {
	goose.SetDialect("postgres")

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}

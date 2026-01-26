package pgkit

import (
	"database/sql"

	"github.com/nrf24l01/go-web-utils/config"
	"github.com/pressly/goose/v3"
)

// ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
// ░░░░ЗАПУСКАЕМ░ГУСЕЙ-РАЗВЕДЧИКОВ░░░░
// ░░░░░▄▀▀▀▄░░░▄▀▀▀▀▄░░░▄▀▀▀▄░░░░░
// ▄███▀░◐░░░▌░▐0░░░░0▌░▐░░░◐░▀███▄
// ░░░░▌░░░░░▐░▌░▐▀▀▌░▐░▌░░░░░▐░░░░
// ░░░░▐░░░░░▐░▌░▌▒▒▐░▐░▌░░░░░▌░░░░
// ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
func RunMigrations(db *sql.DB, pg_cfg *config.PGConfig) error {
	goose.SetDialect("postgres")

	if err := goose.Up(db, pg_cfg.Migrations); err != nil {
		return err
	}

	return nil
}

package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersChecksTable = "users_checks"
)

type Config struct {
	Host    string
	Port    string
	DBName  string
	SSLMode string
}

type Repository struct {
	db *sqlx.DB
}

func New(cfg Config) (*Repository, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CountChecks(userId int64, timeLimit int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	addCheck := fmt.Sprintf("INSERT INTO %s (id, time) VALUES ($1, now()) ", usersChecksTable)
	if _, err := tx.Exec(addCheck, userId); err != nil {
		tx.Rollback()
		return 0, err
	}

	var checks int
	countChecks := fmt.Sprintf(`SELECT COUNT(*) AS TotalCount FROM %s
							WHERE time >= current_timestamp - INTERVAL '1 SECOND' * $2 AND id = $1`, usersChecksTable)
	row := tx.QueryRow(countChecks, userId, timeLimit)
	if err := row.Scan(&checks); err != nil {
		tx.Rollback()
		return 0, err
	}

	delOld := fmt.Sprintf(`DELETE FROM %s
							WHERE EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - time)) > $1`, usersChecksTable)
	if _, err := tx.Exec(delOld, timeLimit); err != nil {
		tx.Rollback()
		return 0, err
	}

	return checks, tx.Commit()
}

package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersChecksTable = "users_checks"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Repository struct {
	db *sqlx.DB
}

func New(cfg Config) (*Repository, error) {
	db, err := sqlx.Open("postgres",
		// fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		// cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
		fmt.Sprintf("host=%s dbname=%s sslmode=%s", cfg.Host, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetChecks(userId int64) (int, error) {
	var checks int
	query := fmt.Sprintf("SELECT checks FROM %s WHERE user_id=$1", usersChecksTable)
	err := r.db.Get(&checks, query, userId)

	return checks, err
}

func (r *Repository) AddCheck(userId int64) (int, error) {
	var checks int
	query := fmt.Sprintf("UPDATE %s SET checks = checks+1 WHERE id=$1 RETURNING checks", usersChecksTable)
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&checks); err != nil {
		return 0, err
	}

	return checks, nil
}


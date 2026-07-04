package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type TreeWalker interface {
	InOrderWalk(fn func(key string, translations []string) bool)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, db.Ping()
}

func (p *PostgresStorage) LoadToTree(insertFn func(eng, rus string) error) error {
	rows, err := p.db.Query("SELECT word_eng, word_rus FROM translations")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var eng, rus string
		if err := rows.Scan(&eng, &rus); err != nil {
			return err
		}
		_ = insertFn(eng, rus)
	}
	return rows.Err()
}

func (p *PostgresStorage) SaveFromTree(walker TreeWalker) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO translations (word_eng, word_rus) 
		VALUES ($1, $2) 
		ON CONFLICT (word_eng, word_rus) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	walker.InOrderWalk(func(key string, translations []string) bool {
		for _, tr := range translations {
			if _, err := stmt.Exec(key, tr); err != nil {
				return false
			}
		}
		return true
	})

	return tx.Commit()
}

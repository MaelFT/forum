package forum

import (
    "database/sql"
)

type SQLiteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
    return &SQLiteRepository{
        db: db,
    }
}

func (r *SQLiteRepository) Migrate() error {
    query := `
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        mail TEXT NOT NULL UNIQUE,
        checkedmail INTEGER NOT NULL,
        password TEXT NOT NULL,
        categories TEXT,
        date DATE
    );
    `

    _, err := r.db.Exec(query)
    return err
}
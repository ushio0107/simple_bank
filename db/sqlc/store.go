// Implement ACID, roll back when fails.
package db

import (
	"context"
	"database/sql"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	// This is an anonymous field.
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	return nil
}

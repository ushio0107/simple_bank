// Implement ACID, roll back when fails.
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute db queries and transactions.
type SQLStore struct {
	// This is an anonymous field.
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within database
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{} // The second bracket means that we're creating a new empty object of that type.

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, account entries, and update accounts' balance within a single database transaction.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// For debug, log information.
		txName := ctx.Value(txKey)

		fmt.Println(txName, " create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, " create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, " create entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Use TDD, keep writing test case and find bug, then keep improving it.
		/*
			Wrong implement, because 2 concurrent transactions can get the same value of the account 1.

			ac1, err := q.GetAccount(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: ac1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			ac2, err := q.GetAccount(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: ac2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
		*/

		// Let the smaller account ID to transfer first everytime to avoid deadlock.
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		/*
			First Approach
			fmt.Println(txName, " get account1")
			ac1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

			fmt.Println(txName, " update account1")
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: ac1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			fmt.Println(txName, " get account2")
			ac2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

			fmt.Println(txName, " update account2")
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: ac2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
		*/

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (ac1 Account, ac2 Account, err error) {
	ac1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	ac2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}

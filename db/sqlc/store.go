package db

import (
	"database/sql"
	"fmt"
)


type Store struct{
	*Queries
	db *sql.DB
}

//create new store

func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db)
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries)error)error{
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
} 


type TransferTxParams struct{
	FromAccountID int64 `json:from_account_id`
	ToAccountID int64 `json:to_account_id`
	Amount int64 `json:amount`
}

type TransferTxResults struct{
	Transfer Transfer `json:transfer`
	FromAccount Account `json:from_account`
	ToAccount Account `json:to_account`
	FromEntry Entry `json:from_entry`
	ToEntry Entry `json:to_entry`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResults error){
	var result TransferTxResults
	err := store.execTx(ctx func(q *Queries) error {
		var err error
		result, err = q.createTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount
		})
		if err != nil{
			return err
		}

		FromEntryResult, FromEntryErr := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if FromEntryErr != nil {
			return FromEntryErr
		}
		ToEntryResult, ToEntryErr := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: -arg.Amount,
		})
		if ToEntryErr != nil {
			return ToEntryErr
		}

		//TODO: update accounts' balance
		return nil;
	})
	return result, err
}
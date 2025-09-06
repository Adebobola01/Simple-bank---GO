package db

import (
	"context"
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
		Queries: New(db),
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

type TransferTxResults struct{
	Transfer Transfer `json:transfer`
	FromAccount Account `json:from_account`
	ToAccount Account `json:to_account`
	FromEntry Entry `json:from_entry`
	ToEntry Entry `json:to_entry`
}

type TransferTxParams struct{
	FromAccountID int64 `json:from_account_id`
	ToAccountID int64 `json:to_account_id`
	Amount int64 `json:amount`
}



func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (Transfer, error){
	
	var result Transfer
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID: args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil{
			return err
		}

		//handle balance change
		err = transferBalanceChangeHandler(ctx, q, args.FromAccountID, args.ToAccountID, args.Amount)
		if err != nil{
			//delete transfer then return
			errorResponse := q.DeleteTransfer(ctx, result.ID)
			if errorResponse != nil{
				return errorResponse
			}
			return err
		}


		_, FromEntryErr := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount: -args.Amount,
		})
		if FromEntryErr != nil {
			return FromEntryErr
		}
		_, ToEntryErr := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount: -args.Amount,
		})
		if ToEntryErr != nil {
			return ToEntryErr
		}

		//TODO: update accounts' balance
		return nil;
	})
	return result, err
}


func transferBalanceChangeHandler(ctx context.Context, q *Queries, fromAccountId int64, toAccounId int64, amount int64) error {

	//update sender account
	acc1, err := q.GetAccount(ctx, fromAccountId)
	if err != nil{
		return err
	}

	if acc1.Balance < amount {
		return fmt.Errorf("Error: Balance Insufficient")
	}

	fromAccArgs := UpdateAccountParams{
		ID: acc1.ID,
		Balance: acc1.Balance - amount,
	}
	err = q.UpdateAccount(ctx, fromAccArgs)
	if err != nil{
		return err
	}

	//handler recipient balance change
	acc2, err := q.GetAccount(ctx, toAccounId)
	if err != nil{
		return err
	}
	
	toAccountArgs := UpdateAccountParams{
		ID: acc2.ID,
		Balance: acc2.Balance + amount,
	}

	err = q.UpdateAccount(ctx, toAccountArgs)

	if err != nil{
		return err
	}
	return err
}
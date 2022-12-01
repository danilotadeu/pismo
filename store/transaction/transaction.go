package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

//go:generate mockgen -destination ../../mock/store/transaction/transaction_store_mock.go -package mockStoreTransaction . Store
type Store interface {
	CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*int64, error)
}

type storeImpl struct {
	db *sql.DB
}

// NewApp init a account
func NewStore(db *sql.DB) Store {
	return &storeImpl{
		db: db,
	}
}

// CreateTransaction create a transaction..
func (a *storeImpl) CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*int64, error) {
	query := fmt.Sprintf("INSERT INTO transactions(account_id,operation_type_id,amount) VALUES ('%d','%d','%f')", accountID, operationTypeID, amount)
	res, err := a.db.Exec(query)

	if err != nil {
		log.Println("store.transaction.CreateTransaction.Exec", err.Error())
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println("store.account.CreateTransaction.LastInsertId", err.Error())
		return nil, err
	}
	return &lastId, nil
}

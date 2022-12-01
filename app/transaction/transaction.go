package transaction

import (
	"context"
	"errors"
	"log"

	accountModel "github.com/danilotadeu/pismo/model/account"
	"github.com/danilotadeu/pismo/model/transaction"
	"github.com/danilotadeu/pismo/store"
)

//go:generate mockgen -destination ../../mock/app/transaction/transaction_app_mock.go -package mockAppTransaction . App
type App interface {
	CreateTransaction(ctx context.Context, transactionBody transaction.TransactionRequest) (*int64, error)
}

type appImpl struct {
	store *store.Container
}

// NewApp init a transaction
func NewApp(store *store.Container) App {
	return &appImpl{
		store: store,
	}
}

// CreateAccount create a account..
func (a *appImpl) CreateTransaction(ctx context.Context, transactionBody transaction.TransactionRequest) (*int64, error) {
	_, ok := transaction.OperationTypes[transactionBody.OperationTypeID]
	if ok {
		_, err := a.store.Account.GetAccount(ctx, transactionBody.AccountID)
		if err != nil {
			log.Println("app.transaction.CreateTransaction.GetAccount", err.Error())
			if errors.Is(err, accountModel.ErrorAccountNotFound) {
				return nil, accountModel.ErrorAccountNotFound
			}
			return nil, err
		}
		valueAmount := transactionBody.Amount
		_, ok := transaction.OperationTypesBuyOrWithdraw[transactionBody.OperationTypeID]
		if ok {
			valueAmount = -transactionBody.Amount
		}
		transactionID, err := a.store.Transaction.CreateTransaction(ctx, transactionBody.AccountID, transactionBody.OperationTypeID, valueAmount)
		if err != nil {
			log.Println("app.transaction.CreateTransaction.CreateTransaction", err.Error())
			return nil, err
		}

		return transactionID, nil
	} else {
		return nil, transaction.ErrorTransactionTypeNotFound
	}
}

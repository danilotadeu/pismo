package account

import (
	"context"
	"log"

	accountModel "github.com/danilotadeu/pismo/model/account"
	"github.com/danilotadeu/pismo/store"
)

//go:generate mockgen -destination ../../mock/app/account/account_app_mock.go -package mockAppAccount . App
type App interface {
	CreateAccount(ctx context.Context, documentNumber string) (*int64, error)
	GetAccount(ctx context.Context, accountId int64) (*accountModel.AccountResultQuery, error)
	GetAllAccounts(ctx context.Context) ([]*accountModel.AccountResultQuery, error)
}

type appImpl struct {
	store *store.Container
}

// NewApp init a account
func NewApp(store *store.Container) App {
	return &appImpl{
		store: store,
	}
}

// CreateAccount create a account..
func (a *appImpl) CreateAccount(ctx context.Context, documentNumber string) (*int64, error) {
	count, err := a.store.Account.GetAccountByDocumentNumber(ctx, documentNumber)
	if err != nil {
		log.Println("app.account.CreateAccount.GetAccountByDocumentNumber", err.Error())
		return nil, err
	}

	if count.Count > 0 {
		return nil, accountModel.ErrorAccountExists
	}

	id, err := a.store.Account.StoreCreateAccount(ctx, documentNumber)
	if err != nil {
		log.Println("app.account.CreateAccount.StoreCreateAccount", err.Error())
		return nil, err
	}
	return id, nil
}

// GetAccount get a account..
func (a *appImpl) GetAccount(ctx context.Context, accountId int64) (*accountModel.AccountResultQuery, error) {
	account, err := a.store.Account.GetAccount(ctx, accountId)
	if err != nil {
		log.Println("app.account.GetAccount.GetAccount", err.Error())
		return nil, err
	}
	return account, nil
}

// GetAllAccounts get aall accounts..
func (a *appImpl) GetAllAccounts(ctx context.Context) ([]*accountModel.AccountResultQuery, error) {
	accounts, err := a.store.Account.GetAllAccounts(ctx)
	if err != nil {
		log.Println("app.account.GetAccount.GetAccount", err.Error())
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, accountModel.ErrorAccountListIsEmpty
	}

	return accounts, nil
}

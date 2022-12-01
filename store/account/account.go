package account

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	accountModel "github.com/danilotadeu/pismo/model/account"
)

// Store is a contract to Account..
//
//go:generate mockgen -destination ../../mock/store/account/account_store_mock.go -package mockStoreAccount . Store
type Store interface {
	StoreCreateAccount(ctx context.Context, documentNumber string) (*int64, error)
	GetAccount(ctx context.Context, accountId int64) (*accountModel.AccountResultQuery, error)
	GetAccountByDocumentNumber(ctx context.Context, documentNumber string) (*accountModel.AccountCountResultQuery, error)
	GetAllAccounts(ctx context.Context) ([]*accountModel.AccountResultQuery, error)
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

// StoreCreateAccount create a account..
func (a *storeImpl) StoreCreateAccount(ctx context.Context, documentNumber string) (*int64, error) {
	query := fmt.Sprintf("INSERT INTO accounts(document_number) VALUES ('%s')", documentNumber)
	res, err := a.db.Exec(query)

	if err != nil {
		log.Println("store.account.StoreCreateAccount.Exec", err.Error())
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println("store.account.StoreCreateAccount.LastInsertId", err.Error())
		return nil, err
	}
	return &lastId, nil
}

// GetAccount get a account..
func (a *storeImpl) GetAccount(ctx context.Context, accountId int64) (*accountModel.AccountResultQuery, error) {
	res, err := a.db.Query("SELECT * FROM accounts WHERE id = ?", accountId)
	if err != nil {
		log.Println("store.account.GetAccount.Query", err.Error())
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var account accountModel.AccountResultQuery
		err := res.Scan(&account.ID, &account.DocumentNumber)
		if err != nil {
			log.Println("store.account.GetAccount.Scan", err.Error())
			return nil, err
		}
		return &account, nil
	} else {
		return nil, accountModel.ErrorAccountNotFound
	}
}

// GetAccount get a account by document number..
func (a *storeImpl) GetAccountByDocumentNumber(ctx context.Context, documentNumber string) (*accountModel.AccountCountResultQuery, error) {
	res, err := a.db.Query("SELECT COUNT(*) as count FROM accounts WHERE document_number = ?", documentNumber)
	if err != nil {
		log.Println("store.account.GetAccountByDocumentNumber.Query", err.Error())
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var count accountModel.AccountCountResultQuery
		err := res.Scan(&count.Count)
		if err != nil {
			log.Println("store.account.GetAccountByDocumentNumber.Scan", err.Error())
			return nil, err
		}
		return &count, nil
	}

	return nil, accountModel.ErrorAccountNotFound
}

// GetAllAccounts get a account by document number..
func (a *storeImpl) GetAllAccounts(ctx context.Context) ([]*accountModel.AccountResultQuery, error) {
	res, err := a.db.Query("SELECT * FROM accounts")
	if err != nil {
		log.Println("store.account.GetAccountByDocumentNumber.Query", err.Error())
		return nil, err
	}
	defer res.Close()

	var listAccounts []*accountModel.AccountResultQuery

	for res.Next() {
		var item = accountModel.AccountResultQuery{}
		err = res.Scan(&item.ID, &item.DocumentNumber)
		if err != nil {
			log.Println("store.account.GetAccountByDocumentNumber.Scan", err.Error())
			return nil, err
		}
		listAccounts = append(listAccounts, &item)
	}
	return listAccounts, nil
}

package store

import (
	"database/sql"
	"log"

	"github.com/danilotadeu/pismo/store/account"
	"github.com/danilotadeu/pismo/store/transaction"
	_ "github.com/go-sql-driver/mysql"
)

//Container ...
type Container struct {
	Account     account.Store
	Transaction transaction.Store
}

//Register store container
func Register(db *sql.DB) *Container {
	container := &Container{
		Account:     account.NewStore(db),
		Transaction: transaction.NewStore(db),
	}

	log.Println("Registered -> Store")
	return container
}

package app

import (
	"log"

	"github.com/danilotadeu/pismo/app/account"
	"github.com/danilotadeu/pismo/app/transaction"
	"github.com/danilotadeu/pismo/store"
)

// Container ...
type Container struct {
	Account     account.App
	Transaction transaction.App
}

// Register app container
func Register(store *store.Container) *Container {
	container := &Container{
		Account:     account.NewApp(store),
		Transaction: transaction.NewApp(store),
	}
	log.Println("Registered -> App")
	return container
}

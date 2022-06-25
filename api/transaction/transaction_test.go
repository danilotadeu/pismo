package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/danilotadeu/pismo/app"
	accountModel "github.com/danilotadeu/pismo/model/account"
	modelTransaction "github.com/danilotadeu/pismo/model/transaction"
	transactionModel "github.com/danilotadeu/pismo/model/transaction"
	"github.com/danilotadeu/pismo/store"
	"github.com/danilotadeu/pismo/tests_personal"

	"github.com/gofiber/fiber/v2"
)

func TestTransactionHandler(t *testing.T) {
	store := store.Register(nil)
	apps := app.Register(store)
	api := apiImpl{
		apps: apps,
	}
	tests := []struct {
		name                   string
		method                 string
		contentType            string
		header                 map[string]string
		url                    string
		urlReq                 string
		handlerFunc            func(c *fiber.Ctx) error
		Input                  modelTransaction.TransactionRequest
		want                   int
		bodyShow               bool
		PrepareMockTransaction func(ctx context.Context, transactionBody transactionModel.TransactionRequest) (*int64, error)
	}{
		{
			name:        "Should return error bad request",
			method:      "POST",
			contentType: "application/json",
			header:      nil,
			url:         "/transactions",
			urlReq:      "/transactions",
			handlerFunc: api.transactionCreate,
			Input: modelTransaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          0,
			},
			want:     http.StatusBadRequest,
			bodyShow: true,
		},
		{
			name:        "Should return error with account dont exist",
			method:      "POST",
			contentType: "application/json",
			header:      nil,
			url:         "/transactions",
			urlReq:      "/transactions",
			handlerFunc: api.transactionCreate,
			Input: modelTransaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			want:     http.StatusNotFound,
			bodyShow: true,
			PrepareMockTransaction: func(ctx context.Context, transactionBody transactionModel.TransactionRequest) (*int64, error) {
				return nil, accountModel.ErrorAccountNotFound
			},
		},
		{
			name:        "Should return error with transaction type invalid",
			method:      "POST",
			contentType: "application/json",
			header:      nil,
			url:         "/transactions",
			urlReq:      "/transactions",
			handlerFunc: api.transactionCreate,
			Input: modelTransaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			want:     http.StatusNotFound,
			bodyShow: true,
			PrepareMockTransaction: func(ctx context.Context, transactionBody transactionModel.TransactionRequest) (*int64, error) {
				return nil, transactionModel.ErrorTransactionTypeNotFound
			},
		},
		{
			name:        "Should return error with internal server error",
			method:      "POST",
			contentType: "application/json",
			header:      nil,
			url:         "/transactions",
			urlReq:      "/transactions",
			handlerFunc: api.transactionCreate,
			Input: modelTransaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			want:     http.StatusInternalServerError,
			bodyShow: true,
			PrepareMockTransaction: func(ctx context.Context, transactionBody transactionModel.TransactionRequest) (*int64, error) {
				return nil, errors.New("Internal Server Error")
			},
		},
		{
			name:        "Should return success created",
			method:      "POST",
			contentType: "application/json",
			header:      nil,
			url:         "/transactions",
			urlReq:      "/transactions",
			handlerFunc: api.transactionCreate,
			Input: modelTransaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			want:     http.StatusCreated,
			bodyShow: true,
			PrepareMockTransaction: func(ctx context.Context, transactionBody transactionModel.TransactionRequest) (*int64, error) {
				var id int64
				id = 1
				return &id, nil
			},
		},
	}

	for _, tt := range tests {
		requestBody, err := json.Marshal(&tt.Input)
		if err != nil {
			t.Errorf("Error json.Marshal:%s", err.Error())
			return
		}
		t.Run(tt.name, func(t *testing.T) {
			if tt.bodyShow {
				fmt.Println("Json: ", string(requestBody))
			}
			createTransactionFunc = tt.PrepareMockTransaction
			tests_personal.TestNewRequest(t, tt.url, tt.urlReq, tt.method,
				tt.handlerFunc,
				bytes.NewBuffer(requestBody),
				tt.contentType, tt.header, tt.want, tt.bodyShow)
		})
	}
}

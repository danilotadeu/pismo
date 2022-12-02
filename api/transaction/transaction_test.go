package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danilotadeu/pismo/app"
	mockAppTransaction "github.com/danilotadeu/pismo/mock/app/transaction"
	accountModel "github.com/danilotadeu/pismo/model/account"
	transactionModel "github.com/danilotadeu/pismo/model/transaction"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestTransactionCreate(t *testing.T) {
	endpoint := "/transactions/"
	var transactionID int64 = 2
	cases := map[string]struct {
		InputBody          transactionModel.TransactionRequest
		NotParse           bool
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockTransactionApp *mockAppTransaction.MockApp)
	}{
		"should return success with transactionID created": {
			InputBody: transactionModel.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          4,
			},
			ExpectedErr: nil,
			PrepareMockApp: func(mockTransactionApp *mockAppTransaction.MockApp) {
				mockTransactionApp.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(&transactionID, nil)
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		"should return error to bind": {
			InputBody:          transactionModel.TransactionRequest{},
			ExpectedErr:        nil,
			NotParse:           true,
			PrepareMockApp:     func(mockTransactionApp *mockAppTransaction.MockApp) {},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should return error with bad request data": {
			InputBody: transactionModel.TransactionRequest{
				OperationTypeID: 5,
			},
			ExpectedErr:        nil,
			PrepareMockApp:     func(mockTransactionApp *mockAppTransaction.MockApp) {},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return error to create transaction": {
			InputBody: transactionModel.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          4,
			},
			ExpectedErr: nil,
			PrepareMockApp: func(mockTransactionApp *mockAppTransaction.MockApp) {
				mockTransactionApp.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should return error to create transaction with account not found": {
			InputBody: transactionModel.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          4,
			},
			ExpectedErr: nil,
			PrepareMockApp: func(mockTransactionApp *mockAppTransaction.MockApp) {
				mockTransactionApp.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil, accountModel.ErrorAccountNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockTransactionApp := mockAppTransaction.NewMockApp(ctrl)
			cs.PrepareMockApp(mockTransactionApp)

			h := apiImpl{
				apps: &app.Container{
					Transaction: mockTransactionApp,
				},
			}

			app := fiber.New()
			app.Post(endpoint, h.transactionCreate)
			var requestBody []byte
			var err error

			if !cs.NotParse {
				requestBody, err = json.Marshal(cs.InputBody)
				if err != nil {
					t.Errorf("Error json.Marshal: %s", err.Error())
					return
				}
			}

			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody)).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}

package transaction

import (
	"context"
	"fmt"
	"testing"

	mockStoreAccount "github.com/danilotadeu/pismo/mock/store/account"
	mockStoreTransaction "github.com/danilotadeu/pismo/mock/store/transaction"
	accountModel "github.com/danilotadeu/pismo/model/account"
	"github.com/danilotadeu/pismo/model/transaction"
	"github.com/danilotadeu/pismo/store"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateTransaction(t *testing.T) {
	var transactionIDExpected int64 = 2
	cases := map[string]struct {
		inputTransaction      transaction.TransactionRequest
		prepareMock           func(accountStore *mockStoreAccount.MockStore, transactionStore *mockStoreTransaction.MockStore)
		expectedTransactionID *int64
		expectedErr           error
	}{
		"should create a transaction with success": {
			inputTransaction: transaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			prepareMock: func(accountStore *mockStoreAccount.MockStore, transactionStore *mockStoreTransaction.MockStore) {
				accountStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, nil)
				var transactionID int64
				transactionID = 2
				transactionStore.EXPECT().CreateTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&transactionID, nil)
			},
			expectedTransactionID: &transactionIDExpected,
			expectedErr:           nil,
		},
		"should return error when getAccount": {
			inputTransaction: transaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			prepareMock: func(accountStore *mockStoreAccount.MockStore, transactionStore *mockStoreTransaction.MockStore) {
				accountStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

			},
			expectedTransactionID: nil,
			expectedErr:           fmt.Errorf("error"),
		},
		"should return error account not found when getAccount": {
			inputTransaction: transaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			prepareMock: func(accountStore *mockStoreAccount.MockStore, transactionStore *mockStoreTransaction.MockStore) {
				accountStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, accountModel.ErrorAccountNotFound)

			},
			expectedTransactionID: nil,
			expectedErr:           accountModel.ErrorAccountNotFound,
		},
		"should return error when a transaction type is not valid": {
			inputTransaction: transaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 5,
				Amount:          1,
			},
			prepareMock:           func(accountStore *mockStoreAccount.MockStore, transactionStore *mockStoreTransaction.MockStore) {},
			expectedTransactionID: nil,
			expectedErr:           transaction.ErrorTransactionTypeNotFound,
		},
		"should return error to create transaction": {
			inputTransaction: transaction.TransactionRequest{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          1,
			},
			prepareMock: func(accountStore *mockStoreAccount.MockStore, transactionStore *mockStoreTransaction.MockStore) {
				accountStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, nil)
				transactionStore.EXPECT().CreateTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedTransactionID: nil,
			expectedErr:           fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			accountStoreMock := mockStoreAccount.NewMockStore(ctrl)
			transactionStoreMock := mockStoreTransaction.NewMockStore(ctrl)

			cs.prepareMock(accountStoreMock, transactionStoreMock)
			app := NewApp(&store.Container{
				Account:     accountStoreMock,
				Transaction: transactionStoreMock,
			})

			// when
			transactionID, err := app.CreateTransaction(ctx, cs.inputTransaction)

			// then
			assert.Equal(t, cs.expectedTransactionID, transactionID)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

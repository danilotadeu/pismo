package account

import (
	"context"
	"fmt"
	"testing"

	mockStoreAccount "github.com/danilotadeu/pismo/mock/store/account"
	accountModel "github.com/danilotadeu/pismo/model/account"
	"github.com/danilotadeu/pismo/store"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateAccount(t *testing.T) {
	var accountID int64 = 21
	cases := map[string]struct {
		inputDocNumber    string
		prepareMock       func(accountStore *mockStoreAccount.MockStore)
		expectedAccountID *int64
		expectedErr       error
	}{
		"should create a account with success": {
			inputDocNumber: "12345",
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAccountByDocumentNumber(gomock.Any(), gomock.Any()).
					Return(&accountModel.AccountCountResultQuery{
						Count: 0,
					}, nil)
				accountStore.EXPECT().StoreCreateAccount(gomock.Any(), gomock.Any()).Return(&accountID, nil)
			},
			expectedAccountID: &accountID,
			expectedErr:       nil,
		},
		"should return error when get account by document number": {
			inputDocNumber: "12345",
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAccountByDocumentNumber(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedAccountID: nil,
			expectedErr:       fmt.Errorf("error"),
		},
		"should return error when account exists": {
			inputDocNumber: "12345",
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAccountByDocumentNumber(gomock.Any(), gomock.Any()).
					Return(&accountModel.AccountCountResultQuery{
						Count: 1,
					}, nil)
			},
			expectedAccountID: nil,
			expectedErr:       accountModel.ErrorAccountExists,
		},
		"should return error when create a account": {
			inputDocNumber: "12345",
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAccountByDocumentNumber(gomock.Any(), gomock.Any()).
					Return(&accountModel.AccountCountResultQuery{
						Count: 0,
					}, nil)
				accountStore.EXPECT().StoreCreateAccount(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedAccountID: nil,
			expectedErr:       fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			accountStoreMock := mockStoreAccount.NewMockStore(ctrl)

			cs.prepareMock(accountStoreMock)
			app := NewApp(&store.Container{
				Account:     accountStoreMock,
				Transaction: nil,
			})

			// when
			accountID, err := app.CreateAccount(ctx, cs.inputDocNumber)

			// then
			assert.Equal(t, cs.expectedAccountID, accountID)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestGetAccount(t *testing.T) {
	cases := map[string]struct {
		inputAccountID int64
		prepareMock    func(accountStore *mockStoreAccount.MockStore)
		expectedData   *accountModel.AccountResultQuery
		expectedErr    error
	}{
		"should get a account with success": {
			inputAccountID: 123,
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&accountModel.AccountResultQuery{
					ID:             1,
					DocumentNumber: "12345",
				}, nil)
			},
			expectedData: &accountModel.AccountResultQuery{
				ID:             1,
				DocumentNumber: "12345",
			},
			expectedErr: nil,
		},
		"should return error when get a account": {
			inputAccountID: 123,
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedData: nil,
			expectedErr:  fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			accountStoreMock := mockStoreAccount.NewMockStore(ctrl)

			cs.prepareMock(accountStoreMock)
			app := NewApp(&store.Container{
				Account:     accountStoreMock,
				Transaction: nil,
			})

			// when
			account, err := app.GetAccount(ctx, cs.inputAccountID)

			// then
			assert.Equal(t, cs.expectedData, account)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestGetAllAccounts(t *testing.T) {
	cases := map[string]struct {
		prepareMock  func(accountStore *mockStoreAccount.MockStore)
		expectedData []*accountModel.AccountResultQuery
		expectedErr  error
	}{
		"should get a all account with success": {
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAllAccounts(gomock.Any()).Return([]*accountModel.AccountResultQuery{
					{
						ID:             1,
						DocumentNumber: "1",
					},
					{
						ID:             2,
						DocumentNumber: "2",
					},
				}, nil)
			},
			expectedData: []*accountModel.AccountResultQuery{
				{
					ID:             1,
					DocumentNumber: "1",
				},
				{
					ID:             2,
					DocumentNumber: "2",
				},
			},
			expectedErr: nil,
		},
		"should return error when get a all account": {
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAllAccounts(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedData: nil,
			expectedErr:  fmt.Errorf("error"),
		},
		"should return error when get a all account is empty": {
			prepareMock: func(accountStore *mockStoreAccount.MockStore) {
				accountStore.EXPECT().GetAllAccounts(gomock.Any()).Return([]*accountModel.AccountResultQuery{}, nil)
			},
			expectedData: nil,
			expectedErr:  accountModel.ErrorAccountListIsEmpty,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			accountStoreMock := mockStoreAccount.NewMockStore(ctrl)

			cs.prepareMock(accountStoreMock)
			app := NewApp(&store.Container{
				Account:     accountStoreMock,
				Transaction: nil,
			})

			// when
			account, err := app.GetAllAccounts(ctx)

			// then
			assert.Equal(t, cs.expectedData, account)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

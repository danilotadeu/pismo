package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danilotadeu/pismo/app"
	mockAppAccount "github.com/danilotadeu/pismo/mock/app/account"
	accountModel "github.com/danilotadeu/pismo/model/account"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestHandlerAccountCreate(t *testing.T) {
	endpoint := "/accounts/"
	var accountID int64 = 2
	cases := map[string]struct {
		InputBody          accountModel.AccountRequest
		NotParse           bool
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockAccountApp *mockAppAccount.MockApp)
	}{
		"should return success with accountID created": {
			InputBody: accountModel.AccountRequest{
				DocumentNumber: "123878445456614",
			},
			ExpectedErr: nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(&accountID, nil)
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		"should return error when document number is empty": {
			InputBody: accountModel.AccountRequest{
				DocumentNumber: "12387",
			},
			ExpectedErr:        nil,
			PrepareMockApp:     func(mockAccountApp *mockAppAccount.MockApp) {},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return error to create account": {
			InputBody: accountModel.AccountRequest{
				DocumentNumber: "123878445456614",
			},
			ExpectedErr: nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should return error to create account with account exists": {
			InputBody: accountModel.AccountRequest{
				DocumentNumber: "123878445456614",
			},
			ExpectedErr: nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil, accountModel.ErrorAccountExists)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should return error to bind": {
			InputBody: accountModel.AccountRequest{
				DocumentNumber: "/",
			},
			NotParse:           true,
			ExpectedErr:        nil,
			PrepareMockApp:     func(mockAccountApp *mockAppAccount.MockApp) {},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockAccountApp := mockAppAccount.NewMockApp(ctrl)
			cs.PrepareMockApp(mockAccountApp)

			h := apiImpl{
				apps: &app.Container{
					Account: mockAccountApp,
				},
			}

			app := fiber.New()
			app.Post(endpoint, h.accountCreate)
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

func TestHandlerGetAllAccounts(t *testing.T) {
	endpoint := "/accounts/list"
	cases := map[string]struct {
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockAccountApp *mockAppAccount.MockApp)
	}{
		"should return success with all account list": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().GetAllAccounts(gomock.Any()).Return([]*accountModel.AccountResultQuery{
					{
						ID:             1,
						DocumentNumber: "12345",
					},
					{
						ID:             2,
						DocumentNumber: "12346",
					},
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
		},
		"should return error to get all account list": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().GetAllAccounts(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should return error to get all account list when account list is empty": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().GetAllAccounts(gomock.Any()).Return(nil, accountModel.ErrorAccountListIsEmpty)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockAccountApp := mockAppAccount.NewMockApp(ctrl)
			cs.PrepareMockApp(mockAccountApp)

			h := apiImpl{
				apps: &app.Container{
					Account: mockAccountApp,
				},
			}

			app := fiber.New()
			app.Get(endpoint, h.allAccounts)
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
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

func TestHandlerGetAccount(t *testing.T) {
	endpoint := "/accounts/:accountId"
	cases := map[string]struct {
		InputParamID       string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockAccountApp *mockAppAccount.MockApp)
	}{
		"should return success with account": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&accountModel.AccountResultQuery{
					ID:             1,
					DocumentNumber: "12345",
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
		},
		"should return error when param is not int": {
			InputParamID:       "A",
			ExpectedErr:        nil,
			PrepareMockApp:     func(mockAccountApp *mockAppAccount.MockApp) {},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should return error when get account": {
			InputParamID: "2",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should return error when get account not found": {
			InputParamID: "2",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockAccountApp *mockAppAccount.MockApp) {
				mockAccountApp.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, accountModel.ErrorAccountNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockAccountApp := mockAppAccount.NewMockApp(ctrl)
			cs.PrepareMockApp(mockAccountApp)

			h := apiImpl{
				apps: &app.Container{
					Account: mockAccountApp,
				},
			}

			app := fiber.New()
			app.Get(endpoint, h.account)
			req := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(endpoint, ":accountId", cs.InputParamID), nil).WithContext(ctx)
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

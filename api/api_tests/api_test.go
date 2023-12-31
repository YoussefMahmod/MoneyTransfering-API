package api_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/YoussefMahmod/MoneyTransfering-API/api"
	"github.com/YoussefMahmod/MoneyTransfering-API/models"
	"github.com/YoussefMahmod/MoneyTransfering-API/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
)

func TestPostAccount(t *testing.T) {
	// t.Parallel()
	testAccount, err := setupAccount()
	assert.NoError(t, err)

	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	res, err := restyPostAccount(server, testAccount)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	resultAccount, err := setupAccount()
	assert.NoError(t, err)

	err = json.Unmarshal(res.Body(), &resultAccount)
	assert.NoError(t, err)

	matchAccounts(t, testAccount, resultAccount)
}

func TestGetAccount(t *testing.T) {
	// t.Parallel()
	testAccount, err := setupAccount()
	assert.NoError(t, err)

	// mock server
	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	// post an account
	_, err = restyPostAccount(server, testAccount)
	assert.NoError(t, err)

	// try to get this account
	res, err := restyGetAccount(server, testAccount)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	resultAccount, err := setupAccount()
	assert.NoError(t, err)

	err = json.Unmarshal(res.Body(), &resultAccount)
	assert.NoError(t, err)

	// validate that the account returned correctly
	matchAccounts(t, testAccount, resultAccount)
}

func TestDelAccount(t *testing.T) {
	// t.Parallel()
	testAccount, err := setupAccount()
	assert.NoError(t, err)

	// mock server
	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	// post an account
	_, err = restyPostAccount(server, testAccount)
	assert.NoError(t, err)

	// try to delete it
	res, err := restyDelAccount(server, testAccount)
	assert.NoError(t, err)
	assert.EqualValues(t, res.RawResponse.StatusCode, 204)
	assert.Zero(t, res.RawResponse.Body)
}

func TestPatchAccount(t *testing.T) {
	// t.Parallel()
	testAccount, err := setupAccount()
	assert.NoError(t, err)

	// setup the server
	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	// post an account
	_, err = restyPostAccount(server, testAccount)
	assert.NoError(t, err)

	// alter the account
	testAccount.SetName(utils.RandStringBytes(7))
	testAccount.SetBalance(decimal.NewFromFloat(utils.RandFloat(99.9, 110.99)))

	// patch it
	res, err := restyPatchAccount(server, testAccount)
	assert.NoError(t, err)

	resultAccount, err := setupAccount()
	assert.NoError(t, err)

	err = json.Unmarshal(res.Body(), &resultAccount)
	assert.NoError(t, err)

	// check the results are correct
	matchAccounts(t, testAccount, resultAccount)
}

func TestBulkInsertionAccount(t *testing.T) {
	// t.Parallel()
	testBulkAccounts, err := setupMultipleAccounts()
	assert.NoError(t, err)
	assert.NotNil(t, testBulkAccounts)

	// mock the server
	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	// post bulk of accounts
	_, err = restyPostAccountsInBulk(server, testBulkAccounts)
	assert.NoError(t, err)

	// get all accounts
	res, err := restyGetAllAccounts(server)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// validate the results are correct
	resultBulkAccounts, err := models.NewListAccounts(res.Body())
	sort.Slice(resultBulkAccounts, func(i, j int) bool {
		return resultBulkAccounts[i].GetName() == resultBulkAccounts[j].GetName()
	})

	sort.Slice(resultBulkAccounts, func(i, j int) bool {
		return resultBulkAccounts[i].GetBalance().Cmp(resultBulkAccounts[j].GetBalance()) == -1
	})
	sort.Slice(testBulkAccounts, func(i, j int) bool {
		return testBulkAccounts[i].GetBalance().Cmp(testBulkAccounts[j].GetBalance()) == -1
	})

	assert.NoError(t, err)
	assert.Len(t, resultBulkAccounts, 5)
	for i := range resultBulkAccounts {
		matchAccounts(t, testBulkAccounts[i], resultBulkAccounts[i])
	}
}

func TestPostAccountTxn(t *testing.T) {
	// t.Parallel()
	senderAccount, err := setupAccount()
	assert.NoError(t, err)

	recieverAccount, err := setupAccount()
	assert.NoError(t, err)

	senderMoney := senderAccount.GetBalance()
	recieverMoney := recieverAccount.GetBalance()
	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	restyPostAccount(server, senderAccount)
	restyPostAccount(server, recieverAccount)

	testTxn, err := setupAccountTxn(senderAccount, recieverAccount)
	assert.NoError(t, err)
	assert.NotNil(t, testTxn)

	res, err := restyCreateAccountTxn(server, senderAccount, testTxn)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = restyGetAccount(server, senderAccount)
	assert.NoError(t, err)

	err = json.Unmarshal(res.Body(), &senderAccount)
	assert.NoError(t, err)

	res, err = restyGetAccount(server, recieverAccount)
	assert.NoError(t, err)

	err = json.Unmarshal(res.Body(), &recieverAccount)
	assert.NoError(t, err)

	assert.Equal(t, senderAccount.GetBalance(), decimal.NewFromInt(0))
	assert.Equal(t, recieverAccount.GetBalance(), senderMoney.Add(recieverMoney))
}

func TestGetAccountSentTxn(t *testing.T) {
	// t.Parallel()
	senderAccount, err := setupAccount()
	assert.NoError(t, err)

	recieverAccount, err := setupAccount()
	assert.NoError(t, err)

	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	restyPostAccount(server, senderAccount)
	restyPostAccount(server, recieverAccount)

	res, _ := restyGetAccountSentTxns(server, senderAccount)
	assert.Equal(t, res.StatusCode(), 404)

	senderTestTxn1, err := setupAccountTxn(senderAccount, recieverAccount)
	assert.NoError(t, err)
	restyCreateAccountTxn(server, senderAccount, senderTestTxn1)

	res, err = restyGetAccountSentTxns(server, senderAccount)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 200)

	var resBody []interface{}

	err = json.Unmarshal(res.Body(), &resBody)
	assert.NoError(t, err)
	assert.Len(t, resBody, 1)
}

func TestGetAccountRecievedTxn(t *testing.T) {
	// t.Parallel()
	senderAccount, err := setupAccount()
	assert.NoError(t, err)

	recieverAccount, err := setupAccount()
	assert.NoError(t, err)

	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	restyPostAccount(server, senderAccount)
	restyPostAccount(server, recieverAccount)

	res, err := restyGetAccountRecievedTxns(server, recieverAccount)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 404)

	senderTestTxn1, _ := setupAccountTxn(senderAccount, recieverAccount)
	restyCreateAccountTxn(server, senderAccount, senderTestTxn1)

	res, err = restyGetAccountRecievedTxns(server, recieverAccount)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 200)

	var resBody []interface{}

	err = json.Unmarshal(res.Body(), &resBody)
	assert.NoError(t, err)
	assert.Len(t, resBody, 1)
}

func TestGetAllTxns(t *testing.T) {
	// t.Parallel()
	acc1, err := setupAccount()
	assert.NoError(t, err)

	acc2, err := setupAccount()
	assert.NoError(t, err)

	acc3, err := setupAccount()
	assert.NoError(t, err)

	acc4, err := setupAccount()
	assert.NoError(t, err)

	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	restyPostAccount(server, acc1)
	restyPostAccount(server, acc2)
	restyPostAccount(server, acc3)
	restyPostAccount(server, acc4)

	senderTestTxn1, err := setupAccountTxn(acc1, acc2)
	assert.NoError(t, err)
	restyCreateAccountTxn(server, acc1, senderTestTxn1)

	senderTestTxn2, err := setupAccountTxn(acc3, acc4)
	assert.NoError(t, err)
	restyCreateAccountTxn(server, acc3, senderTestTxn2)

	res, err := restyGetAllTxns(server)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 200)

	var resBody []interface{}

	err = json.Unmarshal(res.Body(), &resBody)
	assert.NoError(t, err)
	assert.Len(t, resBody, 2)
}

func TestGetTxnByID(t *testing.T) {
	// t.Parallel()
	acc1, err := setupAccount()
	assert.NoError(t, err)

	acc2, err := setupAccount()
	assert.NoError(t, err)

	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	restyPostAccount(server, acc1)
	restyPostAccount(server, acc2)

	senderTestTxn1, err := setupAccountTxn(acc1, acc2)
	assert.NoError(t, err)
	restyCreateAccountTxn(server, acc1, senderTestTxn1)

	res, err := restyGetTxnByID(server, senderTestTxn1)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 200)

	resBody, err := models.NewTransaction(res.Body())

	assert.NoError(t, err)
	assert.Equal(t, resBody.GetSenderID(), acc1.GetID())
	assert.Equal(t, resBody.GetRecieverID(), senderTestTxn1.GetRecieverID())
}

func TestDelTxnByID(t *testing.T) {
	// t.Parallel()
	acc1, err := setupAccount()
	assert.NoError(t, err)

	acc2, err := setupAccount()
	assert.NoError(t, err)

	handler := api.NewServer()
	server := httptest.NewServer(handler.Router)

	defer server.Close()

	restyPostAccount(server, acc1)
	restyPostAccount(server, acc2)

	senderTestTxn1, err := setupAccountTxn(acc1, acc2)
	assert.NoError(t, err)
	restyCreateAccountTxn(server, acc1, senderTestTxn1)

	res, err := restyDelTxnByID(server, senderTestTxn1)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 204)

	res, err = restyGetTxnByID(server, senderTestTxn1)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode(), 404)
}

func setupAccount() (models.IAccount, error) {
	name := utils.RandStringBytes(10)
	balance := decimal.NewFromFloat(utils.RandFloat(10.99, 1000.99))
	byt := []byte(fmt.Sprintf("{\"name\":\"%v\",\"balance\":%v}", name, balance))

	return models.NewAccount(byt)
}

func setupMultipleAccounts() ([]models.IAccount, error) {
	var namesList [5]string
	var balancesList [5]decimal.Decimal
	for i := 0; i < 5; i++ {
		namesList[i] = utils.RandStringBytes(10)
		balancesList[i] = decimal.NewFromFloat(utils.RandFloat(10.99, 1000.99))
	}

	byt := []byte(fmt.Sprintf("[{\"name\":\"%v\",\"balance\":%v},{\"name\":\"%v\",\"balance\":%v},{\"name\":\"%v\",\"balance\":%v},{\"name\":\"%v\",\"balance\":%v},{\"name\":\"%v\",\"balance\":%v}]",
		namesList[0], balancesList[0],
		namesList[1], balancesList[1],
		namesList[2], balancesList[2],
		namesList[3], balancesList[3],
		namesList[4], balancesList[4]))

	return models.NewListAccounts(byt)
}

func restyPostAccount(server *httptest.Server, account models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(account).
		Post(server.URL + "/api/v1/accounts/")
}

func restyPostAccountsInBulk(server *httptest.Server, bulk []models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bulk).
		Post(server.URL + "/api/v1/accounts/bulk")
}

func restyGetAllAccounts(server *httptest.Server) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Get(server.URL + "/api/v1/accounts/")
}

func restyGetAccount(server *httptest.Server, account models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Get(server.URL + fmt.Sprintf("/api/v1/accounts/%v", account.GetID()))
}

func restyDelAccount(server *httptest.Server, account models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Delete(server.URL + fmt.Sprintf("/api/v1/accounts/%v", account.GetID()))
}

func restyPatchAccount(server *httptest.Server, account models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(account).
		Patch(server.URL + fmt.Sprintf("/api/v1/accounts/%v", account.GetID()))
}

func setupAccountTxn(sender models.IAccount, reciever models.IAccount) (models.ITransaction, error) {
	byt := []byte(fmt.Sprintf("{\"reciever_id\":\"%v\",\"amount\":\"%v\"}", reciever.GetID(), sender.GetBalance()))

	return models.NewTransaction(byt)
}

func restyCreateAccountTxn(server *httptest.Server, account models.IAccount, txn models.ITransaction) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(txn).
		Post(server.URL + fmt.Sprintf("/api/v1/accounts/%v/transactions", account.GetID()))
}

func restyGetAccountSentTxns(server *httptest.Server, account models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Get(server.URL + fmt.Sprintf("/api/v1/accounts/%v/transactions/sent", account.GetID()))
}

func restyGetAccountRecievedTxns(server *httptest.Server, account models.IAccount) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Get(server.URL + fmt.Sprintf("/api/v1/accounts/%v/transactions/recieved", account.GetID()))
}

func restyGetAllTxns(server *httptest.Server) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Get(server.URL + "/api/v1/transactions/")
}

func restyGetTxnByID(server *httptest.Server, txn models.ITransaction) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Get(server.URL + fmt.Sprintf("/api/v1/transactions/%v", txn.GetID()))
}

func restyDelTxnByID(server *httptest.Server, txn models.ITransaction) (*resty.Response, error) {
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		Delete(server.URL + fmt.Sprintf("/api/v1/transactions/%v", txn.GetID()))
}

func matchAccounts(t *testing.T, acc1 models.IAccount, acc2 models.IAccount) {
	assert.Equal(t, acc1.GetID(), acc2.GetID())
	assert.Equal(t, acc1.GetName(), acc2.GetName())
	assert.Equal(t, acc1.GetBalance(), acc2.GetBalance())
}

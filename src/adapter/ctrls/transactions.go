package ctrls

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/usecases"
	"github.com/wakuwaku3/account-book.api/src/enterprise/helpers"

	"github.com/wakuwaku3/account-book.api/src/adapter/ctrls/responses"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/core"

	"github.com/labstack/echo"
)

type (
	transactions struct {
		useCase usecases.Transactions
		clock   helpers.Clock
	}
	// Transactions is TransactionsController
	Transactions interface {
		GetTransactions(c echo.Context) error
		GetTransaction(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
	getTransactionsResponse struct {
		Transactions []getTransactionResponse `json:"transactions"`
	}
	getTransactionResponse struct {
		TransactionID string    `json:"id"`
		Amount        int       `json:"amount"`
		Category      int       `json:"categoryId,string"`
		Date          time.Time `json:"date"`
		Notes         *string   `json:"notes,omitempty"`
		Editable      bool      `json:"editable"`
	}
	transactionRequest struct {
		Amount   *int    `json:"amount,omitempty"`
		Category *int    `json:"categoryId,string,omitempty"`
		Notes    *string `json:"notes,omitempty"`
	}
	createTransactionResponse struct {
		TransactionID string `json:"id"`
	}
)

// NewTransactions is create instance
func NewTransactions(useCase usecases.Transactions, clock helpers.Clock) Transactions {
	return &transactions{useCase, clock}
}

func (t *transactions) GetTransactions(c echo.Context) error {
	var err error
	selectedMonth := t.clock.Now()
	month := c.QueryParam("month")
	if month != "" {
		selectedMonth, err = time.Parse("2006-01-02", month)
		if err != nil {
			return err
		}
	}
	res, err := t.useCase.GetTransactions(&usecases.GetTransactionsArgs{
		SelectedMonth: selectedMonth,
	})
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, getTransactionsResponse{
		Transactions: convertTransactions(res.Transactions),
	})
}
func convertTransactions(transactions []usecases.GetTransactionResult) []getTransactionResponse {
	x := make([]getTransactionResponse, len(transactions))
	for i, transaction := range transactions {
		x[i] = convertTransaction(transaction)
	}
	return x
}
func convertTransaction(transaction usecases.GetTransactionResult) getTransactionResponse {
	return getTransactionResponse{
		TransactionID: transaction.TransactionID,
		Amount:        transaction.Amount,
		Category:      transaction.Category,
		Date:          transaction.Date,
		Notes:         transaction.Notes,
		Editable:      transaction.Editable,
	}
}
func (t *transactions) GetTransaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	res, err := t.useCase.GetTransaction(&id)
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertTransaction(*res))
}
func (t *transactions) Create(c echo.Context) error {
	request := new(transactionRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.Create(request.convert())
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, createTransactionResponse{
		TransactionID: res.TransactionID,
	})
}
func (t *transactionRequest) convert() *usecases.TransactionArgs {
	return &usecases.TransactionArgs{
		Amount:   t.Amount,
		Category: t.Category,
		Notes:    t.Notes,
	}
}

func (t *transactions) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	request := new(transactionRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if err := t.useCase.Update(&id, request.convert()); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *transactions) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	if err := t.useCase.Delete(&id); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}

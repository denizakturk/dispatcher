package transaction

import (
	"github.com/denizakturk/dispatcher/middleware"
	"github.com/denizakturk/dispatcher/model"
)

type Transaction[T any] interface {
	*T
	SetRequest(body []byte) error
	GetRequest() any
	GetResponse() any
	Transact() error
	GetRunables() []middleware.MiddlewareRunable
	AddRunable(runable middleware.MiddlewareRunable)
	SetRunables(runables []middleware.MiddlewareRunable)
}

type TransactionBucketItemInterface interface {
	GetName() string
	GetTransaction() model.ServerInterface
}

type TransactionBUcketItem struct {
	Name        string
	Transaction TransactionBucketItemInterface
}

type TransactionBucketItem struct {
	Name        string
	Transaction model.ServerInterface
}

func (t TransactionBucketItem) GetName() string {
	return t.Name
}

func (t TransactionBucketItem) GetTransaction() model.ServerInterface {
	return t.Transaction
}

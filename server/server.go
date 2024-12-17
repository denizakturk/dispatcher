package server

import (
	"dispatcher/department"
	"dispatcher/model"
	"dispatcher/transaction"
	"encoding/json"
	"log"
	"net/http"
)

type Server[T any, TI transaction.Transaction[T]] struct {
}

func (Server[T, TI]) Init(document model.Document) model.Document {
	var ta TI = new(T)
	jsonByteData, err := json.Marshal(document.Form)
	if err != nil {
		outputErrDoc := model.Document{Department: document.Department, Transaction: document.Transaction, Error: err}
		return outputErrDoc
	}

	if err != nil {
		outputErrDoc := model.Document{Department: document.Department, Transaction: document.Transaction, Error: err}
		return outputErrDoc
	}
	validator := model.DocumentFormValidater{Request: string(jsonByteData)}
	err = validator.Validate(ta.GetRequest())
	if err != nil {
		outputErrDoc := model.Document{Department: document.Department, Transaction: document.Transaction, Error: err.Error()}
		return outputErrDoc
	}
	err = ta.Transact()
	if err != nil {
		outputErrDoc := model.Document{Department: document.Department, Transaction: document.Transaction, Error: err}
		return outputErrDoc
	}
	ta.SetRequest(jsonByteData)
	err = ta.Transact()
	if err != nil {
		outputErrDoc := model.Document{Department: document.Department, Transaction: document.Transaction, Error: err}
		return outputErrDoc
	}
	document.Output = ta.GetResponse()

	return document
}

func ServJsonApi(register *department.RegisterDispatcher) {
	http.Handle("/", register)
	log.Fatal(http.ListenAndServe(":"+register.Port, nil))
}

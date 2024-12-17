package department

import (
	"dispatcher/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func RegisterMainFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	byteJson, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	document := &model.Document{}
	err = json.Unmarshal(byteJson, document)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, department := range DispatcherHolder {
		if department.Name == document.Department {
			for _, v := range department.Transactions {
				if v.GetName() == document.Transaction {
					outputDoc := v.GetTransaction().Init(*document)

					response, err := json.Marshal(outputDoc)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					fmt.Fprint(w, string(response))
					return
				}
			}
		}
	}
	outputDoc := model.Document{Department: document.Department, Transaction: document.Transaction, Error: errors.New("transaction not found")}
	w.WriteHeader(http.StatusBadRequest)
	response, err := json.Marshal(outputDoc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, string(response))
	return
}

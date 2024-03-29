package registrant

import (
	"net/http"

	"github.com/denizakturk/dispatcher/handling"
	"github.com/denizakturk/dispatcher/model"
)

type RequestFilter func(rw http.ResponseWriter, req *http.Request) bool

var RequestFilterHolder []RequestFilter

func NewRegisterDispatch() RegisterDispatcher {
	dispatch := RegisterDispatcher{}

	dispatch.MainFunc = func(rw http.ResponseWriter, req *http.Request) {
		if RequestFilterHolder != nil {
			for _, filter := range RequestFilterHolder {
				if !filter(rw, req) {
					return
				}
			}
		}
		body, err := handling.RequestHandle(req)
		if err != nil {
			errDoc := &model.Document{}
			documentarist := model.NewDocumentarist(rw, errDoc)
			documentarist.WriteError(err)
			return
		}
		inputDoc, err := handling.RequestBodyToDocument(body)
		if err != nil {
			errDoc := &model.Document{}
			documentarist := model.NewDocumentarist(rw, errDoc)
			documentarist.WriteError(err)
			return
		}
		documentarist := model.NewDocumentarist(rw, inputDoc)
		documentation := NewDocumentation(&documentarist)

		documentation.DocumentEnforcer()
		documentarist.Write()
	}

	dispatch.Port = "9000"
	return dispatch
}

type RegisterDispatcher struct {
	MainFunc func(http.ResponseWriter, *http.Request)
	Port     string
}

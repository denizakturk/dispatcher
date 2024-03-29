package registrant

import (
	"encoding/json"
	"errors"

	"github.com/denizakturk/dispatcher/constants"
	"github.com/denizakturk/dispatcher/handling"
	"github.com/denizakturk/dispatcher/model"
)

func NewDocumentation(documentairst *model.Documentarist) *Documentation {
	documentation := &Documentation{Documentarist: documentairst}
	return documentation
}

type Documentation struct {
	Documentarist *model.Documentarist
	Transaction   model.Transaction
}

func (r *Documentation) DocumentEnforcer() {
	r.TransactionEnforcer(r.Documentarist.Input, r.Documentarist.Output)
}

func (r Documentation) TransactionEnforcer(inputDoc *model.Document, outputDoc *model.Document) {
	outputDoc.Department = inputDoc.Department
	outputDoc.Transaction = inputDoc.Transaction
	var lastResponse interface{} = nil
	transactionHolder, err := r.TransactionMatcher(inputDoc)
	if err != nil {
		outputDoc.Type = constants.DOC_TYPE_ERROR
		outputDoc.Error = err.Error()
		return
	}

	if inputDoc.Type == constants.DOC_TYPE_PROCEDURE {
		r.TransactionProceduring(transactionHolder, outputDoc)
	} else {

		err = r.ParameterValidator(transactionHolder, inputDoc.Form)
		if err != nil {
			outputDoc.Type = constants.DOC_TYPE_ERROR
			outputDoc.Error = err.Error()
			return
		}
		err = r.InitTransaction(transactionHolder)
		if err != nil {
			outputDoc.Type = constants.DOC_TYPE_ERROR
			outputDoc.Error = err.Error()
			return
		}

		err = r.DocumentVerification(inputDoc, transactionHolder)

		if err != nil {
			outputDoc.Type = constants.DOC_TYPE_ERROR
			outputDoc.Error = err.Error()
			return
		}

		err = r.ProcessTransact(transactionHolder, outputDoc)
		if err != nil {
			outputDoc.Type = constants.DOC_TYPE_ERROR
			outputDoc.Error = err.Error()
			return
		}

		if outputDoc.Type == constants.DOC_TYPE_RESULT {
			lastResponse = outputDoc.Output
		}

		if lastResponse != nil && inputDoc.ChainRequestOption != nil {
			lastResponse = r.ResponseTransformer(lastResponse, inputDoc.ChainRequestOption)
		}
	}

	if inputDoc.Dispatchings != nil {
		for _, val := range inputDoc.Dispatchings {
			val.Form.FromInterface(lastResponse)
			outputDocNew := &model.Document{}
			if val.Security == nil || val.Security.Licence == "" {
				val.Security = &model.Security{}
				val.Security.Licence = inputDoc.Security.Licence
			}
			r.TransactionEnforcer(val, outputDocNew)
			outputDoc.Dispatchings = append(outputDoc.Dispatchings, outputDocNew)
			if outputDocNew.Type == constants.DOC_TYPE_ERROR {
				return
			}
		}
	}
}

func (r Documentation) TransactionMatcher(inputDoc *model.Document) (transactionHolder model.TransactionHolder, err error) {

	for _, department := range DepartmentRegistering {
		if department.Name == inputDoc.Department {
			for _, findedTransaction := range department.Transactions {
				if findedTransaction.Name == inputDoc.Transaction {
					return findedTransaction, nil
				}
			}
		}
	}

	return transactionHolder, errors.New(constants.TRANSACTION_NOT_FOUND)
}

func (r *Documentation) InitTransaction(transactionHolder model.TransactionHolder) (err error) {
	r.Transaction, err = transactionHolder.InitTransaction(*r.Documentarist.Input)
	return
}

func (r Documentation) ParameterValidator(transactionHolder model.TransactionHolder, form model.DocumentForm) error {
	formString, _ := formToString(form)
	documentFormValidater := handling.DocumentFormValidater{Request: formString}
	reqType := transactionHolder.Type.GetRequest()

	return documentFormValidater.Validate(reqType)
}

func (r Documentation) DocumentVerification(inputDoc *model.Document, transactionHolder model.TransactionHolder) error {
	if transactionHolder.Options.GetOptions().Security.LicenceChecker == true {
		if inputDoc.Security != nil && inputDoc.Security.Licence != "" {
			isValidToken := transactionHolder.LicenceValidator(inputDoc.Security.Licence)
			if !isValidToken {
				return errors.New("licence not valid")
			}
		} else {
			return errors.New("licence not found")
		}
	}
	return nil
}

func (r Documentation) TransactionProceduring(transactionHolder model.TransactionHolder, outputDoc *model.Document) {
	outputDoc.Output = model.Analysis(transactionHolder.Type.GetResponse(), &[]string{})
	outputDoc.Procedure = model.Analysis(transactionHolder.Type.GetRequest(), &[]string{})
	transactionOptions := transactionHolder.Options.GetOptions()
	outputDoc.Options = &transactionOptions
	outputDoc.Type = constants.DOC_TYPE_PROCEDURE
}

func (r Documentation) ProcessTransact(transactionHolder model.TransactionHolder, outputDoc *model.Document) error {
	response, err := r.Transaction.Transact()
	if err != nil {
		return err
	}
	outputDoc.Output = response
	outputDoc.Type = constants.DOC_TYPE_RESULT
	return err
}

func (r Documentation) ResponseTransformer(response interface{}, chainRequestOption model.ChainRequestOption) (responseMap map[string]interface{}) {
	responseByte, _ := json.Marshal(response)
	json.Unmarshal(responseByte, &responseMap)
	for key, val := range chainRequestOption {
		if _, ok := responseMap[key]; ok {
			responseMap[val.(string)] = responseMap[key]
		}
	}

	return
}

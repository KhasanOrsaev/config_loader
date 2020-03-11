package api

import (
	"encoding/json"
	"github.com/go-errors/errors"
)

func (response *Response) ErrorWrap(err error, code int) []byte {
	response.Status.Error = &ErrorResponse{Error: err.Error(), StackTrace:err.(*errors.Error).ErrorStack()}
	response.Status.Code = code
	res, _ := json.Marshal(response)
	return res
}
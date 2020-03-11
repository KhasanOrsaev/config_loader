package api

import "encoding/json"

type (
	Response struct {
		Status struct{
			Code int `json:"code"`
			Error *ErrorResponse `json:"error"`
		} `json:"status"`
		Content *json.RawMessage `json:"content"`
	}

	ErrorResponse struct {
		Error string `json:"message"`
		StackTrace string `json:"stack_trace"`
	}
)


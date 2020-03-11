package methods

import (
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/api"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/config"
	"encoding/json"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
)

type ConfGetRequest struct {
	ServiceName string `json:"service_name"`
	ServiceEnv string `json:"service_env"`
}

type CondGetResponse struct {
	Result json.RawMessage `json:"result"`
}

func ConfGetHandler(w http.ResponseWriter, r *http.Request)  {
	response := api.Response{}
	conf:= config.GetConfig()
	response.Status.Code = 0
	b,err := ioutil.ReadAll(r.Body)
	if err!=nil {
		e := errors.Wrap(err, -1)
		api.ErrorHandler("error on reading request", e, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
		res := response.ErrorWrap(e, 1)
		w.WriteHeader(400)
		w.Write(res)
		return
	}
	if len(b)==0 {
		e := errors.New("empty request")
		api.ErrorHandler("error on reading request", e, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
		res := response.ErrorWrap(e, 1)
		w.WriteHeader(400)
		w.Write(res)
		return
	}
	var request ConfGetRequest
	err = json.Unmarshal(b, &request)
	if err!=nil {
		e := errors.Wrap(err, -1)
		api.ErrorHandler("error on parsing request", err, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
		res := response.ErrorWrap(e, 1)
		w.WriteHeader(400)
		w.Write(res)
		return
	}
	if request.ServiceName == "" || request.ServiceEnv == "" {
		e := errors.New("error. service name or service env is empty.")
		api.ErrorHandler("error on validate", err, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
		res := response.ErrorWrap(e, 1)
		w.WriteHeader(400)
		w.Write(res)
		return
	}
	var result []uint8
	err = conf.DB.QueryRow(conf.Webservice.Functions["config-get"], request.ServiceName, request.ServiceEnv).Scan(&result)
	if err!=nil {
		e := errors.Wrap(err, -1)
		api.ErrorHandler("error on query", err, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
		res := response.ErrorWrap(e, 3)
		w.WriteHeader(500)
		w.Write(res)
		return
	}
	response.Content = (*json.RawMessage)(&result)
	res, err := json.Marshal(response)
	if err!=nil {
		e := errors.Wrap(err, -1)
		api.ErrorHandler("error on marshaling response", err, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
		res := response.ErrorWrap(e, 3)
		w.WriteHeader(500)
		w.Write(res)
		return
	}
	w.WriteHeader(200)
	w.Write(res)
}

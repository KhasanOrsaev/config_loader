package methods

import (
	"encoding/json"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/api"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/config"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
)
type ConfApplyRequest struct {
	ServiceName string `json:"service_name"`
	ServiceEnv string `json:"service_env"`
	Data json.RawMessage `json:"data"`
}

// ConfApplyHandler метод для записи новых событий
func ConfApplyHandler(w http.ResponseWriter, r *http.Request)  {
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
	var request ConfApplyRequest
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
	err = conf.DB.QueryRow(conf.Webservice.Functions["config-apply"], r.Context().Value("user_id").(string),
		request.ServiceName, request.ServiceEnv, string(request.Data)).Scan(&result)
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

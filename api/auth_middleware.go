package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/config"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthData struct {
	Scopes []string `json:"scopes"`
	UserID *uuid.UUID `json:"user_id"`
}

func AuthOnly(h http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := isAuth(r, role)
		if err != nil {
			ErrorHandler("auth error", err, map[string]interface{}{"stack_trace":err.(*errors.Error).ErrorStack()}, false)
			response := Response{}
			res := response.ErrorWrap(err, 1)
			w.WriteHeader(401)
			w.Write(res)
		}
		if userID == nil {
			response := Response{}
			res := response.ErrorWrap(errors.New("Not authorized"), 1)
			w.WriteHeader(401)
			w.Write(res)
			return
		}
		h(w, r.WithContext(context.WithValue(r.Context(), "user_id", userID.String())))
	}
}

func isAuth(r *http.Request, role string) (*uuid.UUID, error) {
	conf := config.GetConfig()
	httpClient := &http.Client{}
	if r.Header.Get("Authorization") == "" {
		return nil, errors.New("Empty authorization header")
	}
	header := strings.Split(r.Header.Get("Authorization"), " ")
	var body []byte
	if header[0]=="Basic" {
		body,_ = json.Marshal(map[string][]string{"scopes":strings.Split(role,",")})
	}
	url :=  "http://" + config.GetEnv(conf.Login.Host, "localhost") + ":" + config.GetEnv(conf.Login.Port, "8080") + conf.Login.Route
	request, err := http.NewRequest("post", url, bytes.NewReader(body))
	if err!= nil {
		return nil,errors.Wrap(err, -1)
	}
	request.Header = r.Header
	response, err := httpClient.Do(request)
	if err!= nil {
		return nil,errors.Wrap(err, -1)
	}
	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}
	fmt.Println(response.Status)
	var authResponse Response
	b, err := ioutil.ReadAll(response.Body)
	if err!= nil {
		return nil,errors.Wrap(err, -1)
	}
	err = json.Unmarshal(b, &authResponse)
	if err!= nil {
		return nil,errors.Wrap(err, -1)
	}
	if authResponse.Content == nil {
		return nil, nil
	}
	if authResponse.Status.Error != nil {
		e := errors.New(authResponse.Status.Error.Error)
		ErrorHandler("error on auth", err, map[string]interface{}{"stack_trace":e.ErrorStack()}, false)
	}
	data := AuthData{}
	err = json.Unmarshal(*authResponse.Content, &data)
	if err!= nil {
		return nil,errors.Wrap(err, -1)
	}
	if strInArray(data.Scopes, role) {
		return data.UserID, nil
	}
	return nil,nil
}

func strInArray(arr []string, s string) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

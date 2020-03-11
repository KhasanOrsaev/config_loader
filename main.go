package main

import (
	"flag"
	"fmt"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/api"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/config"
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/methods"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
	"time"
)

func init()  {
	var confPath string
	flag.StringVar(&confPath, "conf", "./config.yaml", "configuration file")
	flag.Parse()
	file,err := ioutil.ReadFile(confPath)
	if err != nil {
		api.ErrorHandler("error on opening config file", err, nil, true)
	}
	err = config.InitConfig(file)
	if err != nil {
		api.ErrorHandler("error on initing config", err, nil, true)
	}
}

func main() {
	conf := config.GetConfig()
	logger := config.GetLogger()
	logger.Info("service start at:", time.Now().Format(time.RFC3339))
	// убрал проверку на права, так как не нужны, а если нужны, то раскомментить
	//http.HandleFunc(conf.Webservice.Routes["event-get"], api.AuthOnly(methods.EventGetHandler, "event:read"))
	http.HandleFunc(conf.Webservice.Routes["config-apply"], api.AuthOnly(methods.ConfApplyHandler, "config:write"))
	fmt.Println()
	http.HandleFunc(conf.Webservice.Routes["config-get"], methods.ConfGetHandler)
	//http.HandleFunc(conf.Webservice.Routes["config-apply"], methods.ConfApplyHandler)
	err := http.ListenAndServe(":" + conf.Webservice.Port, nil)
	if err != nil {
		e := errors.Wrap(err, -1)
		api.ErrorHandler("error on starting server", e, nil, true)
	}
}
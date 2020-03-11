package api

import (
	"git.fin-dev.ru/dmp/cs_conf_webservice.git/config"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(msg string, err error, add map[string]interface{}, fatal bool) {
	logger := config.GetLogger()
	// если возникли проблемы при инициализации логера, используем дефолтный
	if logger == nil {
		logger = logrus.WithField("module", "events-api-webservice")
	}
	logger.WithField("message", msg)
	for i := range add {
		logger.WithField(i, add[i])
	}
	if fatal {
		logger.Fatal(err.Error())
	} else {
		logger.Error(err.Error())
	}
}

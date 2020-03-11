package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type myLogFormatter struct {
}

// форматирование логов
func (f *myLogFormatter) Format(entry *log.Entry) ([]byte, error) {

	dt := time.Now().String()
	extra, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf(config.Log.Format, err)
	}
	l := fmt.Sprintf("[%s] %s.%s message: %s context: %s extra: %s", dt, config.Log.ServiceName, entry.Level.String(), entry.Message, "[]", string(extra))
	return append([]byte(l), '\n'), nil
}

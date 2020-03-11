package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Connections map[string]struct{
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Database string `yaml:"database"`
		User string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"connections"`
	Webservice struct{
		Port string `yaml:"port"`
		Routes map[string]string `yaml:"routes"`
		Functions map[string]string `yaml:"functions"`
	} `yaml:"webservice"`
	DB *sql.DB
	Login struct{
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Route string `yaml:"route"`
	} `yaml:"auth-webservice"`
	Log struct{
		Format      string `yaml:"format"`
		ServiceName string `yaml:"service_name"`
		Level       uint32 `yaml:"level"`
	} `yaml:"log"`
}

var (
	config Config
	logger *log.Entry
)

func InitConfig(f []byte) error {
	err := yaml.Unmarshal(f, &config)
	config.DB, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		GetEnv(config.Connections["postgres"].User, "events"),
		GetEnv(config.Connections["postgres"].Password, "LGbQbugEoYBL1ZQFLAgpDKYk"),
		GetEnv(config.Connections["postgres"].Host, "localhost"),
		GetEnv(config.Connections["postgres"].Database, "test_events")))
	if err != nil {
		return err
	}
	config.Webservice.Port = GetEnv(config.Webservice.Port, "8083")
	// init logging
	if _, err := os.Stat("./var/log"); os.IsNotExist(err) {
		_ = os.MkdirAll("./var/log", os.ModePerm)
	}
	file, err := os.OpenFile("./var/log/"+config.Log.ServiceName+"_"+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	log.SetFormatter(&myLogFormatter{})
	log.SetLevel(log.Level(config.Log.Level))
	logger = log.WithField("module", config.Log.ServiceName)
	//
	return nil
}

func GetConfig() *Config {
	return &config
}

func GetLogger() *log.Entry {
	return logger
}
// Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

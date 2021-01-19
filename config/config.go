package config

import (
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/vrischmann/envconfig"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	GrpcPort                string        `envconfig:"default=9504"`
	HttpPort                string        `envconfig:"default=9505"`
	APIHost                 string        `envconfig:"default=http://api.com"`
	BreakerErrorThreshold   int           `envconfig:"default=3"`
	BreakerSuccessThreshold int           `envconfig:"default=3"`
	BreakerTimeout          time.Duration `envconfig:"default=15s"`
	HttpClientTimeout       time.Duration `envconfig:"default=15s"`
	DatabaseFile            string        `envconfig:"default=logwatcher.db"`
	SecretKey               string        `envconfig:"default=mnNska26ajAS2smnNska26ajAS2sVcE3"`
	MaxWorkerPool           int           `envconfig:"default=100"`
}

func readFromFileAndEnv(conf interface{}) (err error) {
	file, err := os.Open("appsettings.json")
	if err == nil {
		defer file.Close()
		data, inErr := ioutil.ReadAll(file)
		if inErr != nil {
			err = inErr
			return
		}
		maps := make(map[string]string)
		inErr = jsoniter.Unmarshal(data, &maps)
		if inErr != nil {
			err = inErr
			return
		}
		for k, v := range maps {
			inErr = os.Setenv(k, v)
			if inErr != nil {
				err = inErr
				return
			}
		}
	} else {
		logger.Warn("%v", err)
	}

	err = envconfig.Init(conf)
	if err != nil {
		return
	}
	return
}

func New() (conf *Config, err error) {
	conf = new(Config)
	err = readFromFileAndEnv(&conf)
	return
}

package log

import (
	"encoding/json"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"go.uber.org/zap"
)

func NewStdLogger() *kratoszap.Logger {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "levelKey": "level",
		  "levelEncoder": "capital"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	zapLogger := zap.Must(cfg.Build())

	klogger := kratoszap.NewLogger(zapLogger)

	return klogger
}

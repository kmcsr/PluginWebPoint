
package api

import (
	"os"

	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
)

var DEBUG = os.Getenv("DEBUG") == "true"
var loger = initLogger()

func initLogger()(loger logger.Logger){
	loger = logrus.Logger
	if DEBUG {
		loger.SetLevel(logger.TraceLevel)
	}else{
		loger.SetLevel(logger.InfoLevel)
	}
	loger.Debug("API Logger debug mode on")
	return
}

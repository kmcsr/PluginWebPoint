
package api

import (
	"os"

	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/golog"
)

var loger = initLogger()

func initLogger()(loger logger.Logger){
	loger = golog.Logger
	if os.Getenv("DEBUG") == "true" {
		loger.SetLevel(logger.TraceLevel)
	}else{
		loger.SetLevel(logger.InfoLevel)
	}
	return
}

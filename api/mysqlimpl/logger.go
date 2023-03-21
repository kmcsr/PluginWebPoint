
package mysqlimpl

import (
	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"

	"github.com/kmcsr/PluginWebPoint/api"
)

var loger = initLogger()

func initLogger()(loger logger.Logger){
	loger = logrus.Logger
	if api.DEBUG {
		loger.SetLevel(logger.TraceLevel)
	}else{
		loger.SetLevel(logger.InfoLevel)
	}
	return
}

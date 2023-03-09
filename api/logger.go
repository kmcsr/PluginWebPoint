
package api

import (
	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/golog"
)

var loger logger.Logger = golog.Logger

func SetLogger(newLogger logger.Logger){
	loger = newLogger
}

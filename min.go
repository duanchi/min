package min

import (
	"github.com/gin-gonic/gin"
	"github.com/xormplus/xorm"
	"github.com/duanchi/min/bean"
	"github.com/duanchi/min/log"
)

var HttpServer *gin.Engine
var Db *xorm.Engine
var Config interface{}
var Log *log.Logger
var DbMap map[string]*xorm.Engine

func GetBean(name string) interface{} {
	return bean.Get(name)
}


